package service

import (
	"context"
	"errors"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/subprocess"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"sync"
	"time"
)

var (
	ErrServiceBusy              = errors.New("service is busy")
	ErrServiceHasNotBeenStarted = errors.New("service has not been started")
	ErrServiceRunning           = errors.New("service is running")
)

type Status int

var statusNames = []string{
	"Inactive",
	"Active",
	"Starting",
	"Stopping",
	"WaitingActive",
	"Abnormal",
}

func (s Status) String() string {
	return statusNames[s]
}

const (
	StatusInactive Status = iota
	StatusActive
	StatusStarting
	StatusStopping
	StatusWaitingActive
	StatusAbnormal
)

type Interface interface {
	PrepareProcess(ctx context.Context) (map[string]*subprocess.SubProcess, error)
	WaitActive(ctx context.Context) bool
	GracefulShutdown(ctx context.Context, process map[string]*subprocess.SubProcess)
}

func InterruptGracefulShutdown(ctx context.Context, process map[string]*subprocess.SubProcess) {
	for _, p := range process {
		err := p.Interrupt()
		if err != nil {
			log.ErrorC(ctx, "interrupt failed", "process", p.Name(), "err", err)
		}
	}
}

type serviceRuntime struct {
	process          map[string]*subprocess.SubProcess
	processCtx       context.Context
	processCancel    context.CancelFunc
	waitActiveCtx    context.Context
	waitActiveCancel context.CancelFunc
}

func newServiceRuntime(process map[string]*subprocess.SubProcess) *serviceRuntime {
	processCtx, processCancel := context.WithCancel(context.Background())
	waitActiveCtx, waitActiveCancel := context.WithCancel(context.Background())
	return &serviceRuntime{
		process:          process,
		processCtx:       processCtx,
		processCancel:    processCancel,
		waitActiveCtx:    waitActiveCtx,
		waitActiveCancel: waitActiveCancel,
	}
}

type Service struct {
	instance Interface
	logger   *log.Logger

	runtime  *serviceRuntime
	status   Status
	busyLock sync.Mutex
}

func New(ins Interface, logger *log.Logger) *Service {
	return &Service{instance: ins, logger: logger}
}

func (s *Service) Status() Status {
	return s.status
}

func (s *Service) Running() bool {
	return s.runtime != nil
}

func (s *Service) Control(f func(process map[string]*subprocess.SubProcess) error) error {
	if !s.busyLock.TryLock() {
		return ErrServiceBusy
	}
	defer s.busyLock.Unlock()
	if s.runtime == nil {
		return ErrServiceHasNotBeenStarted
	}
	return f(s.runtime.process)
}

func (s *Service) Start(ctx context.Context) error {
	if !s.busyLock.TryLock() {
		return ErrServiceBusy
	}
	defer s.busyLock.Unlock()
	return s.start(ctx)
}

func (s *Service) Stop(ctx context.Context) error {
	if !s.busyLock.TryLock() {
		return ErrServiceBusy
	}
	defer s.busyLock.Unlock()
	s.stop(ctx)
	return nil
}

func (s *Service) Restart(ctx context.Context) error {
	if !s.busyLock.TryLock() {
		return ErrServiceBusy
	}
	defer s.busyLock.Unlock()
	return s.restart(ctx)
}

func (s *Service) start(ctx context.Context) error {
	if s.runtime != nil {
		return ErrServiceRunning
	}

	s.status = StatusStarting
	ctx = log.AppendCtx(ctx, "tag", "starting")
	process, err := s.instance.PrepareProcess(ctx)
	if err != nil {
		return errutil.Wrap(err)
	}

	svcRt := newServiceRuntime(process)
	err = s.startProcess(ctx, svcRt)
	if err != nil {
		return errutil.Wrap(err)
	}

	s.status = StatusWaitingActive
	go s.waitActive(svcRt.waitActiveCtx)

	s.runtime = svcRt
	return nil
}

// waitActive 使用上下文 serviceRuntime.waitActiveCtx
// 如何确保没有线程冲突？
// 1. waitActive 启动后，在正常退出前只可能遇到 stopProcess 线程
// 2. stopProcess 会先调用 serviceRuntime.waitActiveCancel 关闭 waitActive 线程然后再去停进程
func (s *Service) waitActive(waitActiveCtx context.Context) {
	ok := s.instance.WaitActive(waitActiveCtx)
	if errors.Is(waitActiveCtx.Err(), context.Canceled) {
		s.logger.Warn("context canceled, waitActive exited")
		return
	}

	if !ok {
		s.status = StatusAbnormal
		s.logger.Warn("set status to abnormal, waitActive exited")
		return
	}

	s.status = StatusActive
	s.logger.Warn("set status to active, waitActive exited")
}

func (s *Service) stop(ctx context.Context) {
	if s.runtime == nil {
		s.logger.DebugC(ctx, "process has not started yet, no need stop")
		return
	}

	s.status = StatusStopping
	ctx = log.AppendCtx(ctx, "tag", "stopping")
	s.stopProcess(ctx, s.runtime)
	s.status = StatusInactive
	s.runtime = nil
}

func (s *Service) restart(ctx context.Context) error {
	ctx = log.AppendCtx(ctx, "tag", "restarting")
	s.logger.InfoC(ctx, "begin restart")
	s.stop(ctx)
	err := s.start(ctx)
	if err != nil {
		s.logger.ErrorC(ctx, "restart failed at starting", "err", err)
		return errutil.Wrap(err)
	}
	s.logger.InfoC(ctx, "restart success")
	return nil
}

func (s *Service) startProcess(ctx context.Context, svcRt *serviceRuntime) error {
	s.logger.InfoC(ctx, "begin start process")
	for name, p := range svcRt.process {
		err := p.Start(svcRt.processCtx)
		if err != nil {
			s.logger.ErrorC(ctx, "start failed, begin kill all processes", "process", name, "err", err)
			s.stopProcess(ctx, svcRt)
			return errutil.Wrap(err)
		}
	}
	s.logger.InfoC(ctx, "start process success")
	return nil
}

func (s *Service) stopProcess(ctx context.Context, svcRt *serviceRuntime) {

	s.logger.InfoC(ctx, "begin graceful shutdown process")
	s.instance.GracefulShutdown(ctx, svcRt.process)

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer timeoutCancel()
	for _, p := range svcRt.process {
		if p.Ctx() != nil {
			select {
			case <-p.Ctx().Done():
			case <-timeoutCtx.Done():
			}
		}
	}
	if timeoutCtx.Err() == nil {
		s.logger.InfoC(ctx, "graceful shutdown process success")
		return
	}

	s.logger.ErrorC(ctx, "graceful shutdown process timeout, begin kill")
	svcRt.processCancel()
	return
}

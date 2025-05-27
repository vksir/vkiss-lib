package service

import (
	"context"
	"errors"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/subprocess"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServiceBusy = errors.New("service is busy")
)

const (
	StatusInactive = iota
	StatusActive
	StatusStarting
	StatusStopping
	StatusUpdating
)

var (
	statusToString = map[int32]string{
		StatusInactive: "Inactive",
		StatusActive:   "Active",
		StatusStarting: "Starting",
		StatusStopping: "Stopping",
		StatusUpdating: "Updating",
	}
)

type Interface interface {
	PrepareProcess(ctx context.Context) (process map[string]*subprocess.SubProcess, err error)
}

type GracefulShutdown func(ctx context.Context, logger *log.Logger, process map[string]*subprocess.SubProcess)

type Service struct {
	core             Interface
	gracefulShutdown GracefulShutdown
	logger           *log.Logger

	process map[string]*subprocess.SubProcess
	status  atomic.Int32
	free    sync.Mutex
}

func (s *Service) Status() string {
	return statusToString[s.status.Load()]
}

func (s *Service) SetActive() {
	s.status.Store(StatusActive)
}

func (s *Service) Process() map[string]*subprocess.SubProcess {
	return s.process
}

func (s *Service) Start(ctx context.Context) error {
	ok := s.free.TryLock()
	if !ok {
		return ErrServiceBusy
	}
	defer s.free.Unlock()
	s.status.Store(StatusStarting)

	log.AppendCtx(ctx, "tag", "starting")
	s.logger.InfoC(ctx, "begin start")
	process, err := s.core.PrepareProcess(ctx)
	if err != nil {
		return errutil.Wrap(err)
	}
	for name, p := range process {
		err = p.Start(ctx)
		if err != nil {
			s.logger.ErrorC(ctx, "start failed, begin kill all processes", "process", name, "err", err)
			if err := s.Stop(ctx); err != nil {
				s.logger.ErrorC(ctx, "stop failed", "err", err)
			}
			return errutil.Wrap(err)
		}
	}
	s.process = process
	s.logger.InfoC(ctx, "start success")
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	ok := s.free.TryLock()
	if !ok {
		return ErrServiceBusy
	}
	defer s.free.Unlock()
	s.status.Store(StatusStopping)
	s.status.Store(StatusInactive)

	log.AppendCtx(ctx, "tag", "stopping")
	if s.process == nil {
		s.logger.DebugC(ctx, "process has not started yet, no need stop")
		return nil
	}
	s.logger.InfoC(ctx, "begin graceful shutdown process")
	s.gracefulShutdown(ctx, s.logger, s.process)

	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	for _, p := range s.process {
		select {
		case <-p.Ctx().Done():
		case <-ctx.Done():
		}
	}

	if err := ctx.Err(); !errors.Is(err, context.DeadlineExceeded) {
		s.process = nil
		s.logger.InfoC(ctx, "graceful shutdown process success", "err", err)
		return nil
	}

	s.logger.InfoC(ctx, "graceful shutdown process timeout, begin kill")
	for _, p := range s.process {
		if err := p.Kill(); err != nil {
			s.logger.ErrorC(ctx, "kill failed", "process", p.Name())
		}
	}
	s.process = nil
	return errutil.Wrap(subprocess.ErrGraceFulShutdownTimeout)
}

func (s *Service) Restart(ctx context.Context) error {
	ok := s.free.TryLock()
	if !ok {
		return ErrServiceBusy
	}
	defer s.free.Unlock()

	log.AppendCtx(ctx, "tag", "restarting")
	s.logger.InfoC(ctx, "begin restart")
	err := s.Stop(ctx)
	if err != nil {
		s.logger.ErrorC(ctx, "restart failed at stopping", "err", err)
		return errutil.Wrap(err)
	}
	err = s.Start(ctx)
	if err != nil {
		s.logger.ErrorC(ctx, "restart failed at starting", "err", err)
		return errutil.Wrap(err)
	}
	s.logger.InfoC(ctx, "restart success")
	return nil
}

type Option interface {
	apply(s *Service)
}

func New(core Interface, opts ...Option) *Service {
	s := &Service{
		core:             core,
		gracefulShutdown: defaultGracefulShutdown,
		logger:           log.DefaultLogger(),
	}
	for _, opt := range opts {
		opt.apply(s)
	}
	return s
}

type optionFunc func(s *Service)

func (f optionFunc) apply(s *Service) {
	f(s)
}

func SetGracefulShutdown(f GracefulShutdown) Option {
	return optionFunc(func(s *Service) {
		s.gracefulShutdown = f
	})
}

func defaultGracefulShutdown(ctx context.Context, logger *log.Logger, process map[string]*subprocess.SubProcess) {
	for _, p := range process {
		err := p.Interrupt()
		if err != nil {
			logger.ErrorC(ctx, "interrupt failed", "process", p.Name(), "err", err)
		}
	}
}

package util

import "context"

type Runtime struct {
	parentCtx context.Context
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewRuntime(ctx ...context.Context) *Runtime {
	parentCtx := context.Background()
	if len(ctx) > 0 && ctx[0] != nil {
		parentCtx = ctx[0]
	}

	c, cancel := context.WithCancel(parentCtx)
	return &Runtime{parentCtx: parentCtx, ctx: c, cancel: cancel}
}

func (r *Runtime) Reset() {
	if r != nil {
		r.cancel()
		c, cancel := context.WithCancel(r.parentCtx)
		r.ctx = c
		r.cancel = cancel
	}
}

func (r *Runtime) Ctx() context.Context {
	if r != nil {
		return r.ctx
	}
	return nil
}

func (r *Runtime) Active() bool {
	if r != nil {
		return r.ctx.Err() == nil
	}
	return false
}

func (r *Runtime) Cancel() {
	if r != nil {
		r.cancel()
	}
}

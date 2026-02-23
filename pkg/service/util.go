package service

import (
	"context"
	"github.com/vksir/vkiss-lib/pkg/log"
)

func InterruptGracefulShutdown(ctx context.Context, process map[string]*SubProcess) {
	for _, p := range process {
		err := p.Interrupt()
		if err != nil {
			log.ErrorC(ctx, "interrupt failed", "process", p.Name(), "err", err)
		}
	}
}

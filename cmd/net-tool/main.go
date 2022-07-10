package main

import (
	"context"
	"flag"
	"net-tool/internal/udpredirect"
	"os"
	"os/signal"
)

var (
	devIP  = flag.String("d", "192.168.1.106", "Device IP.")
	filter = flag.String("f", "port 14001", "BPF filter expression.")
	srcIP  = flag.String("src-ip", "10.0.1.3", "Redirect src IP.")
	dstIP  = flag.String("dst-ip", "10.0.1.2", "Redirect dst IP.")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	flag.Parse()
	udpredirect.Run(ctx, *devIP, *filter, *srcIP, *dstIP)
}

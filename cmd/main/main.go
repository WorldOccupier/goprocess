package main

import (
	"context"
	"goprocess/metadata"
	"goprocess/processor"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	metadata.StartUpdater(ctx, &wg)

	processor.Process(ctx)

	<-ctx.Done()
	wg.Wait()
}

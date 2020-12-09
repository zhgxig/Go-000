package main

import (
	"context"
	"errors"
	"example/service"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	g := errgroup.Group{}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sm := http.NewServeMux()
	sm.HandleFunc("/service", service.Service)

	httpServerProduct := http.Server{
		Handler: sm,
		Addr:    ":8080",
	}

	httpServerDebug := http.Server{
		Handler: sm,
		Addr:    ":8090",
	}

	g.Go(func() error {
		fmt.Printf("Product Env Listen on %v\n", httpServerProduct.Addr)
		return httpServerProduct.ListenAndServe()
	})

	g.Go(func() error {
		fmt.Printf("Debug Env Listen on %v\n", httpServerDebug.Addr)
		return httpServerDebug.ListenAndServe()
	})

	g.Go(func() error {
		var msg string
		select {
		case s := <-signalCh:
			msg = fmt.Sprintf("Got signal: %s", s)
		case <-ctx.Done():
			msg = fmt.Sprintf("ctx err: %v", ctx.Err())
		}

		if err := httpServerProduct.Shutdown(ctx); err != nil {
			return err
		}
		if err := httpServerDebug.Shutdown(ctx); err != nil {
			return err
		}
		return errors.New(msg)
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("error: %s", err)
	}
}

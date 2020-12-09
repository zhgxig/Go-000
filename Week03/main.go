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

	signalMsg := make(chan os.Signal, 1)
	signal.Notify(signalMsg, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nsm1 := http.NewServeMux()
	nsm1.HandleFunc("/service1", service.Service1)

	nsm2 := http.NewServeMux()
	nsm2.HandleFunc("/service2", service.Service2)

	Server1 := http.Server{
		Handler: nsm1,
		Addr:    ":15000",
	}

	Server2 := http.Server{
		Handler: nsm2,
		Addr:    ":16000",
	}

	g.Go(func() error {
		fmt.Printf("Server1 addr %v\n", Server1.Addr)
		return Server1.ListenAndServe()
	})

	g.Go(func() error {
		fmt.Printf("Server2 addr %v\n", Server2.Addr)
		return Server2.ListenAndServe()
	})

	g.Go(func() error {
		var msg string
		select {
		case s := <-signalMsg:
			msg = fmt.Sprintf("signal msg: %s", s)
		case <-ctx.Done():
			msg = fmt.Sprintf("ctx: %v", ctx.Err())
		}

		if err := Server1.Shutdown(ctx); err != nil {
			return err
		}
		if err := Server2.Shutdown(ctx); err != nil {
			return err
		}
		return errors.New(msg)
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("err is : %s", err)
	}
}

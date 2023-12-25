package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"martinstromberg.se/gabriel/internal/server"
)

func main() {
    s, err := server.CreateSmtpServer("0.0.0.0")
    if err != nil {
        fmt.Println("Gabriel was unable to start:", err.Error())
        os.Exit(1)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <- sigCh
        cancel()
        fmt.Println("Received termination signal. Stopping...")
    }()

    err = s.Start(ctx, &wg)
    if err != nil {
        fmt.Println("Gabriel was unable to start", err.Error())
        os.Exit(1)
    }

    wg.Wait()
    fmt.Println("All workers stopped. Exiting...")
}


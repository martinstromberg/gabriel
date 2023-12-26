package network

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Handler interface {
    HandleConnection(c *TcpClient)
}

type TcpListener struct {
    handler     Handler
    context     context.Context
    waitGroup   *sync.WaitGroup
    tlsConfig   *tls.Config
}

func (tl *TcpListener) StartListener(port int) error {
    l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    defer l.Close()

    tcpL, ok := l.(*net.TCPListener)
    if !ok {
        err = errors.New("Unable to create TCPListener")
        fmt.Println(err.Error())

        return err
    }

    fmt.Println("Listening on port", port)

    dur := 5 * time.Second

    for {
        select {
        case <- tl.context.Done():
            fmt.Println("Shutting down listener on port", port)
            tl.waitGroup.Done()
            return nil
        default:
            tcpL.SetDeadline(time.Now().Add(dur))
            conn, err := tcpL.Accept()
            if err != nil {
                if !errors.Is(err, os.ErrDeadlineExceeded) {
                    fmt.Println("Error accepting connection:", err.Error())
                }

                continue
            }

            tcpC, ok := conn.(*net.TCPConn)
            if !ok {
                conn.Close()
                continue
            }

            c := NewClient(tcpC, tl.tlsConfig)
            go tl.handler.HandleConnection(c)
        }
    }
}

func NewListener(h Handler, ctx context.Context, wg *sync.WaitGroup, tlsConf *tls.Config) *TcpListener {
    l := &TcpListener{
        handler:    h,
        context:    ctx,
        waitGroup:  wg,
        tlsConfig:  tlsConf,
    }

    return l
}

package network

import (
	"context"
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

func StartListener(p int, h Handler, c context.Context, wg *sync.WaitGroup) error {
    l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", p))
    if err != nil {
        return err
    }

    defer l.Close()

    tcpL, ok := l.(*net.TCPListener)
    if !ok {
        return errors.New("Unable to create TCPListener")
    }

    fmt.Println("Listening on port", p)

    dur := 5 * time.Second

    for {
        select {
        case <- c.Done():
            fmt.Println("Shutting down listener on port", p)
            wg.Done()
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

            go h.HandleConnection(CreateClient(tcpC))
        }
    }
}


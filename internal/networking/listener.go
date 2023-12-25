package networking

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Represents a function that will handle an incoming connection
type ConnectionHandler interface {
    HandleConnection(conn net.Conn)
}

/*
Starts a listener on the provided address with incoming connections being 
handled by the ConnectionHandler until cancellation is signalled
*/
func StartListener(a string, ch ConnectionHandler, c context.Context, wg *sync.WaitGroup) error {
    l, err := net.Listen("tcp", a)
    if err != nil {
        return err
    }

    defer l.Close()

    var (
        ok bool
        tcpl *net.TCPListener
    )

    if tcpl, ok = l.(*net.TCPListener); !ok {
        return nil
    }

    fmt.Println("Listening on", a)

    dur := 5 * time.Second

    for {
        select {
        case <- c.Done():
            fmt.Println("Shutting down listener on ", a)
            wg.Done()
            return nil
        default:
            tcpl.SetDeadline(time.Now().Add(dur))
            conn, err := tcpl.Accept()
            if err != nil {
                if !errors.Is(err, os.ErrDeadlineExceeded) {
                    fmt.Println("Error accepting connection: ", err.Error())
                }

                continue;
            }

            go ch.HandleConnection(conn)
        }
    }
}


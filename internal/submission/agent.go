package submission

import (
	"context"
	"fmt"
	"io"
	"sync"

	"martinstromberg.se/gabriel/internal/network"
)

type Agent struct {
}

func (a *Agent) HandleConnection(c *network.TcpClient){
    defer c.Close()

    client := CreateClient(c, a)

    err := client.PerformHandshake()
    if err != nil {
        fmt.Println("Handshake failed", err.Error())
        return
    }

    var (
        sender        string
        recipients  []string
        headers     []string
        body          string
    )

    for {
        cmd, err := client.ReadCommand()
        if err != nil {
            if err == io.EOF {
                fmt.Println("From: ", sender)

                for _, v := range recipients {
                    fmt.Println("To: ", v)
                }

                for _, v := range headers {
                    fmt.Printf("'%s'\r\n", v)
                }

                fmt.Println(body)
            }
            return
        }

        switch cmd.Verb() {
            case ClientMail:
                s := cmd.(*Sender)
                sender = s.Address()
                client.writeString(fmt.Sprintf("250 2.1.0 Sender <%s> OK\r\n", sender))
                break
                
            case ClientRecipient:
                r := cmd.(*Recipient)
                recipients = append(recipients, r.Address())
                client.writeString(fmt.Sprintf("250 2.1.5 Recipient <%s> OK\r\n", r.Address()))
                break

            case ClientData:
                d := cmd.(*Data)
                headers = d.Headers()
                body = d.Body()
                client.writeString("250 2.0.0 Message received and queued for delivery\r\n")
                break

            case ClientReset:
            case ClientNoOp:
            case ClientQuit:
            case ClientVerify:
            case ClientExpand:
            case ClientHelp:
            default:
                return
        }
    }
}

func (a *Agent) Start(ctx context.Context, wg *sync.WaitGroup) {
    wg.Add(1)
    go network.StartListener(1587, a, ctx, wg)
}

func CreateAgent() (*Agent, error) {
    sa := &Agent{}

    return sa, nil
}

package submission

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"sync"

	"martinstromberg.se/gabriel/internal/network"
)

type Agent struct {
    port            int
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
            } else {
                fmt.Println(err.Error())
            }

            return
        }

        if !client.Authenticated() && cmd.Verb() != ClientAuthentication {
            client.writeString("530 5.7.1 Authentication required\r\n")
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

            case ClientStartTls:
                client.writeString("220 Ready to start TLS\r\n")
                if err := client.EnableTLS(); err != nil {
                    fmt.Println("TLS Error: ", err.Error())
                    return
                }
                break

            case ClientAuthentication:
                a.handleAuthentication(client, cmd.(*Authentication))
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

func (a *Agent) handleAuthentication(c *Client, cmd *Authentication) {
    panic("Not Implemented")
}

func (a *Agent) Start(ctx context.Context, wg *sync.WaitGroup) {
    cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
    if err != nil {
        panic("Unable to load certificates for TLS")
    }

    tlsConf := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         "smtp.martinstromberg.se",
        Certificates:       []tls.Certificate{cert},
    }

    l := network.NewListener(
        a,
        ctx,
        wg,
        tlsConf,
    )

    wg.Add(1)
    go l.StartListener(a.port)
}

func CreateAgent(port int) (*Agent, error) {
    sa := &Agent{
        port:       port,
    }

    return sa, nil
}

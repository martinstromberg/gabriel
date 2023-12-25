package agents

import (
	"context"
	"fmt"
	"net"
	"sync"

	"martinstromberg.se/gabriel/internal/commands"
	"martinstromberg.se/gabriel/internal/networking"
	"martinstromberg.se/gabriel/internal/submissions"
)

const (
    SUBMISSIONS_AGENT_PORT = 1587 // Remove the 1 later ;D
)

type SubmissionAgent struct {

}

func (sa *SubmissionAgent) HandleConnection(c net.Conn) {
    fmt.Println("Accepting connection from", c.RemoteAddr().String())
    defer c.Close()

    sc := submissions.CreateSubmissionClient(c)
    sc.SendGreeting("smtp.martinstromberg.se")

    cmd, err := sc.ReadCommand()
    if err != nil {
        fmt.Println("Error reading command", err.Error())
    }

    if !commands.IsHello(cmd.Keyword()) {
        fmt.Printf("Expected Hello, found '%s'\n", cmd.Keyword())
        return
    }

    hello := cmd.(*commands.Hello)
    sc.Name = hello.Name()
    
    sc.SendCapabilities("smtp.martinstromberg.se")

    for {
        cmd, err = sc.ReadCommand()
        if err != nil {
            fmt.Println("Error:", err.Error())
            break;
        }

        switch cmd.Keyword() {
        case commands.ClientMail:
            mailCmd, ok := cmd.(*commands.MailFrom)
            if !ok {
                break;
            }

            sc.SetSender(mailCmd.Sender())
            submissions.WriteString(c, fmt.Sprintf(
                "250 2.1.0 Sender <%s> OK\r\n",
                mailCmd.Sender().Address,
            ))

            break;

        case commands.ClientRecipient:
            rcptCmd, ok := cmd.(*commands.Recipient)
            if !ok {
                break;
            }

            sc.AddRecipient(rcptCmd.Address())
            submissions.WriteString(c, fmt.Sprintf(
                "250 2.1.5 Recipient <%s> OK\r\n",
                rcptCmd.Address().Address,
            ))

            break;

        case commands.ClientData:
            _, ok := cmd.(*commands.Data)
            if !ok {
                break;
            }

            submissions.WriteString(
                c,
                "250 2.0.0 Message received and queued for delivery\r\n",
            )

            // TODO: set data
        default:
            return
        }
    }
}

func (sa *SubmissionAgent) Start(ctx context.Context, wg *sync.WaitGroup) error {
    address := fmt.Sprintf("0.0.0.0:%d", SUBMISSIONS_AGENT_PORT)

    wg.Add(1)
    go networking.StartListener(address, sa, ctx, wg)
    
    return nil
}

func CreateSubmissionAgent() (*SubmissionAgent, error) {
    sa := &SubmissionAgent{}

    return sa, nil
}

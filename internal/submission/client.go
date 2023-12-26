package submission

import (
	"errors"
	"fmt"
	"strings"

	"martinstromberg.se/gabriel/internal/network"
)

const requestMaxSize int = 31457280

type Client struct {
    agent               *Agent
    authenticated       bool
    tcp                 *network.TcpClient
}

func (c *Client) Authenticated() bool {
    return c.authenticated
}

func (c *Client) readLine() ([]byte, error) {
    return c.tcp.ReadBytesUntil([]byte("\r\n"), requestMaxSize)
}

func (c *Client) writeString(str string) error {
    fmt.Println("MSA:", str)

    return c.tcp.WriteBytes([]byte(str))
}

func (c *Client) sendGreeting() error {
    return c.writeString(
        fmt.Sprintf("220 %s ESMTP Gabriel\r\n", "smtp.martinstromberg.se"),
    )
}

func (c *Client) PerformHandshake() error {
    err := c.sendGreeting()
    if err != nil {
        return err
    }

    cmd, err := c.ReadCommand()
    if err != nil {
        return err
    }

    if cmd.Verb() != ClientHello && cmd.Verb() != ClientExtendedHello {
        return fmt.Errorf("Expected Hello, got '%s'", cmd.Verb())
    }

    h, ok := cmd.(*Hello)
    if !ok {
        return errors.New("Unable to type cast command to Hello")
    }

    return c.sendCapabilities(h.Name())
}

func (c *Client) sendCapabilities(clientName string) error {
    extensions := []string{
        "STARTTLS",
        "AUTH PLAIN",
        "AUTH=REQUIRED",
        fmt.Sprintf("SIZE %d", requestMaxSize),
    }

    var sb strings.Builder
    var host string = "smtp.martinstromberg.se"
    
    sb.WriteString(fmt.Sprintf(
        "250-%s Hello %s, pleased to meet you\r\n",
        host,
        clientName,
    ))

    for i, v := range extensions {
        prefix := "-"
        if i == len(extensions) - 1 {
            prefix = " "
        }

        sb.WriteString(fmt.Sprintf("250%s%s\r\n", prefix, v))
    }

    return c.writeString(sb.String())
}

func (c *Client) EnableTLS() error {
    if err := c.tcp.EnableTLS(); err != nil {
        return err
    }

    cmd, err := c.ReadCommand()
    if err != nil {
        return err
    }

    h, ok := cmd.(*Hello)
    if !ok {
        return fmt.Errorf("Expected HELO/EHLO, got '%s'", cmd.Verb())
    }

    return c.sendCapabilities(h.Name())
}

func (c *Client) ReadCommand() (Command, error) {
    lb, err := c.readLine()
    if err != nil {
        return nil, err
    }

    fmt.Printf("MUA: %s", string(lb))

    verb := strings.Trim(strings.SplitN(string(lb), " ", 2)[0], "\r\n ")
    switch verb {
    case ClientHello:
    case ClientExtendedHello:
        return parseHelloFromBytes(lb)

    case ClientAuthentication:
        return parseAuthenticationFromBytes(lb)

    case ClientStartTls:
        return &StartTls{}, nil

    case ClientMail:
        return parseSenderFromBytes(lb)

    case ClientRecipient:
        return parseRecipientFromBytes(lb)

    case ClientData:
        return c.handleDataCommand()

    case ClientReset:
        return nil, errors.New("Not Implemented")
    case ClientNoOp:
        return nil, errors.New("Not Implemented")
    case ClientQuit:
        return nil, errors.New("Not Implemented")
    case ClientVerify:
        return nil, errors.New("Not Implemented")
    case ClientExpand:
        return nil, errors.New("Not Implemented")
    case ClientHelp:
        return nil, errors.New("Not Implemented")
    }

    return nil, fmt.Errorf("'%s' is not an accepted verb", verb)
}

func (c *Client) handleDataCommand() (*Data, error) {
    err := c.writeString(
        "354 Enter message, ending with \".\" on a line by itself\r\n",
    )

    if err != nil {
        return nil, err
    }

    buf, err := c.tcp.ReadBytesUntil([]byte("\r\n.\r\n"), requestMaxSize)
    if err != nil {
        return nil, err
    }

    // fmt.Printf("%s", string(buf))

    return parseDataFromBytes(buf)
}


func CreateClient(t *network.TcpClient, a *Agent) *Client {
    return &Client{
        tcp:    t,
        agent:  a,
    }
}

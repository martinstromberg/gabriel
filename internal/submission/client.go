package submission

import (
	"errors"
	"fmt"
	"strings"

	"martinstromberg.se/gabriel/internal/network"
)

const requestMaxSize int = 31457280

type Client struct {
    agent   *Agent
    tcp     *network.TcpClient
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
    var sb strings.Builder
    var host string = "smtp.martinstromberg.se"
    
    sb.WriteString(fmt.Sprintf(
        "250-%s Hello %s, pleased to meet you\r\n",
        host,
        clientName,
    ))

    sb.WriteString("250-STARTTLS\r\n")
    sb.WriteString("250-AUTH LOGIN PLAIN\r\n")
    sb.WriteString(fmt.Sprintf("250 SIZE %d\r\n", requestMaxSize))

    return c.writeString(sb.String())
}

func (c *Client) ReadCommand() (Command, error) {
    lb, err := c.readLine()
    if err != nil {
        return nil, err
    }

    fmt.Printf("MUA: %s", string(lb))

    verb := string(lb[:4])
    switch verb {
    case ClientHello:
    case ClientExtendedHello:
        return parseHelloFromBytes(lb)

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

    fmt.Printf("%s", string(buf))

    return parseDataFromBytes(buf)
}


func CreateClient(t *network.TcpClient, a *Agent) *Client {
    return &Client{
        tcp:    t,
        agent:  a,
    }
}

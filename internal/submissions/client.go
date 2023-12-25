package submissions

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"strings"

	"martinstromberg.se/gabriel/internal/commands"
)

type SubmissionClient struct {
    conn            net.Conn
    Name            string
    currentMail     Mail
}

func writeBytes(c net.Conn, buf []byte) error {
    _, err := c.Write(buf)

    return err
}

func WriteString(c net.Conn, str string) error {
    if err := writeBytes(c, []byte(str)); err != nil {
        return err
    }

    fmt.Printf("MSA: %s", str)

    return nil
}

func (sc *SubmissionClient) SendCommand(cmd commands.ServerCommand) error {
    buf := cmd.Bytes()
    if len(buf) == 0 {
        return errors.New("Provided command did not produce any bytes")
    }

    if err := writeBytes(sc.conn, buf); err != nil {
        return err
    }

    fmt.Println("MSA:", string(buf))

    return nil
}

func readLine(c net.Conn) ([]byte, error) {
    return readUntil(c, []byte("\r\n"))
}

func readUntil(c net.Conn, cutoff []byte) ([]byte, error) {
    var (
        buf bytes.Buffer
        bufSize int
    )

    for {
        tmp := make([]byte, 64)
        n, err := c.Read(tmp)
        if err != nil {
            if err == io.EOF {
                return nil, errors.New("CRLF not found before EOF")
            }

            return nil, fmt.Errorf("error while reading: %w", err)
        }

        bufSize += n
        buf.Write(tmp[:n])

        if bytes.HasSuffix(buf.Bytes(), cutoff) {
            break
        }
    }

    return buf.Bytes(), nil
}

func (sc *SubmissionClient) ResetTransaction() {
    sc.currentMail = Mail{}
}

func (sc *SubmissionClient) SetSender(address *mail.Address) {
    sc.currentMail.Sender = address
}

func (sc *SubmissionClient) AddRecipient(address *mail.Address) {
    sc.currentMail.Recipients = append(sc.currentMail.Recipients, address)
}

func (sc *SubmissionClient) SendGreeting(host string) error {
    greeting := commands.CreateGreeting(host, "Welcome friend!")

    return sc.SendCommand(greeting)
}

func (sc *SubmissionClient) SendCapabilities(host string) error {
    str := fmt.Sprintf(
        "250-%s Hello %s, pleased to meet you\r\n" +
        "250-STARTTLS\r\n" + 
        "250-AUTH LOGIN PLAIN\r\n" +
        "250 SIZE 31457280\r\n",
        host,
        sc.Name,
    )

    if err := writeBytes(sc.conn, []byte(str)); err != nil {
        return err
    }

    fmt.Println("MSA:", str)

    return nil
}

func (sc *SubmissionClient) ReadCommand() (commands.ClientCommand, error) {
    lineBytes, err := readLine(sc.conn)
    if err != nil {
        return nil, err
    }

    fmt.Println("MUA:", strings.TrimRight(string(lineBytes), "\r\n"))
    p := strings.SplitN(string(lineBytes), " ", 2)
    kw := strings.Trim(strings.ToUpper(p[0]), "\r\n ")

    switch kw {
    case commands.ClientHello:
    case commands.ClientExtendedHello:
        return commands.CreateHello(kw, strings.TrimRight(p[1], "\r\n "))

    case commands.ClientMail:
        return commands.CreateMailFromLine(strings.TrimRight(p[1], "\r\n "))

    case commands.ClientRecipient:
        return commands.ParseRecipient(strings.TrimRight(p[1], "\r\n "))

    case commands.ClientData:
        WriteString(
            sc.conn,
            "354 Enter message, ending with \".\" on a line by itself\r\n",
        )

        data, err := readUntil(sc.conn, []byte("\r\n.\r\n"))
        if err != nil {
            return nil, err
        }

        fmt.Println("MUA:", string(data))
        return commands.ParseData(string(data)), nil
    }
    
    return nil, errors.New(fmt.Sprintf("'%s' is not a recognized command", kw))
}

func CreateSubmissionClient(c net.Conn) *SubmissionClient {
    sc := &SubmissionClient{
        conn:           c,
        Name:           "",
    }

    return sc
}


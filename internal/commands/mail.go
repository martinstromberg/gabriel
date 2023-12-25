package commands

import (
	"errors"
	"net/mail"
	"strings"
)

type MailFrom struct {
    sender          *mail.Address
}

func (m *MailFrom) Keyword() string {
    return ClientMail
}

func (m *MailFrom) Sender() *mail.Address {
    return m.sender
}

func CreateMailFromLine(line string) (*MailFrom, error) {
    if len(line) == 0 {
        return nil, errors.New("Argument 'line' cannot be empty")
    }

    if strings.HasPrefix(line, ClientMail) {
        line, _ = strings.CutPrefix(line, ClientMail)
        line = strings.TrimLeft(line, " ")
    }

    parts := strings.SplitN(line, ":", 2)
    if len(parts) != 2 {
        return nil, errors.New("Invalid format for MAIL command")
    }

    addr, err := mail.ParseAddress(parts[1])
    if err != nil {
        return nil, err
    }

    m := &MailFrom{
        sender: addr,
    }

    return m, nil
}


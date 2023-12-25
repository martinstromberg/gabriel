package commands

import (
	"errors"
	"net/mail"
	"strings"
)

type Recipient struct {
    address         *mail.Address
}

func (r *Recipient) Address() *mail.Address {
    return r.address
}

func (r *Recipient) Keyword() string {
    return ClientRecipient
}

func ParseRecipient(l string) (*Recipient, error) {
    if strings.HasPrefix(l, ClientRecipient) {
        l, _ := strings.CutPrefix(l, ClientRecipient)
        l = strings.Trim(l, "\r\n ")
    }

    p := strings.SplitN(l, ":", 2)
    if len(p) != 2 {
        return nil, errors.New("Invalid recipient format")
    }

    a, err := mail.ParseAddress(p[1])
    if err != nil {
        return nil, err
    }

    r := &Recipient{
        address: a,
    }

    return r, nil
}

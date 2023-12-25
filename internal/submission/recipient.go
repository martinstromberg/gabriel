package submission

import (
    "errors"
    "fmt"
    "net/mail"
    "strings"
)

type Recipient struct {
    address     string
}

func (_ *Recipient) Verb() string {
    return ClientRecipient
}

func (r *Recipient) Address() string {
    return r.address
}

func parseRecipientFromBytes(buf []byte) (*Recipient, error) {
    if len(buf) < 13 { // should be shortest valid
        return nil, errors.New("Invalid format for Recipient command")
    }

    verb := strings.ToUpper(string(buf[:4]))
    if verb != ClientRecipient {
        return nil, fmt.Errorf("Got '%s', expected RCPT", verb)
    }

    str := strings.Trim(string(buf[8:]), "\r\n ")
    
    addr, err := mail.ParseAddress(str)
    if err != nil {
        return nil, err
    }

    r := &Recipient{
        address: addr.Address,
    }

    return r, nil
}

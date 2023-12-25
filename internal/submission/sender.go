package submission

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Sender struct {
    address         string
}

func (_ *Sender) Verb() string {
    return ClientMail
}

func (s *Sender) Address() string {
    return s.address
}

func parseSenderFromBytes(buf []byte) (*Sender, error) {
    if len(buf) < 15 { // should be shortest valid
        return nil, errors.New("Invalid format for Sender command")
    }

    verb := strings.ToUpper(string(buf[:4]))
    if verb != ClientMail {
        return nil, fmt.Errorf("Got '%s', expected MAIL", verb)
    }

    str := strings.Trim(string(buf[10:]), "\r\n ")
    
    addr, err := mail.ParseAddress(str)
    if err != nil {
        return nil, err
    }

    s := &Sender{
        address: addr.Address,
    }

    return s, nil
}

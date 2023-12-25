package submission

import (
	"errors"
	"fmt"
	"strings"
)

type Hello struct {
    extended        bool
    name            string
}

func (h *Hello) Extended() bool {
    return h.extended
}

func (h *Hello) Name() string {
    return h.name
}

func (h *Hello) Verb() string {
    if h.extended {
        return ClientExtendedHello
    }

    return ClientHello
}

func parseHelloFromBytes(buf []byte) (*Hello, error) {
    if len(buf) < 6 {
        return nil, errors.New("Not enough bytes for valid Hello")
    }

    verb := strings.ToUpper(string(buf[:4]))
    if verb != ClientHello && verb != ClientExtendedHello {
        return nil, fmt.Errorf("Verb '%s' is not valid for Hello", verb)
    }

    n := strings.Trim(string(buf[5:]), "\r\n ")
    h := &Hello{
        extended:   verb == ClientExtendedHello,
        name:       n,
    }

    return h, nil
}


package commands

import (
	"errors"
	"fmt"
)

type Hello struct {
    keyword         string
    name            string
}

func IsHello(kw string) bool {
    return kw == ClientExtendedHello || kw == ClientHello
}

func (h *Hello) Name() string {
    return h.name;
}

func (h *Hello) Keyword() string {
    return h.keyword
}

func CreateHello(kw string, name string) (*Hello, error) {
    if !IsHello(kw) {
        return nil, errors.New(
            fmt.Sprintf(
                "Keyword '%s' is invalid for Hello",
                kw,
            ),
        )
    }

    h := &Hello{
        keyword:    kw,
        name:       name,
    }
    
    return h, nil
}


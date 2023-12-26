package submission

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type Authentication struct {
    mechanism           string
    initialResponse     string
}

func (_ *Authentication) Verb() string {
    return ClientAuthentication
}

func (a *Authentication) Mechanism() string {
    return a.mechanism
}

func (a *Authentication) InitialResponse() string {
    return a.initialResponse
}

func parseAuthenticationFromBytes(buf []byte) (*Authentication, error) {
    if len(buf) < 6 {
        return nil, errors.New("Not enough bytes for valid Authentication")
    }

    verb := strings.ToUpper(string(buf[:4]))
    if verb != ClientAuthentication {
        return nil, fmt.Errorf("Verb '%s' is not valid for Authentication", verb)
    }

    str := strings.Trim(string(buf[5:]), "\r\n ")
    p := strings.SplitN(str, " ", 2)

    var init string = ""

    if len(p[1]) > 0 && p[1] != "=" {
        bytes, err := base64.StdEncoding.DecodeString(p[1])
        if err != nil {
            return nil, err
        }

        init = string(bytes)
    }

    a := &Authentication{
        mechanism: p[0],
        initialResponse: init,
    }

    return a, nil
}

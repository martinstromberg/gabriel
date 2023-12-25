package commands

import "fmt"

type Greeting struct {
    host        string
    message     string
}

func (_ *Greeting) Keyword() string {
    return ServerGreeting
}

func (g *Greeting) Bytes() []byte {
    msg := fmt.Sprintf(
        "220 %s ESMTP %s\r\n",
        g.host,
        g.message,
    )

    return []byte(msg)
}

func CreateGreeting(host string, message string) *Greeting {
    g := &Greeting{
        host:       host,
        message:    message,
    }

    return g
}

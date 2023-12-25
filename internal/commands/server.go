package commands

const (
    ServerGreeting          string  = "220"
    ServerOkay                      = "250"
)

type ServerCommand interface {
    Bytes() []byte
}

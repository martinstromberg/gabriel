package commands

const (
    OkayCapabilities    string  = "OKCAP"
    OkayQueued                  = "2.0.0"
    OkaySender                  = "2.1.0"
    OkayRecipient               = "2.1.5"
    OkayNoContent               = "2.5.0"
)

type Okay struct {
    data        []string
    status        string
}

func (o *Okay) Bytes() []byte {
    switch o.status {
    case OkayCapabilities:

        return []byte("")

    default:
        return make([]byte, 0)
    }
}

func (o *Okay) Keyword() string {
    return ServerOkay
}

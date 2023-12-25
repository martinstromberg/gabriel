package submission

const (
    ClientData              string  = "DATA"
    ClientExpand                    = "EXPN"
    ClientExtendedHello             = "EHLO"
    ClientHello                     = "HELO"
    ClientHelp                      = "HELP"
    ClientMail                      = "MAIL"
    ClientNoOp                      = "NOOP"
    ClientQuit                      = "QUIT"
    ClientRecipient                 = "RCPT"
    ClientReset                     = "RSET"
    ClientVerify                    = "VRFY"
)

type Command interface {
    Verb() string
}


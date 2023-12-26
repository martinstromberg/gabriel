package submission

const (
    ClientAuthentication    string  = "AUTH"
    ClientData                      = "DATA"
    ClientExpand                    = "EXPN"
    ClientExtendedHello             = "EHLO"
    ClientHello                     = "HELO"
    ClientHelp                      = "HELP"
    ClientMail                      = "MAIL"
    ClientNoOp                      = "NOOP"
    ClientQuit                      = "QUIT"
    ClientRecipient                 = "RCPT"
    ClientReset                     = "RSET"
    ClientStartTls                  = "STARTTLS"
    ClientVerify                    = "VRFY"
)

type Command interface {
    Verb() string
}


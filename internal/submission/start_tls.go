package submission

type StartTls struct {}

func (_ *StartTls) Verb() string {
    return ClientStartTls
}


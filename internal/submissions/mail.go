package submissions

import "net/mail"

type Mail struct {
    Sender                *mail.Address
    Recipients          []*mail.Address
    Content               *mail.Message
}



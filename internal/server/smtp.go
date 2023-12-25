package server

import (
    "context"
    "sync"

    "martinstromberg.se/gabriel/internal/submission"
)

type SmtpServer struct {
    msa         *submission.Agent
}

func (s *SmtpServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
    s.msa.Start(ctx, wg)

    return nil
}

func CreateSmtpServer(address string) (*SmtpServer, error) {
    sa, err := submission.CreateAgent()
    if err != nil {
        return nil, err
    }

    s := &SmtpServer{
        msa:        sa,
    }

    return s, nil
}


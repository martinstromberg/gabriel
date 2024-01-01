package server

import (
    "context"
    "sync"

    "martinstromberg.se/gabriel/internal/submission"
)

type SmtpServer struct {
    config      *Config
    msa         *submission.Agent
}

func (s *SmtpServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
    s.msa.Start(ctx, wg)

    return nil
}

func CreateSmtpServer(c *Config) (*SmtpServer, error) {
    var (
        err error
        sa *submission.Agent = nil
    )

    if c.SubmissionsAgent {
        sa, err = submission.CreateAgent(c.ServerPort)
        if err != nil {
            return nil, err
        }
    }

    s := &SmtpServer{
        config:     c,
        msa:        sa,
    }

    return s, nil
}


package server

import (
	"context"
	"net"
	"sync"

	"martinstromberg.se/gabriel/internal/agents"
)

type SmtpServer struct {
    msa         *agents.SubmissionAgent
}

func (s *SmtpServer) HandleConnection (c net.Conn) {

}

func (s *SmtpServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
    err := s.msa.Start(ctx, wg)
    if err != nil {
        return err
    }

    return nil
}

func CreateSmtpServer(address string) (*SmtpServer, error) {
    sa, err := agents.CreateSubmissionAgent()
    if err != nil {
        return nil, err
    }

    s := &SmtpServer{
        msa:        sa,
    }

    return s, nil
}

package network

import (
	"bytes"
    "crypto/tls"
	"fmt"
	"io"
	"net"
)

const (
    readBufferSize      int     = 64
)

type TcpClient struct {
    conn            net.Conn
    tlsConfig       *tls.Config
}

func (c *TcpClient) Close() error {
    return c.conn.Close()
}

func (c *TcpClient) EnableTLS() error {
    tlsConn := tls.Server(c.conn, c.tlsConfig)
    if err := tlsConn.Handshake(); err != nil {
        return err
    }

    c.conn = tlsConn

    return nil
}

func (c *TcpClient) ReadBytesUntil(p []byte, ms int) ([]byte, error) {
    var (
        buf         bytes.Buffer
        bufSize     int
    )

    for {
        tmp := make([]byte, readBufferSize)
        n, err := c.conn.Read(tmp)
        if err != nil {
            if err == io.EOF {
                return nil, err
            }

            return nil, fmt.Errorf("Error while rading: %w", err)
        }

        bufSize += n
        buf.Write(tmp[:n])

        if bytes.HasSuffix(buf.Bytes(), p) {
            break;
        }

        if bufSize >= ms {
            err = fmt.Errorf("Buffer exceeds maximum allowed size of %d", ms)
            return nil, err
        }
    }

    return buf.Bytes(), nil
}

func (c *TcpClient) WriteBytes(buf []byte) error {
    _, err := c.conn.Write(buf)

    return err
}

func NewClient(c *net.TCPConn, tlsConf *tls.Config) *TcpClient {
    return &TcpClient{
        conn:       c,
        tlsConfig:  tlsConf,
    }
}

package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
    readBufferSize      int     = 64
)

type TcpClient struct {
    conn    *net.TCPConn
}

func (c *TcpClient) Close() error {
    return c.conn.Close()
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

func CreateClient(c *net.TCPConn) *TcpClient {
    return &TcpClient{
        conn: c,
    }
}

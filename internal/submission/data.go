package submission

import (
	"strings"
)

type Data struct {
    headers         []string
    body            string
}

func (d *Data) Body() string {
    return d.body
}

func (d *Data) Headers() []string {
    return d.headers
}

func (_ *Data) Verb() string {
    return ClientData
}

func parseDataFromBytes(buf []byte) (*Data, error) {
    p := strings.SplitN(string(buf[:len(buf) - 3]), "\r\n\r\n", 2)
    h := strings.Split(p[0], "\r\n")

    b := strings.Trim(p[1], "\r\n ")

    d := &Data{
        headers: h,
        body: b,
    }
    
    return d, nil
}

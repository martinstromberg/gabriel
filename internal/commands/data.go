package commands

type Data struct {
    data        string
}

func (d *Data) Keyword() string {
    return ClientData
}

func ParseData(data string) *Data {
    return &Data{
        data: data,
    }
}


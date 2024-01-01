package server

import (
    "gopkg.in/ini.v1"
)

const (
    AppModeDevelopment          string      = "development"
    AppModeProduction           string      = "production"
)

type Config struct {
    AppMode                     string

    CertFilePath                string
    KeyFilePath                 string

    ServerPort                  int
    ServerHost                  string
    ServerHostName              string

    SubmissionsAgent            bool
}

func (c *Config) Development() bool {
    return c.AppMode == AppModeDevelopment
}

func parseConfig(cfg *ini.File) (*Config, error) {
    tlsSection := cfg.Section("tls")
    serverSection := cfg.Section("server")
    submissionSection := cfg.Section("submission")

    port, err := serverSection.Key("port").Int()
    if err != nil {
        return nil, err
    }

    msaEnabled, err := submissionSection.Key("enabled").Bool()
    if err != nil {
        return nil, err
    }

    c := &Config{
        AppMode: cfg.Section("").Key("app_mode").In(
            AppModeProduction,
            []string{
                AppModeDevelopment,
                AppModeProduction,
            },
        ),

        CertFilePath: tlsSection.Key("cert_file_path").String(),
        KeyFilePath: tlsSection.Key("key_file_path").String(),

        ServerHostName: serverSection.Key("host_name").String(),
        ServerHost: serverSection.Key("host").String(),
        ServerPort: port,

        SubmissionsAgent: msaEnabled,
    }

    return c, nil
}

func LoadConfig(path string) (*Config, error) {
    cfg, err := ini.Load(path)
    if err != nil {
        return nil, err
    }

    return parseConfig(cfg)
}

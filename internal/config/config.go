package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Server        string `yaml:"server"`
		ListenAddress string `yaml:"listen_address"`
		ListenPort    int    `yaml:"listen_port"`
		CaCertFile    string `yaml:"ca_cert_file"`
		CertFile      string `yaml:"cert_file"`
		KeyFile       string `yaml:"key_file"`
	}

	Db struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		DbName   string `yaml:"database"`
		UserName string `yaml:"user"`
		Password string `yaml:"password"`
	}

	Ldap struct {
		ServerURI  string `yaml:"server_uri"`
		SearchBase string `yaml:"search_base"`
		BindDN     string `yaml:"bind_dn"`
		BindPW     string `yaml:"bind_pw"`
		StartTLS   bool   `yaml:"start_tls"`
	}
}

func Load(file string) (*Config, error) {
	// load the config
	fh, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer fh.Close()

	decoder := yaml.NewDecoder(fh)

	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// set the file paths - relative to the config file location
	configdir := filepath.Dir(file)

	if !strings.HasPrefix(cfg.Service.CaCertFile, "/") {
		cfg.Service.CaCertFile = filepath.Join(configdir, cfg.Service.CaCertFile)
	}
	if !strings.HasPrefix(cfg.Service.CertFile, "/") {
		cfg.Service.CertFile = filepath.Join(configdir, cfg.Service.CertFile)
	}
	if !strings.HasPrefix(cfg.Service.KeyFile, "/") {
		cfg.Service.KeyFile = filepath.Join(configdir, cfg.Service.KeyFile)
	}

	return &cfg, nil
}

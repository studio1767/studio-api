package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Id               string `yaml:"id"`
		Secret           string `yaml:"secret"`
		StateSecret      string `yaml:"state_secret"`
		StateKey         []byte
		CookieHashSecret string `yaml:"cookie_hash_secret"`
		CookieHashKey    []byte
		CookieEncSecret  string `yaml:"cookie_enc_secret"`
		CookieEncKey     []byte
		RedirectURLs     []string `yaml:"redirect_urls"`
	}

	Listener string `yaml:"listener"`

	Https struct {
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	}

	Idp struct {
		IssuerURL  string `yaml:"issuer_url"`
		CaCertFile string `yaml:"ca_cert_file"`
	}

	Db struct {
		Server     string `yaml:"server"`
		Port       int    `yaml:"port"`
		DbName     string `yaml:"database"`
		UserName   string `yaml:"user"`
		Password   string `yaml:"password"`
		CaCertFile string `yaml:"ca_cert_file"`
	}
}

func Load(file string, noauth bool) (*Config, error) {
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

	// decode the state secret
	if noauth == false {
		cfg.Service.StateKey, err = decodeKey(cfg.Service.StateSecret, 32)
		if err != nil {
			return nil, err
		}
		cfg.Service.CookieHashKey, err = decodeKey(cfg.Service.CookieHashSecret, 32)
		if err != nil {
			return nil, err
		}
		cfg.Service.CookieEncKey, err = decodeKey(cfg.Service.CookieEncSecret, 16)
		if err != nil {
			return nil, err
		}
	}

	// set the file paths
	configdir := filepath.Dir(file)
	if len(cfg.Https.KeyFile) > 0 && !strings.HasPrefix(cfg.Https.KeyFile, "/") {
		cfg.Https.KeyFile = filepath.Join(configdir, cfg.Https.KeyFile)
	}
	if len(cfg.Https.CertFile) > 0 && !strings.HasPrefix(cfg.Https.CertFile, "/") {
		cfg.Https.CertFile = filepath.Join(configdir, cfg.Https.CertFile)
	}
	if len(cfg.Idp.CaCertFile) > 0 && !strings.HasPrefix(cfg.Idp.CaCertFile, "/") {
		cfg.Idp.CaCertFile = filepath.Join(configdir, cfg.Idp.CaCertFile)
	}
	if len(cfg.Db.CaCertFile) > 0 && !strings.HasPrefix(cfg.Db.CaCertFile, "/") {
		cfg.Db.CaCertFile = filepath.Join(configdir, cfg.Db.CaCertFile)
	}

	// standardize the URL
	cfg.Idp.IssuerURL = strings.TrimRight(cfg.Idp.IssuerURL, "/")

	return &cfg, nil
}

func decodeKey(kvalue string, length int) ([]byte, error) {
	// remove any trailing '=' characters
	kvalue = strings.TrimRight(kvalue, "=")

	// now decode
	kdata, err := base64.RawStdEncoding.DecodeString(kvalue)
	if err != nil {
		return nil, fmt.Errorf("config: failed to decode secret: %w", err)
	}
	if len(kdata) != length {
		return nil, fmt.Errorf("config: key length incorrect")
	}
	return kdata, nil
}

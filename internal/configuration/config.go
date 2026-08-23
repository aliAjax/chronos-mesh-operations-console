package configuration

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Source struct {
	Name     string  `yaml:"name"`
	Address  string  `yaml:"address"`
	OffsetMS float64 `yaml:"offset_ms"`
	DelayMS  float64 `yaml:"delay_ms"`
	Reach    uint8   `yaml:"reach"`
}
type Config struct {
	HTTPAddr      string          `yaml:"http_addr"`
	NTPAddr       string          `yaml:"ntp_addr"`
	NTSAddr       string          `yaml:"nts_addr"`
	MinStratum    uint8           `yaml:"min_stratum"`
	Holdover      time.Duration   `yaml:"holdover"`
	CookieKey     string          `yaml:"cookie_key"`
	Sources       []Source        `yaml:"sources"`
	RatePerSecond int             `yaml:"rate_per_second"`
	Policies      map[string]bool `yaml:"policies"`
}

func Default() Config {
	return Config{HTTPAddr: ":8080", NTPAddr: ":8123", NTSAddr: ":8443", MinStratum: 2, Holdover: time.Hour, CookieKey: "", RatePerSecond: 0, Sources: []Source{{Name: "primary", Address: "127.0.0.1:9123"}, {Name: "secondary", Address: "127.0.0.1:9124"}, {Name: "faulty", Address: "127.0.0.1:9125", OffsetMS: 5000}}}
}
func Load(path string) (Config, error) {
	c := Default()
	if path == "" {
		path = os.Getenv("CHRONOS_CONFIG")
	}
	if path == "" {
		return c, nil
	}
	b, e := os.ReadFile(path)
	if e != nil {
		return c, e
	}
	e = yaml.Unmarshal(b, &c)
	return c, e
}

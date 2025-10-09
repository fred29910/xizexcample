package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	OTLP    OTLPConfig    `yaml:"otlp"`
}

// MetricsConfig holds the metrics configuration.
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	Addr    string `yaml:"addr"`
}

// OTLPConfig holds the OpenTelemetry configuration.
type OTLPConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Endpoint    string  `yaml:"endpoint"`
	Protocol    string  `yaml:"protocol"`
	ServiceName string  `yaml:"service_name"`
	SampleRatio float64 `yaml:"sample_ratio"`
}

// Loader handles loading and watching the configuration file.
type Loader struct {
	Path string
	Cfg  *AppConfig
	mu   sync.RWMutex
	cb   func(*AppConfig)
}

// NewLoader creates a new configuration loader.
func NewLoader(path string, cb func(*AppConfig)) (*Loader, error) {
	ld := &Loader{Path: path, cb: cb}
	if err := ld.load(); err != nil {
		return nil, err
	}
	go ld.watch()
	return ld, nil
}

func (l *Loader) load() error {
	b, err := ioutil.ReadFile(l.Path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var c AppConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	l.mu.Lock()
	l.Cfg = &c
	l.mu.Unlock()

	if l.cb != nil {
		l.cb(&c)
	}
	return nil
}

func (l *Loader) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("watcher init error:", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(l.Path)
	if err != nil {
		log.Println("watcher add error:", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := l.load(); err != nil {
					log.Println("reload config error:", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("watcher error:", err)
		}
	}
}

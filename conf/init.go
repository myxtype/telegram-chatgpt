package conf

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	configPath string
	config     GbeConfig
	configOnce sync.Once
)

func init() {
	flag.StringVar(&configPath, "conf", "", "Config file path. This path must include config.toml file.")
}

func Config() *GbeConfig {
	configOnce.Do(func() {
		flag.Parse()

		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		if configPath == "" {
			if p, err := execPath(); err == nil {
				for _, n := range p {
					viper.AddConfigPath(n)
				}
			} else {
				panic(err)
			}
		} else {
			viper.AddConfigPath(configPath)
		}
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			if e.Op != fsnotify.Remove {
				if err := viper.Unmarshal(&config); err != nil {
					log.Println("Reload config error", err)
				}
			}
		})
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	})
	return &config
}

func execPath() (p []string, err error) {
	p = []string{"./"}
	if _, currentPath, _, ok := runtime.Caller(0); ok {
		p = append(p, filepath.Dir(currentPath))
	}
	if t, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		p = append(p, t)
	}
	return
}

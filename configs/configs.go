package configs

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/luo2pei4/base-server/internal/logger"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

type ServiceConfigs struct {
	SericePort string `toml:"service_port"` // 端口
	LogLevel   string `toml:"log_level"`    // 日志输出等级
	LogFile    string `toml:"log_file"`     // 日志文件路径
}

var (
	serviceConfigsMu   sync.RWMutex
	serviceConfigsFile string
	serviceConfigs     *ServiceConfigs
)

func init() {
	serviceConfigs = &ServiceConfigs{}
}

func LoadServiceConfig(path string) error {
	serviceConfigsMu.Lock()
	defer serviceConfigsMu.Unlock()
	if len(strings.TrimSpace(path)) == 0 {
		return nil
	}
	_, err := os.Lstat(path)
	if err != nil {
		return err
	}
	serviceConfigsFile = path

	file, err := os.ReadFile(serviceConfigsFile)
	if err != nil {
		return err
	}
	if err = toml.Unmarshal(file, serviceConfigs); err != nil {
		return err
	}
	return nil
}

func StartServiceConfigWatch() {
	v := viper.New()
	v.SetConfigFile(serviceConfigsFile)
	v.SetConfigType("toml")
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		LoadServiceConfig(serviceConfigsFile)
		logger.SetLogLevel(serviceConfigs.LogLevel)
	})
}

func GetSerivePort() (string, error) {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	if len(serviceConfigs.SericePort) == 0 {
		return ":8080", nil
	}
	matched, err := regexp.MatchString("^:?[0-9]{4,5}", serviceConfigs.SericePort)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("%s is an invalid port", serviceConfigs.SericePort)
	}
	return serviceConfigs.SericePort, nil
}

func GetLogLevel() string {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	if len(serviceConfigs.LogLevel) == 0 {
		return "info"
	}
	return serviceConfigs.LogLevel
}

func GetLogFile() string {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	if len(serviceConfigs.LogFile) == 0 {
		return "./base-server.log"
	}
	return serviceConfigs.LogFile
}

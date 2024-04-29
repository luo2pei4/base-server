package config

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/luo2pei4/base-server/internal/logger"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

type ServiceConfigs struct {
	ServicePort  string `toml:"service_port"`   // 端口
	LogLevel     string `toml:"log_level"`      // 日志输出等级
	LogFile      string `toml:"log_file"`       // 日志文件路径
	MaxCPURate   int    `toml:"max_cpu_rate"`   // cpu最大使用率
	I18nRootPath string `toml:"i18n_root_path"` // i18n文件根路径
	Language     string `toml:"language"`       // 语言
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
		// 重新加载配置文件
		LoadServiceConfig(serviceConfigsFile)
		// 设置最大CPU使用数量
		SetMaxCPUNum()
		// 设置日志等级
		logger.SetLogLevel(serviceConfigs.LogLevel)
	})
}

func GetSerivePort() (string, error) {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	if len(serviceConfigs.ServicePort) == 0 {
		return ":8080", nil
	}
	matched, err := regexp.MatchString("^:?[0-9]{4,5}", serviceConfigs.ServicePort)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("%s is an invalid port", serviceConfigs.ServicePort)
	}
	return serviceConfigs.ServicePort, nil
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

func GetMaxCPUUsage() int {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	if serviceConfigs.MaxCPURate > 100 {
		return 100
	}
	if serviceConfigs.MaxCPURate <= 0 {
		return 0
	}
	return serviceConfigs.MaxCPURate
}

func Geti18nDir() string {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	return serviceConfigs.I18nRootPath
}

func GetLanguage() string {
	serviceConfigsMu.RLock()
	defer serviceConfigsMu.RUnlock()
	return serviceConfigs.Language
}

func SetMaxCPUNum() {
	maxCPUNum := int((runtime.NumCPU() * serviceConfigs.MaxCPURate) / 100)
	if maxCPUNum < 1 {
		maxCPUNum = 1
	}
	runtime.GOMAXPROCS(maxCPUNum)
	logger.Infof("set max cpu num: %d", maxCPUNum)
}

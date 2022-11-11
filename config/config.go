package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dean2032/go-project-layout/utils"
	"github.com/dean2032/go-project-layout/utils/errors"
	"github.com/dean2032/go-project-layout/utils/logging"
	"go.uber.org/fx"
)

// Config ...
type Config struct {
	Debug      bool   `json:"debug"`
	PprofPath  string `json:"pprof_path"`
	ServerPort string `json:"server_port"`
	PublicDir  string `json:"public_dir"`
	LogDir     string `json:"log_dir"`
	MainDB     string `json:"main_db"`
	Redis      struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	JWTSecret            string `json:"jwt_secret"`
	DBConnectionPoolSize int    `json:"db_connection_pool_size"`
}

// Module ...
var Module = fx.Provide(GetConfig)

// cfg global config
var globalCfg = DefaultConfig()

// GetConfig get global config
func GetConfig() *Config {
	return globalCfg
}

// LoadFromPath ...
func LoadFromPath(configPath string) (*Config, error) {
	cfg := DefaultConfig()
	if !utils.IsFileExist(configPath) {
		return cfg, nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, errors.Wrap(err, fmt.Sprintf("read file %s fail: ", configPath))
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, fmt.Sprintf("load yaml file %s fail: ", configPath))
	}
	fixConfigDirs(cfg)
	globalCfg = cfg
	return cfg, nil
}

func fixConfigDirs(cfg *Config) {
	if cfg.PublicDir == "" {
		cfg.PublicDir = "."
	}

	var err error
	cfg.PublicDir, err = filepath.Abs(cfg.PublicDir)
	if err != nil {
		panic(fmt.Sprintf("get dir of %s fail: ", cfg.PublicDir))
	}
	if cfg.LogDir == "" {
		cfg.LogDir = "."
	}
	cfg.LogDir, err = filepath.Abs(cfg.LogDir)
	if err != nil {
		panic(fmt.Sprintf("get dir of %s fail: ", cfg.LogDir))
	}
}

// DefaultConfig set default config
func DefaultConfig() *Config {
	cfg := &Config{
		Debug:                false,
		ServerPort:           "8888",
		LogDir:               "./log",
		PublicDir:            ".",
		PprofPath:            "/debug/pprof",
		DBConnectionPoolSize: 1000,
		// db DSN example
		MainDB: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
	}
	fixConfigDirs(cfg)
	return cfg
}

// GenLoggerModule ...
func GenLoggerModule(name string) fx.Option {
	return fx.Invoke(func(cfg *Config) error {
		lg, err := logging.InitLogger(name, cfg.LogDir, cfg.Debug)
		if err != nil {
			return err
		}
		logging.ReplaceLogger(lg)
		return nil
	})
}

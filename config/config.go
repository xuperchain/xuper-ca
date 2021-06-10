/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package config

import (
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	DbConfig DbConfig `yaml:"dbConfig,omitempty"`
	CertPath string   `yaml:certPath,omitempty`
	Port     string   `yaml:"port,omitempty"`
	HttpPort string   `yaml:"httpPort,omitempty"`
	CaAdmin  string   `yaml:"caAdmin,omitempty"`
	Log      Log      `yaml:"log,omitempty"`
}

type DbConfig struct {
	DbType string `yaml:"dbType,omitempty"`
	DbPath string `yaml:"dbPath,omitempty"`
}

type Log struct {
	Level string `yaml:"level,omitempty"`
	Path  string `yaml:"path,omitempty"`
}

func InstallCaConfig(configFile string) error {
	// 从配置文件中加载配置
	config = &Config{}
	filePath, fileName := filepath.Split(configFile)
	file := strings.TrimSuffix(fileName, path.Ext(fileName))
	viper.AddConfigPath(filePath)
	viper.SetConfigName(file)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Read config file error!", "err", err.Error())
		return nil
	}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Unmarshal config from file error! error=", err.Error())
		return nil
	}

	printConfig()

	// 监听配置变化, 重启加载
	//viper.WatchConfig()
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	// 配置发生变化则重新加载
	//	config = &Config{}
	//	viper.Unmarshal(config)
	//	printConfig()
	//})

	return nil
}

func printConfig() {
	log.Infof("init config: %+v", config)
}

func GetDBConfig() *DbConfig {
	return &config.DbConfig
}

func GetCertPath() string {
	return config.CertPath
}

func GetServerPort() string {
	return config.Port
}

func GetHttpPort() string {
	return config.HttpPort
}

func GetCaAdmin() string {
	return config.CaAdmin
}

func GetLog() Log {
	return config.Log
}

/*
Copyright © 2020 Marvin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package service

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
)

// 程序配置文件
type CfgFile struct {
	SourceConfig SourceConfig `toml:"source" json:"source"`
	LogConfig    LogConfig    `toml:"log" json:"log"`
}

type SourceConfig struct {
	Username      string   `toml:"username" json:"username"`
	Password      string   `toml:"password" json:"password"`
	Host          string   `toml:"host" json:"host"`
	Port          int      `toml:"port" json:"port"`
	ServiceName   string   `toml:"service-name" json:"service-name"`
	ConnectParams string   `toml:"connect-params" json:"connect-params"`
	SessionParams []string `toml:"session-params" json:"session-params"`
	RootPwd       string   `toml:"rootpwd" json:"rootpwd"`
	PrimaryOracleHome string `toml:"primary-oracle-home" json:"primary-oracle-home"`
	PrimaryGridHome string `toml:"primary-grid-home" json:"primary-grid-home"`
	PrimaryOracleHomeOwner string `toml:"primary-oracle-home-owner" json:"primary-oracle-home-owner"`
	PrimaryGridHomeOwner string `toml:"primary-grid-home-owner" json:"primary-grid-home-owner"`
	StandbyOracleBase string `toml:"standby-oracle-base" json:"standby-oracle-base"`
	StandbyOracleHome string `toml:"standby-oracle-home" json:"standby-oracle-home"`
	StandbyGridHome string `toml:"standby-grid-home" json:"standby-grid-home"`
	StandbyOracleHomeOwner string `toml:"standby-oracle-home-owner" json:"standby-oracle-home-owner"`
	StandbyGridHomeOwner string `toml:"standby-grid-home-owner" json:"standby-grid-home-owner"`
	PrimaryDataDg string 	`toml:"primary-data-dg" json:"primary-data-dg"`
	StandbyDataDg string    `toml:"standby-data-dg" json:"standby-data-dg"`
	StandbyHostIps []string	`toml:"standby-host-ips"json:"standby-host-ips"`
	OracleDBname string
	OracleSid   string
	OracleUniqname string
	IsRAC 		string
}



type LogConfig struct {
	LogLevel   string `toml:"log-level" json:"log-level"`
	LogFile    string `toml:"log-file" json:"log-file"`
	MaxSize    int    `toml:"max-size" json:"max-size"`
	MaxDays    int    `toml:"max-days" json:"max-days"`
	MaxBackups int    `toml:"max-backups" json:"max-backups"`
}

// 读取配置文件
func ReadConfigFile(file string) (*CfgFile, error) {
	cfg := &CfgFile{}
	if err := cfg.configFromFile(file); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// 加载配置文件并解析
func (c *CfgFile) configFromFile(file string) error {
	if _, err := toml.DecodeFile(file, c); err != nil {
		return fmt.Errorf("failed decode toml config file %s: %v", file, err)
	}
	return nil
}



func (c *CfgFile) String() string {
	cfg, err := json.Marshal(c)
	if err != nil {
		return "<nil>"
	}
	return string(cfg)
}





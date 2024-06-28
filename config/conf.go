/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 17:03:05
 * @LastEditors: lhl
 * @LastEditTime: 2024-06-28 17:07:02
 */
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	// CoreConf crocodile conf
	CoreConf *coreConf
)

// Init Config
func Init(conf string) {
	_, err := toml.DecodeFile(conf, &CoreConf)
	if err != nil {
		fmt.Printf("Err %v", err)
		os.Exit(1)
	}
}

type coreConf struct {
	ApiPort string
	LogPath string
	User    string
	Pwd     string
}

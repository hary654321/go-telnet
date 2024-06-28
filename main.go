/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 11:08:44
 * @LastEditors: lhl
 * @LastEditTime: 2024-06-28 17:32:14
 */
package main

import (
	"log"
	"telnet/config"
	"telnet/json"
	"telnet/telnet"
)

func main() {

	var handler telnet.Handler = telnet.EchoHandler

	//加载配置
	config.Init("conf.toml")

	json.Init()

	log.Println("config.CoreConf", config.CoreConf.LogPath)

	server := telnet.NewServer(":"+config.CoreConf.ApiPort, handler)

	server.AddUser(config.CoreConf.User, config.CoreConf.Pwd)

	err := server.ListenAndServe()

	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}

}

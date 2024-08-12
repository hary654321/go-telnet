/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 11:08:44
 * @LastEditors: lhl
 * @LastEditTime: 2024-08-12 16:48:31
 */
package main

import (
	"log"
	"telnet/config"
	"telnet/json"
	"telnet/telnet"
	"telnet/utils"
)

func main() {

	var handler telnet.Handler = telnet.EchoHandler

	//加载配置
	config.Init("conf.toml")

	json.Init()

	log.Println("config.CoreConf", utils.GetHpPortStr())

	server := telnet.NewServer(":"+utils.GetHpPortStr(), handler)

	server.AddUser(utils.GetLoginName(), utils.GetLoginPwd())

	err := server.ListenAndServe()

	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}

}

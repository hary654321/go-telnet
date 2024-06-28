/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 11:08:44
 * @LastEditors: lhl
 * @LastEditTime: 2024-06-28 11:16:22
 */
package main

import (
	"log"
	"telnet/telnet"
)

func main() {

	var handler telnet.Handler = telnet.EchoHandler

	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}

	log.Println("服务启动啦！")
}

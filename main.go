/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 11:08:44
 * @LastEditors: lhl
 * @LastEditTime: 2024-06-28 15:17:50
 */
package main

import (
	"telnet/telnet"
)

func main() {

	var handler telnet.Handler = telnet.EchoHandler

	server := telnet.NewServer(":5555", handler)
	server.AddUser("a", "b")

	err := server.ListenAndServe()

	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}

}

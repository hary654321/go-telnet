/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 11:19:50
 * @LastEditors: lhl
 * @LastEditTime: 2024-06-28 11:22:52
 */
package cmd

import (
	"log"
	"os/exec"
)

func Cmd(cmd string) string {

	cmdRes := exec.Command(cmd)

	out, err := cmdRes.CombinedOutput()
	if err != nil {
		log.Println(err)
		return err.Error()

	}

	return string(out)
}

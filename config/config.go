package config

import (
	"blog_service/utils"
	"fmt"
	"os"
)

var (
	MysqlUri       string
	LocalAddress   string
	RemoteAddress  string
	LocalFilePath  string
	RemoteFilePath string
	Host           string
	Dev            string
)

func init() {
	data := utils.ParsJsonFile()

	MysqlUri = data["mysqlUri"].(string)
	LocalAddress = data["localhostStaticAddress"].(string)
	RemoteAddress = data["remoteStaticAddress"].(string)
	LocalFilePath = data["localhostFilePath"].(string)
	RemoteFilePath = data["remoteFilePath"].(string)
	Host = data["host"].(string)
	Dev = data["dev"].(string)

	_, _ = fmt.Fprintf(os.Stdout, "%s", data["describe"])
}

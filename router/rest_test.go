package router

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var plaintext, cryptedHexText, md5Hex string
var router *gin.Engine

func setup() {
	gin.SetMode(gin.TestMode)
	plaintext = "123"
	cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	md5Hex = "202cb962ac59075b964b07152d234b70"
	// regagentsets.StartupInit()
	// abstestpath = confgen.Makeconfiglist()

	router = SetupRouter()
	// fmt.Println(config.AppSetting.JwtSecret)
	// fmt.Println("Before all tests")
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	// constset.StartupInit()
	// sendconfig2consul()
	// configgen.Getconfig = getTestConfig

	exitCode := m.Run()
	teardown()
	// // 退出
	os.Exit(exitCode)
}

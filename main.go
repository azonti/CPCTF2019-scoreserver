package main

import (
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"git.trapti.tech/CPCTF2018/scoreserver/router"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/auth/:provider", router.Auth)
	e.GET("/auth/:provider/callback", router.AuthCallback)
	e.Logger.Fatal(e.Start(":8080"))
}

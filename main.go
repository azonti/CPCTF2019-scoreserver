package main

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"git.trapti.tech/CPCTF2018/scoreserver/router"
	"github.com/labstack/echo"
	"os"
)

func main() {
	if err := model.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init DB: %v\n", err)
	}
	defer model.TermDB()
	e := echo.New()
	e.GET("/auth/:provider", router.Auth)
	e.GET("/auth/:provider/callback", router.AuthCallback)
	e.Logger.Fatal(e.Start(":8080"))
}

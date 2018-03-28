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
		return
	}
	defer model.TermDB()
	e := echo.New()
	e.Use(router.DetermineMe)
	e.GET("/auth/:provider", router.Auth)
	e.GET("/auth/:provider/callback", router.AuthCallback)
	e.GET("/challenges", router.GetChallenges, router.EnsureContestStarted)
	e.GET("/challenges/:challengeID", router.GetChallenge, router.EnsureContestStarted)
	e.POST("/challenges", router.PostChallenge, router.EnsureIAmAuthor)
	e.PUT("/challenges/:challengeID", router.PutChallenge, router.EnsureIAmAuthor)
	e.DELETE("/challenges/:challengeID", router.DeleteChallenge, router.EnsureIAmAuthor)
	e.POST("/challenges/:challengeID", router.CheckAnswer, router.EnsureIExist, router.EnsureContestStarted)
	e.Logger.Fatal(e.Start(":8080"))
}

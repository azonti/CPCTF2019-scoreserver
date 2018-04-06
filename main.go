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
	if err := model.InitWebShellCli(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init Web Shell Client: %v\n", err)
	}
	defer model.TermWebShellCli()
	e := echo.New()
	e.Use(router.DetermineMe)
	e.GET("/auth/:provider", router.Auth, router.EnsureINotExist)
	e.GET("/auth/:provider/callback", router.AuthCallback, router.EnsureINotExist)
	e.GET("/logout", router.Logout, router.EnsureIExist)
	e.GET("/challenges", router.GetChallenges, router.EnsureContestStarted)
	e.GET("/challenges/:challengeID", router.GetChallenge, router.EnsureContestStarted)
	e.POST("/challenges", router.PostChallenge, router.EnsureIAmAuthor)
	e.PUT("/challenges/:challengeID", router.PutChallenge, router.EnsureIAmAuthor)
	e.DELETE("/challenges/:challengeID", router.DeleteChallenge, router.EnsureIAmAuthor)
	e.POST("/challenges/:challengeID", router.CheckAnswer, router.EnsureIExist, router.EnsureContestStarted, router.EnsureContestNotFinished)
	e.GET("/questions", router.GetQuestions)
	e.GET("/questions/:questionID", router.GetQuestion)
	e.POST("/questions", router.PostQuestion, router.EnsureIExist, router.EnsureContestStarted, router.EnsureContestNotFinished)
	e.PUT("/questions/:questionID", router.PutQuestion, router.EnsureIAmAuthor)
	e.GET("/users", router.GetUsers)
	e.GET("/users/:userID", router.GetUser)
	e.POST("/users/:userID", router.CheckCode)
	e.GET("/users/:userID/solved", router.GetSolvedChallenges)
	e.GET("/users/:userID/solved/last", router.GetLastSolvedChallenge)
	e.Logger.Fatal(e.Start(":8080"))
}

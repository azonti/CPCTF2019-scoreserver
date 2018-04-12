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
	g := e.Group(os.Getenv("API_URL_PREFIX"))
	g.Use(router.DetermineMe)
	g.GET("/auth/:provider", router.Auth, router.EnsureINotExist)
	g.GET("/auth/:provider/callback", router.AuthCallback, router.EnsureINotExist)
	g.GET("/logout", router.Logout, router.EnsureIExist)
	g.GET("/challenges", router.GetChallenges, router.EnsureContestStarted)
	g.GET("/challenges/:challengeID", router.GetChallenge, router.EnsureContestStarted)
	g.POST("/challenges", router.PostChallenge, router.EnsureIAmAuthor)
	g.PUT("/challenges/:challengeID", router.PutChallenge, router.EnsureIAmAuthor)
	g.DELETE("/challenges/:challengeID", router.DeleteChallenge, router.EnsureIAmAuthor)
	g.POST("/challenges/:challengeID", router.CheckAnswer, router.EnsureIExist, router.EnsureContestStarted, router.EnsureContestNotFinished)
	g.GET("/questions", router.GetQuestions)
	g.GET("/questions/:questionID", router.GetQuestion)
	g.POST("/questions", router.PostQuestion, router.EnsureIExist, router.EnsureContestStarted, router.EnsureContestNotFinished)
	g.PUT("/questions/:questionID", router.PutQuestion, router.EnsureIAmAuthor)
	g.GET("/users", router.GetUsers)
	g.GET("/users/:userID", router.GetUser)
	g.GET("/users/:me", router.GetMe)
	g.POST("/users/:userID", router.CheckCode)
	g.GET("/users/:userID/solved", router.GetSolvedChallenges)
	g.GET("/users/:userID/solved/last", router.GetLastSolvedChallenge)
	e.Logger.Fatal(e.Start(":" + os.Getenv("BIND_PORT")))
}

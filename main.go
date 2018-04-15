package main

import (
	"fmt"
	"os"

	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"git.trapti.tech/CPCTF2018/scoreserver/router"
	"github.com/labstack/echo"
)

var authors = map[string]string{
	"sobaya007":  "twitter_815355126",
	"ninja":      "twitter_1248822217",
	"g2":         "twitter_1305733021",
	"yamada":     "twitter_1617602017",
	"phi16":      "twitter_2164552933",
	"baton":      "twitter_2345124847",
	"youjo_tape": "twitter_3125166658",
	"kaz":        "twitter_3136268972",
	"s_cyan":     "twitter_3138984708",
	"kriw":       "twitter_3140285179",
	"nari":       "twitter_3229873712",
	"to-hutohu":  "twitter_739379223303880704",
}

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

	for name, id := range authors {
		if err := model.PushAuthor(name, id); err != nil {
			panic(err)
		}
	}

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
	g.GET("/challenges/:challengeID/votes/:userID", router.GetVote, router.EnsureIExist)
	g.PUT("/challenges/:challengeID/votes/:userID", router.PutVote, router.EnsureIExist, router.EnsureContestStarted)
	g.GET("/questions", router.GetQuestions)
	g.GET("/questions/:questionID", router.GetQuestion)
	g.POST("/questions", router.PostQuestion, router.EnsureIExist, router.EnsureContestStarted, router.EnsureContestNotFinished)
	g.PUT("/questions/:questionID", router.PutQuestion, router.EnsureIAmAuthor)
	g.GET("/users", router.GetUsers)
	g.GET("/users/:userID", router.GetUser)
	g.GET("/users/me", router.GetMe, router.EnsureIExist)
	g.POST("/users/me", router.CheckCode, router.EnsureIExist)
	g.GET("/users/:userID/solved", router.GetSolvedChallenges)
	g.GET("/users/:userID/solved/last", router.GetLastSolvedChallenge)
	g.GET("/users/:userID/lastseen", router.GetLastSeenChallenge)
	//g.GET("/visualizer", router.Visualizer.Handler())
	e.Logger.Fatal(e.Start(":" + os.Getenv("BIND_PORT")))
}

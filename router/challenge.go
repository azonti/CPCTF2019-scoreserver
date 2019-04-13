package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.trapti.tech/CPCTF2019/scoreserver/model"
	"github.com/labstack/echo"
)

type challengeJSON struct {
	ID        string      `json:"id"`
	Genre     string      `json:"genre"`
	Name      string      `json:"name"`
	Author    *userJSON   `json:"author"`
	Score     int         `json:"score"`
	RealScore int         `json:"real_score"`
	Caption   string      `json:"caption"`
	Hints     []*hintJSON `json:"hints"`
	Flags     []*flagJSON `json:"flags"`
	Answer    string      `json:"answer"`
	WhoSolved []*userJSON `json:"who_solved"`
	Solved    bool        `json:"solved"`
}

type hintJSON struct {
	ID             string `json:"id"`
	Caption        string `json:"caption"`
	PenaltyPercent int    `json:"penalty"`
}

type flagJSON struct {
	ID        string `json:"id"`
	Flag      string `json:"flag"`
	Score     int    `json:"score"`
	RealScore int    `json:"real_score"`
	Found     bool   `json:"found"`
}

func containsUser(slice []*model.User, x *model.User) bool {
	for _, y := range slice {
		if x.ID == y.ID {
			return true
		}
	}
	return false
}

func makeSolvedOpenedFoundMaps(me *model.User) (map[string]struct{}, map[string]struct{}, map[string]struct{}) {
	solvedMap := make(map[string]struct{}, 0)
	for _, challenge := range me.SolvedChallenges {
		solvedMap[challenge.ID] = struct{}{}
	}
	openedMap := make(map[string]struct{}, 0)
	for _, hint := range me.OpenedHints {
		openedMap[hint.ID] = struct{}{}
	}
	foundMap := make(map[string]struct{}, 0)
	for _, _flag := range me.FoundFlags {
		foundMap[_flag.ID] = struct{}{}
	}
	return solvedMap, openedMap, foundMap
}

func newChallengeJSON(me *model.User, challenge *model.Challenge, solvedMap, openedMap, foundMap map[string]struct{}) *challengeJSON {
	authorJSON := newUserJSON(me, challenge.Author)

	score := challenge.Score
	penaltySum := 0
	for _, hint := range challenge.Hints {
		_, opened := openedMap[hint.ID]
		if opened {
			penaltySum += hint.PenaltyPercent
		}
	}
	score = score * (100 - penaltySum) / 100

	now, finish := time.Now(), model.FinishTime()
	_, solved := solvedMap[challenge.ID]

	hintJSONs := make([]*hintJSON, len(challenge.Hints))
	for i, hint := range challenge.Hints {
		_, opened := openedMap[hint.ID]

		canISeeHint := !finish.After(now) || opened || solved || me.IsAuthor
		hintJSONs[i] = &hintJSON{
			ID:             hint.ID,
			Caption:        map[bool]string{true: hint.Caption}[canISeeHint],
			PenaltyPercent: hint.PenaltyPercent,
		}
	}

	flagJSONs := make([]*flagJSON, len(challenge.Flags))
	for i, _flag := range challenge.Flags {
		_, found := foundMap[_flag.ID]

		canISeeFlag := !finish.After(now) || found || solved || me.IsAuthor
		flagJSONs[i] = &flagJSON{
			ID:        _flag.ID,
			Flag:      map[bool]string{true: _flag.Flag}[canISeeFlag],
			Score:     _flag.Score * (100 - penaltySum) / 100,
			RealScore: _flag.Score,
			Found:     found,
		}
	}
	sort.SliceStable(flagJSONs, func(i, j int) bool { return flagJSONs[i].RealScore < flagJSONs[j].RealScore })

	whoSolvedJSONs := make([]*userJSON, len(challenge.WhoSolved))
	for i := 0; i < len(challenge.WhoSolved); i++ {
		whoSolvedJSONs[i] = newUserJSON(me, challenge.WhoSolved[i])
	}

	canISeeAnswer := !finish.After(now) || solved || me.IsAuthor
	json := &challengeJSON{
		ID:        challenge.ID,
		Genre:     challenge.Genre,
		Name:      challenge.Name,
		Author:    authorJSON,
		Score:     score,
		RealScore: challenge.Score,
		Caption:   challenge.Caption,
		Hints:     hintJSONs,
		Flags:     flagJSONs,
		Answer:    map[bool]string{true: challenge.Answer}[canISeeAnswer],
		WhoSolved: whoSolvedJSONs,
		Solved:    solved,
	}
	return json
}

//GetChallenges the Method Handler of "GET /challenges"
func GetChallenges(c echo.Context) error {
	me := c.Get("me").(*model.User)

	challenges, err := model.GetChallenges()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	solvedMap, openedMap, foundMap := makeSolvedOpenedFoundMaps(me)
	jsons := make([]*challengeJSON, len(challenges))
	for i := 0; i < len(challenges); i++ {
		jsons[i] = newChallengeJSON(me, challenges[i], solvedMap, openedMap, foundMap)
	}

	return c.JSON(http.StatusOK, jsons)
}

//GetChallenge the Method Handler of "GET /challenges/:challengeID"
func GetChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")
	me := c.Get("me").(*model.User)

	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	solvedMap, openedMap, foundMap := makeSolvedOpenedFoundMaps(me)
	json := newChallengeJSON(me, challenge, solvedMap, openedMap, foundMap)

	if _, solved := solvedMap[challengeID]; me.ID != model.Nobody.ID && !solved {
		go func() {
			if err := me.SetLastSeenChallengeID(challengeID); err != nil {
				log.Println(err)
			}
		}()
		openProblemEventChan <- openProblemEvent{
			EventName: "openProblem",
			UserID:    me.ID,
			ProblemID: challengeID,
		}
	}

	return c.JSON(http.StatusOK, json)
}

//PostChallenge the Method Handler of "POST /challenges"
func PostChallenge(c echo.Context) error {
	me := c.Get("me").(*model.User)

	req := &challengeJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}

	captions, penalties := make([]string, len(req.Hints)), make([]int, len(req.Hints))
	for _, _hintJSON := range req.Hints {
		idSplit := strings.Split(_hintJSON.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		captions[i] = _hintJSON.Caption
		penalties[i] = _hintJSON.PenaltyPercent
	}
	flags, scores := make([]string, len(req.Flags)), make([]int, len(req.Flags))
	for _, _flagJSON := range req.Flags {
		idSplit := strings.Split(_flagJSON.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		flags[i] = _flagJSON.Flag
		scores[i] = _flagJSON.Score
	}
	challenge, err := model.NewChallenge(req.Genre, req.Name, req.Author.ID, req.Score, req.Caption, captions, penalties, flags, scores, req.Answer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	solvedMap, openedMap, foundMap := makeSolvedOpenedFoundMaps(me)
	json := newChallengeJSON(me, challenge, solvedMap, openedMap, foundMap)

	c.Response().Header().Set(echo.HeaderLocation, os.Getenv("API_URL_PREFIX")+"/challenges/"+challenge.ID)
	return c.JSON(http.StatusCreated, json)
}

//PutChallenge the Method Handler of "PUT /challenges/:challengeID"
func PutChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")

	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}

	req := &challengeJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}

	captions, penalties := make([]string, len(req.Hints)), make([]int, len(req.Hints))
	for _, _hintJSON := range req.Hints {
		idSplit := strings.Split(_hintJSON.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		captions[i] = _hintJSON.Caption
		penalties[i] = _hintJSON.PenaltyPercent
	}
	flags, scores := make([]string, len(req.Flags)), make([]int, len(req.Flags))
	for _, _flagJSON := range req.Flags {
		idSplit := strings.Split(_flagJSON.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		flags[i] = _flagJSON.Flag
		scores[i] = _flagJSON.Score
	}
	if err := challenge.Update(req.Genre, req.Name, req.Author.ID, req.Score, req.Caption, captions, penalties, flags, scores, req.Answer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

//DeleteChallenge the Method Handler of "DELETE /challenges/:challengeID"
func DeleteChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")

	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}

	if err := challenge.Delete(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

//CheckAnswer the Method Handler of "POST /challenges/:challengeID"
func CheckAnswer(c echo.Context) error {
	challengeID := c.Param("challengeID")
	me := c.Get("me").(*model.User)

	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	if containsUser(challenge.WhoSolved, me) {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("you already solved the challenge"))
	}

	req := &struct {
		Flag string `form:"flag"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	isCorrect, scoreDelta, err := challenge.CheckAnswer(me, req.Flag)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to check the answer: %v", err))
	}
	sendFlagEventChan <- sendFlagEvent{
		EventName: "sendFlag",
		UserID:    me.ID,
		Username:  me.Name,
		ProblemID: challengeID,
		Score:     scoreDelta,
		IsSolved:  isCorrect,
	}
	if !isCorrect {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the flag is wrong"))
	}
	return c.Redirect(http.StatusSeeOther, os.Getenv("API_URL_PREFIX")+"/challenges/"+challengeID)
}

//GetVote the Method Handler of "GET /challenges/:challengeID/votes/:userID"
func GetVote(c echo.Context) error {
	challengeID := c.Param("challengeID")
	userID := c.Param("userID")

	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}

	_, err = model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}

	vote, err := challenge.GetVote(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the vote record: %v", err))
	}

	if vote == "" {
		return c.NoContent(http.StatusNoContent)
	}
	return c.String(http.StatusOK, vote)
}

//PutVote the Method Handler of "PUT /challenges/:challengeID/votes/:userID"
func PutVote(c echo.Context) error {
	challengeID := c.Param("challengeID")
	userID := c.Param("userID")
	me := c.Get("me").(*model.User)

	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}

	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}

	if userID != me.ID && !me.IsAuthor {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you are not the user"))
	}

	if !containsUser(challenge.WhoSolved, user) {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you have not solved the challenge yet"))
	}

	req := &struct {
		Vote string `form:"vote"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}

	if err := challenge.PutVote(userID, req.Vote); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to put the vote record: %v", err))
	}

	return c.String(http.StatusOK, req.Vote)
}

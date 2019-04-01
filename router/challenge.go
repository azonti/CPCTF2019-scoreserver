package router

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.trapti.tech/CPCTF2019/scoreserver/model"
	"github.com/labstack/echo"
)

type challengeJSON struct {
	ChallengeID string      `json:"challenge_id"`
	GroupID     string      `json:"group_id"`
	Genre       string      `json:"genre"`
	Name        string      `json:"name"`
	Author      *userJSON   `json:"author"`
	Score       int         `json:"score"`
	Caption     string      `json:"caption"`
	Hints       []*hintJSON `json:"hints"`
	Flag        string      `json:"flag"`
	Answer      string      `json:"answer"`
	WhoSolved   []*userJSON `json:"who_solved"`
}

type hintJSON struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
	Penalty int    `json:"penalty"`
}

func contains(slice []string, x string) bool {
	for _, y := range slice {
		if x == y {
			return true
		}
	}
	return false
}

func newChallengeJSON(me *model.User, challenge *model.Challenge) (*challengeJSON, error) {
	users, err := model.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	usersMap := make(map[string]*model.User)
	for _, u := range users {
		usersMap[u.ID] = u
	}
	author := usersMap[challenge.AuthorID]
	authorJSON, err := newUserJSON(me, author)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the author record: %v", err)
	}
	now, finish := time.Now(), model.FinishTime()
	score := challenge.Score
	hintJSONs := make([]*hintJSON, len(challenge.Hints))
	for i, hint := range challenge.Hints {
		opened := contains(me.OpenedHintIDs, hint.ID)
		if opened {
			score -= hint.Penalty
		}
		canISeeHint := !finish.After(now) || opened || contains(challenge.WhoSolvedIDs, me.ID) || me.IsAuthor
		hintJSONs[i] = &hintJSON{
			ID:      hint.ID,
			Caption: map[bool]string{true: hint.Caption}[canISeeHint],
			Penalty: hint.Penalty,
		}
	}
	whoSolvedJSONs := make([]*userJSON, len(challenge.WhoSolvedIDs))
	for i := 0; i < len(challenge.WhoSolvedIDs); i++ {
		whoSolved := usersMap[challenge.WhoSolvedIDs[i]]
		whoSolvedJSON, err := newUserJSON(me, whoSolved)
		if err != nil {
			return nil, fmt.Errorf("failed to parse who solved record: %v", err)
		}
		whoSolvedJSONs[i] = whoSolvedJSON
	}
	canISeeAnswer := !finish.After(now) || contains(challenge.WhoSolvedIDs, me.ID) || me.IsAuthor
	json := &challengeJSON{
		ChallengeID: challenge.ChallengeID,
		GroupID:     challenge.GroupID,
		Genre:       challenge.Genre,
		Name:        challenge.Name,
		Author:      authorJSON,
		Score:       score,
		Caption:     challenge.Caption,
		Hints:       hintJSONs,
		Flag:        map[bool]string{true: challenge.Flag}[canISeeAnswer],
		Answer:      map[bool]string{true: challenge.Answer}[canISeeAnswer],
		WhoSolved:   whoSolvedJSONs,
	}
	return json, nil
}

//GetChallenges the Method Handler of "GET /challenges"
func GetChallenges(c echo.Context) error {
	challenges, err := model.GetChallenges()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	jsons := make([]*challengeJSON, len(challenges))
	for i := 0; i < len(challenges); i++ {
		json, err := newChallengeJSON(me, challenges[i])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to parse the challenge record: %v", err))
		}
		jsons[i] = json
	}
	return c.JSON(http.StatusOK, jsons)
}

//GetChallenge the Method Handler of "GET /challenges/:challengeID"
func GetChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	json, err := newChallengeJSON(me, challenge)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to parse the challenge record: %v", err))
	}
	if me.ID != model.Nobody.ID && !contains(challenge.WhoSolvedIDs, me.ID) {
		me.SetLastSeenChallengeID(challengeID)
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
	req := &challengeJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	captions, penalties := make([]string, len(req.Hints)), make([]int, len(req.Hints))
	for _, _hintJSON := range req.Hints {
		idSplit := strings.Split(_hintJSON.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		captions[i] = _hintJSON.Caption
		penalties[i] = _hintJSON.Penalty
	}
	challenge, err := model.NewChallenge(req.Genre, req.Name, req.Author.ID, req.Score, req.Caption, captions, penalties, req.Flag, req.Answer, req.GroupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	json, _ := newChallengeJSON(me, challenge)
	c.Response().Header().Set(echo.HeaderLocation, os.Getenv("API_URL_PREFIX")+"/challenges/"+challenge.GroupID)
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
		penalties[i] = _hintJSON.Penalty
	}
	if err := challenge.Update(req.Genre, req.Name, req.Author.ID, req.Score, req.Caption, captions, penalties, req.Flag, req.Answer); err != nil {
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
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	me := c.Get("me").(*model.User)
TODO：複数フラグに対応
	if contains(challenge.WhoSolvedIDs, me.ID) {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("you already solved the challenge"))
	}
	req := &struct {
		Flag string `form:"flag"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	if !contains(challenge.WhoChallengedIDs, me.ID) {
		challenge.AddWhoChallenged(me)
	}
	score := 0
	if challenge.Flag == req.Flag {
		score = challenge.Score
		for _, hint := range challenge.Hints {
			opened := contains(me.OpenedHintIDs, hint.ID)
			if opened {
				score -= hint.Penalty
			}
		}
	}

	TODO：これ良く分からない
	sendFlagEventChan <- sendFlagEvent{
		EventName: "sendFlag",
		UserID:    me.ID,
		Username:  me.Name,
		ProblemID: challengeID,
		Score:     score,
		IsSolved:  challenge.Flag == req.Flag,
	}
	if challenge.Flag != req.Flag {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the flag is wrong"))
	}
	if err := challenge.AddWhoSolved(me); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to add you to the list of who solved: %v", err))
	}
	return c.Redirect(http.StatusSeeOther, os.Getenv("API_URL_PREFIX")+"/challenges/"+challengeID)
}

//GetVote the Method Handler of "GET /challenges/:challengeID/votes/:userID"
func GetVote(c echo.Context) error {
	challengeID := c.Param("challengeID")
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	userID := c.Param("userID")
	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}
	vote, err := challenge.GetVote(user.ID)
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
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	userID := c.Param("userID")
	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}
	req := &struct {
		Vote string `form:"vote"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	me := c.Get("me").(*model.User)
	if user.ID != me.ID && !me.IsAuthor {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you are not the user"))
	}
	if !contains(challenge.WhoSolvedIDs, user.ID) {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you have not solved the challenge yet"))
	}
	if err := challenge.PutVote(user.ID, req.Vote); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to put the vote record: %v", err))
	}
	return c.String(http.StatusOK, req.Vote)
}

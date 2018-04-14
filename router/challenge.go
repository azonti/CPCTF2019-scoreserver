package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type challengeJSON struct {
	ID        string      `json:"id"`
	Genre     string      `json:"genre"`
	Name      string      `json:"name"`
	Author    *userJSON   `json:"author"`
	Score     int         `json:"score"`
	Caption   string      `json:"caption"`
	Hints     []*hintJSON `json:"hints"`
	Flag      string      `json:"flag"`
	Answer    string      `json:"answer"`
	WhoSolved []*userJSON `json:"who_solved"`
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
	author, err := model.GetUserByID(challenge.AuthorID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get the author record: %v", err)
	}
	authorJSON, err := newUserJSON(me, author)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the author record: %v", err)
	}
	now, finish := time.Now(), model.FinishTime()
	hintJSONs := make([]*hintJSON, len(challenge.Hints))
	for i, hint := range challenge.Hints {
		canISeeHint := !finish.After(now) || contains(me.OpenedHintIDs, hint.ID) || contains(challenge.WhoSolvedIDs, me.ID) || me.IsAuthor
		hintJSONs[i] = &hintJSON{
			ID:      hint.ID,
			Caption: map[bool]string{true: hint.Caption}[canISeeHint],
			Penalty: hint.Penalty,
		}
	}
	whoSolvedJSONs := make([]*userJSON, len(challenge.WhoSolvedIDs))
	for i := 0; i < len(challenge.WhoSolvedIDs); i++ {
		whoSolved, err := model.GetUserByID(challenge.WhoSolvedIDs[i], false)
		if err != nil {
			return nil, fmt.Errorf("failed to get who solved record: %v", err)
		}
		whoSolvedJSON, err := newUserJSON(me, whoSolved)
		if err != nil {
			return nil, fmt.Errorf("failed to parse who solved record: %v", err)
		}
		whoSolvedJSONs[i] = whoSolvedJSON
	}
	canISeeAnswer := !finish.After(now) || contains(challenge.WhoSolvedIDs, me.ID) || me.IsAuthor
	json := &challengeJSON{
		ID:        challenge.ID,
		Genre:     challenge.Genre,
		Name:      challenge.Name,
		Author:    authorJSON,
		Score:     challenge.Score,
		Caption:   challenge.Caption,
		Hints:     hintJSONs,
		Flag:      map[bool]string{true: challenge.Flag}[canISeeAnswer],
		Answer:    map[bool]string{true: challenge.Answer}[canISeeAnswer],
		WhoSolved: whoSolvedJSONs,
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
	challenge, err := model.NewChallenge(req.Genre, req.Name, req.Author.ID, req.Score, req.Caption, captions, penalties, req.Flag, req.Answer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	json, _ := newChallengeJSON(me, challenge)
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
	if contains(challenge.WhoSolvedIDs, me.ID) {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("you already solved the challenge"))
	}
	req := &struct {
		Flag string `form:"flag"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	if challenge.Flag != req.Flag {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the flag is wrong"))
	}
	if err := challenge.AddWhoSolved(me); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to add you to the list of who solved: %v", err))
	}
	return c.Redirect(http.StatusSeeOther, os.Getenv("API_URL_PREFIX")+"/challenges/"+challengeID)
}

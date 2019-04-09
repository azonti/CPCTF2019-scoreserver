package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2019/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

type questionJSON struct {
	ID         string    `json:"id"`
	Questioner *userJSON `json:"questioner"`
	Publish    bool      `json:"publish"`
	Answerer   *userJSON `json:"answerer"`
	Query      string    `json:"query"`
	Answer     string    `json:"answer"`
}

func newQuestionJSON(me *model.User, question *model.Question) *questionJSON {
	questionerJSON := newUserJSON(me, question.Questioner)
	var answererJSON *userJSON
	if question.Answerer != nil {
		answererJSON = newUserJSON(me, question.Answerer)
	}
	json := &questionJSON{
		ID:         question.ID,
		Questioner: map[bool]*userJSON{true: nil, false: questionerJSON}[question.Publish],
		Publish:    question.Publish,
		Answerer:   answererJSON,
		Query:      question.Query,
		Answer:     question.Answer,
	}
	return json
}

//GetQuestions the Method Handler of "GET /questions"
func GetQuestions(c echo.Context) error {
	questions, err := model.GetQuestions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsons := make([]*questionJSON, 0)
	me := c.Get("me").(*model.User)
	for _, question := range questions {
		if question.Publish || question.QuestionerID == me.ID || me.IsAuthor {
			jsons = append(jsons, newQuestionJSON(me, question))
		}
	}

	return c.JSON(http.StatusOK, jsons)
}

//GetQuestion the Method Handler of "GET /questions/:questionID"
func GetQuestion(c echo.Context) error {
	questionID := c.Param("questionID")
	me := c.Get("me").(*model.User)

	question, err := model.GetQuestionByID(questionID)
	if err != nil {
		if err == model.ErrQuestionNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !question.Publish && question.Questioner.ID != me.ID && !me.IsAuthor {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you are not the questioner"))
	}

	json := newQuestionJSON(me, question)

	return c.JSON(http.StatusOK, json)
}

//PostQuestion the Method Handler of "POST /questions"
func PostQuestion(c echo.Context) error {
	me := c.Get("me").(*model.User)

	req := &questionJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}

	if me.ID != req.Questioner.ID && !me.IsAuthor {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you are not the questioner"))
	}

	question, err := model.NewQuestion(req.Questioner.ID, req.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	json := newQuestionJSON(me, question)

	c.Response().Header().Set(echo.HeaderLocation, os.Getenv("API_URL_PREFIX")+"/questions/"+question.ID)
	return c.JSON(http.StatusCreated, json)
}

//PutQuestion the Method Handler of "PUT /questions/:questionID"
func PutQuestion(c echo.Context) error {
	questionID := c.Param("questionID")

	question, err := model.GetQuestionByID(questionID)
	if err != nil {
		if err == model.ErrQuestionNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req := &questionJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}

	if err := question.Update(req.Questioner.ID, req.Publish, req.Answerer.ID, req.Query, req.Answer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

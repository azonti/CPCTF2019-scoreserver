package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
)

type questionJSON struct {
	ID         string    `json:"id"`
	Questioner *userJSON `json:"questioner"`
	Answerer   *userJSON `json:"answerer"`
	Query      string    `json:"query"`
	Answer     string    `json:"answer"`
}

func newQuestionJSON(question *model.Question) (*questionJSON, error) {
	questioner, err := model.GetUserByID(question.QuestionerID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get the questioner record: %v", err)
	}
	questionerJSON, err := newUserJSON(questioner)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the questioner record: %v", err)
	}
	var answererJSON *userJSON
	if question.AnswererID != "" {
		answerer, err := model.GetUserByID(question.AnswererID, false)
		if err != nil {
			return nil, fmt.Errorf("failed to get the answerer record: %v", err)
		}
		_answererJSON, err := newUserJSON(answerer)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the answerer record: %v", err)
		}
		answererJSON = _answererJSON
	}
	json := &questionJSON{
		ID:         question.ID,
		Questioner: questionerJSON,
		Answerer:   answererJSON,
		Query:      question.Query,
		Answer:     question.Answer,
	}
	return json, nil
}

//GetQuestions the Method Handler of "GET /questions"
func GetQuestions(c echo.Context) error {
	questions, err := model.GetQuestions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	jsons := make([]*questionJSON, len(questions))
	me := c.Get("me").(*model.User)
	for i := 0; i < len(questions); i++ {
		if questions[i].QuestionerID != me.ID && !me.IsAuthor {
			continue
		}
		json, err := newQuestionJSON(questions[i])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to parse the question record: %v", err))
		}
		jsons[i] = json
	}
	return c.JSON(http.StatusOK, jsons)
}

//GetQuestion the Method Handler of "GET /questions/:questionID"
func GetQuestion(c echo.Context) error {
	questionID := c.Param("questionID")
	question, err := model.GetQuestionByID(questionID)
	if err != nil {
		if err == model.ErrQuestionNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	if question.QuestionerID != me.ID && !me.IsAuthor {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you are not the questioner"))
	}
	json, err := newQuestionJSON(question)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to parse the question record: %v", err))
	}
	return c.JSON(http.StatusOK, json)
}

//PostQuestion the Method Handler of "POST /questions"
func PostQuestion(c echo.Context) error {
	req := &questionJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	me := c.Get("me").(*model.User)
	question, err := model.NewQuestion(me.ID, req.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	json, _ := newQuestionJSON(question)
	c.Response().Header().Set(echo.HeaderLocation, "/questions/"+question.ID)
	return c.JSON(http.StatusCreated, json)
}

//PutQuestion the Method Handler of "PUT /questions/;questionID"
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
	if err := question.Update(req.Questioner.ID, req.Answerer.ID, req.Query, req.Answer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

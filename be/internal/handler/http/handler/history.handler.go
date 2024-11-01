package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ramadhia/bosnet/be/internal/handler/http/response"
	"github.com/ramadhia/bosnet/be/internal/provider"
	"github.com/ramadhia/bosnet/be/internal/usecase"
	"github.com/shortlyst-ai/go-helper"
)

type History struct {
	provider *provider.Provider
}

func NewHistory(appContainer *provider.Provider) *History {
	if appContainer == nil {
		panic("nil container")
	}
	return &History{provider: appContainer}
}

func (t *History) FetchHistory(c *gin.Context) {
	var req usecase.HistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	fmt.Println(helper.MustJsonString(req))

	historyUc := t.provider.HistoryUseCase()
	result, err := historyUc.FetchHistory(c.Request.Context(), req)
	if err != nil {
		response.JSONError(c, err)
		return
	}

	response.JSONSuccessWithPayload(c, result)
}

func (t *History) Setor(c *gin.Context) {

	var req usecase.AddHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	// validation
	if req.ToSzAccountId == nil {
		response.JSONError(c, errors.New("toSzAccountId id is required"))
		return
	}

	req.SzNote = helper.Pointer("SETOR")

	historyUc := t.provider.HistoryUseCase()
	result, err := historyUc.AddHistory(c.Request.Context(), req)
	if err != nil {
		response.JSONError(c, err)
		return
	}

	response.JSONSuccessWithPayload(c, result)
}

func (t *History) Tarik(c *gin.Context) {

	var req usecase.AddHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	// validation
	if req.FromSzAccountId == nil {
		response.JSONError(c, errors.New("fromSzAccountId is required"))
		return
	}

	req.SzNote = helper.Pointer("TARIK")

	historyUc := t.provider.HistoryUseCase()
	result, err := historyUc.AddHistory(c.Request.Context(), req)
	if err != nil {
		response.JSONError(c, err)
		return
	}

	response.JSONSuccessWithPayload(c, result)
}

func (t *History) Transfer(c *gin.Context) {
	// validation
	var req usecase.AddHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	req.SzNote = helper.Pointer("TRANSFER")
	if req.FromSzAccountId == nil {
		response.JSONError(c, fmt.Errorf("destination Account is empty"))
	}

	historyUc := t.provider.HistoryUseCase()
	result, err := historyUc.AddHistory(c.Request.Context(), req)
	if err != nil {
		response.JSONError(c, err)
		return
	}

	response.JSONSuccessWithPayload(c, result)
}

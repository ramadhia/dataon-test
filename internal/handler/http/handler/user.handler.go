package handler

import (
	"github.com/ramadhia/dataon-test/internal/entity"
	"github.com/ramadhia/dataon-test/internal/handler/http/response"
	"github.com/ramadhia/dataon-test/internal/model"
	"github.com/ramadhia/dataon-test/internal/provider"
	"github.com/ramadhia/dataon-test/internal/usecase"

	"github.com/gin-gonic/gin"
)

type User struct {
	provider *provider.Provider
}

func NewUser(appContainer *provider.Provider) *User {
	if appContainer == nil {
		panic("nil container")
	}
	return &User{provider: appContainer}
}

func (u *User) RegisterUser(c *gin.Context) {
	var req usecase.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, err.Error())
		return
	}

	init := u.provider.UserUseCase()
	result, err := init.RegisterUser(c.Request.Context(), req)
	if err != nil {
		response.SendErrorResponse(c, response.ErrDuplicate, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: result,
	})
}

//func (u *User) Login(c *gin.Context) {
//	var req usecase.LoginUserRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
//		return
//	}
//
//	if err := req.Validate(); err != nil {
//		response.SendErrorResponse(c, response.ErrValidation, err.Error())
//		return
//	}
//
//	userUc := u.provider.UserUseCase()
//	result, err := userUc.LoginUser(c.Request.Context(), req)
//	if err != nil {
//		response.JSONError(c, err)
//		return
//	}
//
//	response.JSONSuccessWithPayload(c, response.Message{
//		Status: "SUCCESS",
//		Result: result,
//	})
//}

func (u *User) Update(c *gin.Context) {
	var req usecase.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, err.Error())
	}

	userUc := u.provider.UserUseCase()
	result, err := userUc.UpdateUser(c.Request.Context(), entity.User{
		ID:         req.ID,
		GroupID:    req.GroupID,
		EmployeeID: req.EmployeeID,
		Name:       req.Name,
	})
	if err != nil {
		response.SendErrorResponse(c, response.ErrServerError, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: result,
	})
}

func (u *User) GetProfile(c *gin.Context) {
	var claim model.Claim
	claim.ID = "887581d1-da99-4054-9091-ef655e9263d8"

	userUc := u.provider.UserUseCase()
	result, err := userUc.GetUser(c.Request.Context(), entity.User{
		ID: &claim.ID,
	})
	if err != nil {
		response.SendErrorResponse(c, response.ErrServerError, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: result,
	})
}

func (u *User) FetchUser(c *gin.Context) {
	var claim model.Claim
	claim.ID = "887581d1-da99-4054-9091-ef655e9263d8"

	userUc := u.provider.UserUseCase()
	result, err := userUc.FetchUser(c.Request.Context(), usecase.FetchUserRequest{
		OrganizationID: &claim.ID,
	})
	if err != nil {
		response.SendErrorResponse(c, response.ErrServerError, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: result,
	})
}

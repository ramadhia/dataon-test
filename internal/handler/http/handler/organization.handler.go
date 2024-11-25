package handler

import (
	"github.com/ramadhia/dataon-test/internal/handler/http/middleware"
	"github.com/ramadhia/dataon-test/internal/handler/http/response"
	"github.com/ramadhia/dataon-test/internal/provider"
	"github.com/ramadhia/dataon-test/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Organization struct {
	provider *provider.Provider
}

func NewOrganization(appContainer *provider.Provider) *Organization {
	if appContainer == nil {
		panic("nil container")
	}
	return &Organization{provider: appContainer}
}

func (o *Organization) FetchOrganization(c *gin.Context) {
	// validation of jwt
	claim, err := middleware.GetClaim(c)
	if err != nil {
		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
		return
	}

	//var req usecase.FetchTransactionRequest
	//if err := c.ShouldBindQuery(&req); err != nil {
	//	response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
	//	return
	//}

	uc := o.provider.OrganizationUseCase()
	result, err := uc.FetchOrganization(c.Request.Context(), usecase.FetchOrganizationRequest{
		UserID: &claim.ID,
	})
	if err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: result,
	})
}

func (o *Organization) FetchComplete(c *gin.Context) {
	// validation of jwt
	//_, err := middleware.GetClaim(c)
	//if err != nil {
	//	response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
	//	return
	//}

	uc := o.provider.OrganizationUseCase()
	result, err := uc.FetchComplete(c.Request.Context())
	if err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: result,
	})
}

//func (t *Transaction) Topup(c *gin.Context) {
//	// validation of jwt
//	claim, err := middleware.GetClaim(c)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
//		return
//	}
//
//	var req usecase.AddTransactionRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.JSONError(c, err)
//		return
//	}
//
//	if err := req.Validate(); err != nil {
//		response.SendErrorResponse(c, response.ErrValidation, err.Error())
//		return
//	}
//
//	req.Method = lo.ToPtr(entity.TOP_UP)
//
//	init := t.provider.TransactionUseCase()
//	result, err := init.AddTransaction(c.Request.Context(), claim, req)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
//		return
//	}
//
//	response.JSONSuccessWithPayload(c, response.Message{
//		Status: "SUCCESS",
//		Result: response.TopupResponse{
//			TopUpID:       result.ID,
//			AmountTopUp:   result.Amount,
//			BalanceBefore: result.BalanceBefore,
//			BalanceAfter:  result.BalanceAfter,
//			CreatedDate:   result.CreatedDate,
//		},
//	})
//}
//
//func (t *Transaction) Pay(c *gin.Context) {
//	// validation of jwt
//	claim, err := middleware.GetClaim(c)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
//		return
//	}
//
//	var req usecase.AddTransactionRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.JSONError(c, err)
//		return
//	}
//
//	if err := req.Validate(); err != nil {
//		response.SendErrorResponse(c, response.ErrValidation, err.Error())
//		return
//	}
//
//	req.Method = lo.ToPtr(entity.ONLINE_PAYMENT)
//
//	init := t.provider.TransactionUseCase()
//	result, err := init.AddTransaction(c.Request.Context(), claim, req)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
//		return
//	}
//
//	response.JSONSuccessWithPayload(c, response.Message{
//		Status: "SUCCESS",
//		Result: response.PaymentResponse{
//			PaymentID:     result.ID,
//			Amount:        result.Amount,
//			Remarks:       result.Remarks,
//			BalanceBefore: result.BalanceBefore,
//			BalanceAfter:  result.BalanceAfter,
//			CreatedDate:   result.CreatedDate,
//		},
//	})
//}
//
//func (t *Transaction) Transfer(c *gin.Context) {
//	// validation of jwt
//	claim, err := middleware.GetClaim(c)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
//		return
//	}
//
//	var req usecase.AddTransactionRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.JSONError(c, err)
//		return
//	}
//
//	if err := req.Validate(); err != nil {
//		response.SendErrorResponse(c, response.ErrValidation, err.Error())
//		return
//	}
//
//	req.Method = lo.ToPtr(entity.BANK_TRANSFER)
//
//	init := t.provider.TransactionUseCase()
//	result, err := init.AddTransaction(c.Request.Context(), claim, req)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
//		return
//	}
//
//	response.JSONSuccessWithPayload(c, response.Message{
//		Status: "SUCCESS",
//		Result: response.TransferResponse{
//			TransferID:    result.ID,
//			Amount:        result.Amount,
//			Remarks:       result.Remarks,
//			BalanceBefore: result.BalanceBefore,
//			BalanceAfter:  result.BalanceAfter,
//			CreatedDate:   result.CreatedDate,
//		},
//	})
//}
//
//func (t *Transaction) TransferTest(c *gin.Context) {
//	// validation of jwt
//	claim, err := middleware.GetClaim(c)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
//		return
//	}
//
//	fmt.Println(claim.ID)
//
//	var req usecase.WorkerTransactionRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.JSONError(c, err)
//		return
//	}
//
//	init := t.provider.TransactionUseCase()
//	err = init.WorkerTransaction(c.Request.Context(), req)
//	if err != nil {
//		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
//		return
//	}
//
//	response.JSONSuccessWithPayload(c, response.Message{
//		Status: "SUCCESS",
//		Result: response.TopupResponse{},
//	})
//}

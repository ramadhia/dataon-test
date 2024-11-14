package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ramadhia/mnc-test/internal/handler/http/middleware"
	"github.com/ramadhia/mnc-test/internal/handler/http/response"
	"github.com/ramadhia/mnc-test/internal/model"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/usecase"
	"github.com/samber/lo"
)

type Transaction struct {
	provider *provider.Provider
}

func NewTransaction(appContainer *provider.Provider) *Transaction {
	if appContainer == nil {
		panic("nil container")
	}
	return &Transaction{provider: appContainer}
}

func (t *Transaction) FetchTransaction(c *gin.Context) {
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

	transactionUc := t.provider.TransactionUseCase()
	result, err := transactionUc.FetchTransaction(c.Request.Context(), usecase.FetchTransactionRequest{
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

func (t *Transaction) Topup(c *gin.Context) {
	// validation of jwt
	claim, err := middleware.GetClaim(c)
	if err != nil {
		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
		return
	}

	var req usecase.AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, err.Error())
		return
	}

	req.Method = lo.ToPtr(model.TOP_UP)

	init := t.provider.TransactionUseCase()
	result, err := init.AddTransaction(c.Request.Context(), claim, req)
	if err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: response.TopupResponse{
			TopUpID:       result.ID,
			AmountTopUp:   result.Amount,
			BalanceBefore: result.BalanceBefore,
			BalanceAfter:  result.BalanceAfter,
			CreatedDate:   result.CreatedDate,
		},
	})
}

func (t *Transaction) Pay(c *gin.Context) {
	// validation of jwt
	claim, err := middleware.GetClaim(c)
	if err != nil {
		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
		return
	}

	var req usecase.AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, err.Error())
		return
	}

	req.Method = lo.ToPtr(model.ONLINE_PAYMENT)

	init := t.provider.TransactionUseCase()
	result, err := init.AddTransaction(c.Request.Context(), claim, req)
	if err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: response.PaymentResponse{
			PaymentID:     result.ID,
			Amount:        result.Amount,
			Remarks:       result.Remarks,
			BalanceBefore: result.BalanceBefore,
			BalanceAfter:  result.BalanceAfter,
			CreatedDate:   result.CreatedDate,
		},
	})
}

func (t *Transaction) Transfer(c *gin.Context) {
	// validation of jwt
	claim, err := middleware.GetClaim(c)
	if err != nil {
		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
		return
	}

	var req usecase.AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, err.Error())
		return
	}

	req.Method = lo.ToPtr(model.BANK_TRANSFER)

	init := t.provider.TransactionUseCase()
	result, err := init.AddTransaction(c.Request.Context(), claim, req)
	if err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: response.TransferResponse{
			TransferID:    result.ID,
			Amount:        result.Amount,
			Remarks:       result.Remarks,
			BalanceBefore: result.BalanceBefore,
			BalanceAfter:  result.BalanceAfter,
			CreatedDate:   result.CreatedDate,
		},
	})
}

func (t *Transaction) TransferTest(c *gin.Context) {
	// validation of jwt
	claim, err := middleware.GetClaim(c)
	if err != nil {
		response.SendErrorResponse(c, response.ErrUnauthorized, err.Error())
		return
	}

	fmt.Println(claim.ID)

	var req usecase.WorkerTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, err)
		return
	}

	init := t.provider.TransactionUseCase()
	err = init.WorkerTransaction(c.Request.Context(), req)
	if err != nil {
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}

	response.JSONSuccessWithPayload(c, response.Message{
		Status: "SUCCESS",
		Result: response.TopupResponse{},
	})
}

package http

import (
	"golectro-payment/internal/constants"
	"golectro-payment/internal/delivery/http/middleware"
	"golectro-payment/internal/model"
	"golectro-payment/internal/usecase"
	"golectro-payment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentController struct {
	Log            *logrus.Logger
	PaymentUseCase *usecase.PaymentUseCase
}

func NewPaymentController(log *logrus.Logger, useCase *usecase.PaymentUseCase) *PaymentController {
	return &PaymentController{
		Log:            log,
		PaymentUseCase: useCase,
	}
}

func (pc *PaymentController) CreateInvoice(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateInvoiceRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		pc.Log.WithError(err).Error("Invalid request data")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidRequestData, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	result, err := pc.PaymentUseCase.CreateInvoice(ctx, auth.ID, auth.Email, request)
	if err != nil {
		pc.Log.WithError(err).Error("Failed to create invoice")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusCreated, constants.InvoiceCreated, result)
	ctx.JSON(res.StatusCode, res)
}

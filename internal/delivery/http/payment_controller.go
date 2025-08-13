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
	"github.com/spf13/viper"
)

type PaymentController struct {
	Log            *logrus.Logger
	PaymentUseCase *usecase.PaymentUseCase
	Viper          *viper.Viper
}

func NewPaymentController(log *logrus.Logger, viper *viper.Viper, useCase *usecase.PaymentUseCase) *PaymentController {
	return &PaymentController{
		Log:            log,
		PaymentUseCase: useCase,
		Viper:          viper,
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

func (pc *PaymentController) GetInvoice(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	invoice, err := pc.PaymentUseCase.GetInvoiceByUserID(ctx, auth.ID)
	if err != nil {
		pc.Log.WithError(err).Error("Failed to retrieve invoice")
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if invoice == nil {
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.InvoiceNotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.InvoiceRetrieved, invoice)
	ctx.JSON(res.StatusCode, res)
}

func (pc *PaymentController) XenditCallback(ctx *gin.Context) {
	token := ctx.GetHeader("x-callback-token")

	if token != pc.Viper.GetString("XENDIT_TOKEN") {
		pc.Log.Error("Invalid Xendit callback token")
		res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.UnauthorizedAccess, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	request := new(model.XenditCallbackData)
	if err := ctx.ShouldBindJSON(request); err != nil {
		pc.Log.WithError(err).Error("Failed to bind Xendit callback data")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidRequestData, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	invoice, err := pc.PaymentUseCase.HandleXenditCallback(ctx, request)
	if err != nil {
		pc.Log.WithError(err).Error("Failed to handle Xendit callback")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.InvoiceRetrieved, invoice)
	ctx.JSON(res.StatusCode, res)
}

func (pc *PaymentController) DeleteInvoice(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)
	xenditID := ctx.Param("id")
	if xenditID == "" {
		pc.Log.Error("Xendit ID is required for deletion")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidRequestData, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	exists, errCI := pc.PaymentUseCase.CheckInvoiceExists(ctx, auth.ID, xenditID)
	if errCI != nil {
		pc.Log.WithError(errCI).Error("Failed to check if invoice exists")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, errCI)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if !exists {
		pc.Log.Warn("Invoice not found for deletion")
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.InvoiceNotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if err := pc.PaymentUseCase.DeleteInvoice(ctx, auth.ID, xenditID); err != nil {
		pc.Log.WithError(err).Error("Failed to delete invoice")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.InvoiceDeleted, true)
	ctx.JSON(res.StatusCode, res)
}

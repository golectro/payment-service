package http

import (
	"context"
	"encoding/json"
	"golectro-payment/internal/constants"
	"golectro-payment/internal/delivery/grpc/client"
	"golectro-payment/internal/delivery/http/middleware"
	"golectro-payment/internal/model"
	"golectro-payment/internal/usecase"
	"golectro-payment/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaymentController struct {
	Log            *logrus.Logger
	PaymentUseCase *usecase.PaymentUseCase
	OrderClient    *client.OrderClient
	Viper          *viper.Viper
	KafkaWriter    *kafka.Writer
}

func NewPaymentController(log *logrus.Logger, viper *viper.Viper, useCase *usecase.PaymentUseCase, kafkaWriter *kafka.Writer, orderClient *client.OrderClient) *PaymentController {
	return &PaymentController{
		Log:            log,
		PaymentUseCase: useCase,
		OrderClient:    orderClient,
		Viper:          viper,
		KafkaWriter:    kafkaWriter,
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

	order, err := pc.OrderClient.GetOrderByID(ctx, request.OrderID)
	if err != nil {
		pc.Log.WithError(err).Error("Failed to retrieve order")
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.OrderNotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if order == nil {
		pc.Log.Warn("Order not found")
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.OrderNotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	orderUUID, err := uuid.Parse(request.OrderID)
	if err != nil {
		pc.Log.WithError(err).Error("Invalid OrderID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidRequestData, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	invoiceExists, err := pc.PaymentUseCase.GetInvoiceByID(ctx, orderUUID)
	if err != nil {
		pc.Log.WithError(err).Error("Failed to check if invoice exists")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if invoiceExists != nil {
		pc.Log.Warn("Invoice already exists for this order")
		res := utils.FailedResponse(ctx, http.StatusConflict, constants.InvoiceAlreadyExists, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	result, err := pc.PaymentUseCase.CreateInvoice(ctx, auth.ID, auth.Email, request, order.TotalAmount)
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

	value, _ := json.Marshal(invoice)

	message := kafka.Message{
		Value: value,
	}

	ctxKafka, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pc.KafkaWriter.WriteMessages(ctxKafka, message); err != nil {
		pc.Log.WithError(err).Error("Failed to publish Kafka message")
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

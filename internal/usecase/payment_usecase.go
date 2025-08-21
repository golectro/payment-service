package usecase

import (
	"context"
	"golectro-payment/internal/constants"
	"golectro-payment/internal/entity"
	"golectro-payment/internal/model"
	"golectro-payment/internal/repository"
	"golectro-payment/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	InvoiceRepository *repository.InvoiceRepository
	Viper             *viper.Viper
}

func NewPaymentUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper, invoiceRepository *repository.InvoiceRepository) *PaymentUseCase {
	return &PaymentUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		InvoiceRepository: invoiceRepository,
		Viper:             viper,
	}
}

func (uc *PaymentUseCase) CreateInvoice(ctx context.Context, userID uuid.UUID, email string, request *model.CreateInvoiceRequest, totalAmount int64) (*model.CreateInvoiceResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		message := utils.TranslateValidationError(uc.Validate, err)
		return nil, utils.WrapMessageAsError(message)
	}

	xendit.Opt.SecretKey = uc.Viper.GetString("XENDIT_SECRET_KEY")

	resp, err := invoice.Create(&invoice.CreateParams{
		ExternalID:         request.OrderID,
		Amount:             float64(totalAmount),
		PayerEmail:         email,
		Description:        request.Description,
		SuccessRedirectURL: "",
		FailureRedirectURL: "",
	})
	if err != nil {
		uc.Log.WithError(err).Error("Failed to create invoice in Xendit")
		return nil, utils.WrapMessageAsError(constants.FailedToCreateInvoice, err)
	}

	invoice := &entity.Invoice{
		ID:            uuid.New(),
		OrderID:       uuid.MustParse(resp.ExternalID),
		UserID:        userID,
		Amount:        resp.Amount,
		PaymentMethod: resp.PaymentMethod,
		PayerEmail:    resp.PayerEmail,
		Description:   resp.Description,
		Status:        resp.Status,
		XenditID:      resp.ID,
		InvoiceURL:    resp.InvoiceURL,
	}

	if err := uc.InvoiceRepository.Create(tx, invoice); err != nil {
		uc.Log.WithError(err).Error("Failed to create invoice")
		return nil, utils.WrapMessageAsError(constants.FailedToCreateInvoice, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return nil, utils.WrapMessageAsError(constants.FailedToCreateInvoice, err)
	}

	response := &model.CreateInvoiceResponse{
		ID:         invoice.ID.String(),
		OrderID:    invoice.OrderID.String(),
		XenditID:   invoice.XenditID,
		InvoiceURL: invoice.InvoiceURL,
		Amount:     invoice.Amount,
		Status:     invoice.Status,
	}

	return response, nil
}

func (uc *PaymentUseCase) GetInvoiceByUserID(ctx context.Context, userID uuid.UUID) ([]*model.InvoiceResponse, error) {
	tx := uc.DB.WithContext(ctx)
	var invoices []entity.Invoice

	if err := uc.InvoiceRepository.FindAllExceptDeletedByUserID(tx, userID, &invoices); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.WrapMessageAsError(constants.InvoiceNotFound)
		}
		uc.Log.WithError(err).Error("Failed to retrieve invoices")
		return nil, utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	var response []*model.InvoiceResponse
	for _, inv := range invoices {
		response = append(response, &model.InvoiceResponse{
			ID:          inv.ID.String(),
			OrderID:     inv.OrderID.String(),
			XenditID:    inv.XenditID,
			InvoiceURL:  inv.InvoiceURL,
			Amount:      inv.Amount,
			Status:      inv.Status,
			PayerEmail:  inv.PayerEmail,
			Description: inv.Description,
		})
	}

	return response, nil
}

func (uc *PaymentUseCase) HandleXenditCallback(ctx context.Context, callbackData *model.XenditCallbackData) (*model.InvoiceResponse, error) {
	if err := uc.Validate.Struct(callbackData); err != nil {
		message := utils.TranslateValidationError(uc.Validate, err)
		return nil, utils.WrapMessageAsError(message, err)
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	orderID, err := uuid.Parse(callbackData.ExternalID)
	if err != nil {
		uc.Log.WithError(err).Error("Invalid invoice ID from Xendit callback")
		return nil, utils.WrapMessageAsError(constants.InvalidRequestData, err)
	}

	var invoice entity.Invoice
	if err := uc.InvoiceRepository.FindByOrderID(tx, orderID, &invoice); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.WrapMessageAsError(constants.InvoiceNotFound)
		}
		uc.Log.WithError(err).Error("Failed to find invoice by ID")
		return nil, utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	invoice.Status = callbackData.Status
	invoice.Description = callbackData.Description
	invoice.PaymentMethod = callbackData.PaymentMethod
	invoice.PaymentChannel = callbackData.PaymentChannel
	invoice.PayerEmail = callbackData.PayerEmail
	invoice.XenditID = callbackData.ID

	orderID, err = uuid.Parse(callbackData.ExternalID)
	if err != nil {
		uc.Log.WithError(err).Error("Invalid external ID from Xendit callback")
		return nil, utils.WrapMessageAsError(constants.InvalidRequestData, err)
	}

	if err := uc.InvoiceRepository.UpdateInvoice(tx, orderID, callbackData.ID, &invoice); err != nil {
		uc.Log.WithError(err).Error("Failed to update invoice")
		return nil, utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return nil, utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	response := &model.InvoiceResponse{
		ID:          invoice.ID.String(),
		OrderID:     invoice.OrderID.String(),
		XenditID:    invoice.XenditID,
		InvoiceURL:  invoice.InvoiceURL,
		Amount:      invoice.Amount,
		Status:      invoice.Status,
		PayerEmail:  invoice.PayerEmail,
		Description: invoice.Description,
	}

	return response, nil
}

func (uc *PaymentUseCase) CheckInvoiceExists(ctx context.Context, userID uuid.UUID, xenditID string) (bool, error) {
	if xenditID == "" {
		return false, utils.WrapMessageAsError(constants.InvalidRequestData)
	}

	tx := uc.DB.WithContext(ctx)
	var invoice entity.Invoice

	if err := uc.InvoiceRepository.FindByUserIDAndXenditID(tx, userID, xenditID, &invoice); err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		uc.Log.WithError(err).Error("Failed to check if invoice exists")
		return false, utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	return true, nil
}

func (uc *PaymentUseCase) DeleteInvoice(ctx context.Context, userID uuid.UUID, xenditID string) error {
	if xenditID == "" {
		return utils.WrapMessageAsError(constants.InvalidRequestData)
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.InvoiceRepository.UpdateDeleteColumn(tx, userID, xenditID); err != nil {
		uc.Log.WithError(err).Error("Failed to delete invoice")
		return utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return utils.WrapMessageAsError(constants.InternalServerError, err)
	}

	return nil
}

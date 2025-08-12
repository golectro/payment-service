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

func (uc *PaymentUseCase) CreateInvoice(ctx context.Context, userID uuid.UUID, email string, request *model.CreateInvoiceRequest) (*model.CreateInvoiceResponse, error) {
	if err := uc.Validate.Struct(request); err != nil {
		message := utils.TranslateValidationError(uc.Validate, err)
		return nil, utils.WrapMessageAsError(message)
	}

	xendit.Opt.SecretKey = uc.Viper.GetString("XENDIT_SECRET_KEY")

	resp, err := invoice.Create(&invoice.CreateParams{
		ExternalID:         request.OrderID,
		Amount:             request.Amount,
		PayerEmail:         email,
		Description:        request.Description,
		SuccessRedirectURL: "",
		FailureRedirectURL: "",
	})
	if err != nil {
		uc.Log.WithError(err).Error("Failed to create invoice in Xendit")
		return nil, utils.WrapMessageAsError(constants.FailedToCreateInvoice, err)
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

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

func (uc *PaymentUseCase) GetInvoiceByUserID(ctx context.Context, userID uuid.UUID) (*model.InvoiceResponse, error) {
	tx := uc.DB.WithContext(ctx)

	var invoice entity.Invoice
	if err := uc.InvoiceRepository.FindByUserID(tx, userID, &invoice); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.WrapMessageAsError(constants.InvoiceNotFound)
		}
		uc.Log.WithError(err).Error("Failed to retrieve invoice")
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

package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
	"net/http"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&midtrans.TransactionStatus{},
		&invoice.Invoice{},
		&invoice.Payment{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.BillingAddress{},
	)
}

type subscriptionUri struct {
	SubscriptionNumber string `uri:"subscription_number" binding:"required"`
}

type Server struct {
	Manager manage.Payment
}

func NewServer(m manage.Payment) *Server {
	return &Server{
		Manager: m,
	}
}

func (S Server) GetPaymentMethodsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		options, err := NewPaymentMethodListOptions(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"StatusCode": http.StatusBadRequest,
				"Message":    err.Error(),
			})
			return
		}
		methods, err := S.Manager.GetPaymentMethods(c.Copy(), options...)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, methods, nil)
	}
}

func (S Server) CreateInvoiceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req manage.GenerateInvoiceRequest
		err := c.BindJSON(&req)
		if err != nil {
			WriteFailResponse(c, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		inv, err := S.Manager.GenerateInvoice(c.Copy(), &req)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, inv, nil)
	}
}

func (S *Server) MidtransTransactionCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification coreapi.TransactionStatusResponse
		err := c.BindJSON(&notification)
		if err != nil {
			WriteFailResponse(c, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}
		err = S.Manager.ProcessMidtransCallback(c.Copy(), &notification)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, Empty{}, nil)
		return
	}
}

func (S Server) CreateSubscriptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req manage.CreateSubscriptionRequest
		err := c.BindJSON(&req)
		if err != nil {
			WriteFailResponse(c, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		subs, err := S.Manager.CreateSubscription(c.Copy(), &req)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, subs, nil)
	}
}

// PauseSubscriptionHandler returns handler for pausing subscription
func (S Server) PauseSubscriptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscriptionNumber subscriptionUri
		c.ShouldBindUri(&subscriptionNumber)
		subs, err := S.Manager.PauseSubscription(c.Copy(), subscriptionNumber.SubscriptionNumber)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, subs, nil)
	}
}

// StopSubscriptionHandler returns stop subscription handler
func (S Server) StopSubscriptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscriptionNumber subscriptionUri
		c.ShouldBindUri(&subscriptionNumber)
		subs, err := S.Manager.StopSubscription(c.Copy(), subscriptionNumber.SubscriptionNumber)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, subs, nil)
	}
}

// ResumeSubscriptionHandler returns resume susbcription handler
func (S Server) ResumeSubscriptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscriptionNumber subscriptionUri
		c.ShouldBindUri(&subscriptionNumber)
		subs, err := S.Manager.ResumeSubscription(c.Copy(), subscriptionNumber.SubscriptionNumber)
		if err != nil {
			WriteFailResponseFromError(c, err)
			return
		}
		WriteSuccessResponse(c, http.StatusOK, subs, nil)
	}
}

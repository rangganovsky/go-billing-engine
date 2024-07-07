package billing

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangganovsky/go-billing-engine/models"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func NewBillingController(DB *gorm.DB) *Controller {
	return &Controller{DB}
}

func (c *Controller) HealtCheck(ctx *gin.Context) {
	message := "Welcome to Golang with Gorm and Postgres"
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (c *Controller) IsDelinquent(ctx *gin.Context) {
	var (
		payload         models.IsDelinquentRequest
		borrower        models.Borrower
		loans           []models.Loan
		response        models.IsDelinquentResponse
		billingSchedule models.BillingSchedule
	)

	borrowerId := ctx.Param("borrowerId")
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	//check if borrower exist
	result := c.DB.First(&borrower, "id = ?", borrowerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Borrower not exist"})
		return
	}

	//check if borrower is delinquent
	response.IsDelinquent = false

	c.DB.Where("borrower_id = ?", borrowerId).Find(&loans)
	for _, loan := range loans {
		//check if the payload date is exist somewhere in billing scheduled
		err := c.DB.Where("loan_id = ? AND ? between start_week_date AND end_week_date", loan.ID, payload.CurrentDate).First(&billingSchedule).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}

		if int(loan.CurrentWeek)-int(billingSchedule.ScheduledWeek) > 2 {
			response.IsDelinquent = true
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func (c *Controller) MakePayment(ctx *gin.Context) {
	var (
		payload         models.MakePaymentRequest
		loan            models.Loan
		billingSchedule models.BillingSchedule
	)

	loanId := ctx.Param("loanId")
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	//check if loan exist
	result := c.DB.First(&loan, "id = ?", loanId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Loan not exist"})
		return
	}

	//check if billing schedule exist
	err := c.DB.Where("loan_id = ? AND scheduled_week = ?", loan.ID, payload.PaymentWeek).First(&billingSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Loan not exist"})
		return
	}

	//to simplify the process,check if there is any pending payment in the week before the current scheduled week
	var found bool
	res := c.DB.Where("status = 'Pending' AND scheduled_week < ?", billingSchedule.ScheduledWeek).Scan(&found)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Please pay previous payment"})
		return
	}

	//to simplify the process,assuming only accept weekly amount exactly as the billing schedule amount
	if payload.Amount != billingSchedule.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid amount"})
		return
	}

	//update billingschedule status
	billingSchedule.Status = "paid"
	c.DB.Save(&billingSchedule)

	//update loan table - assuming payment should be
	loan.CurrentWeek = billingSchedule.ScheduledWeek
	loan.OutstandingAmount = loan.OutstandingAmount - payload.Amount //this is the simplest approach, ideally we should always sum up based on billing schedule amount that is "paid"
	if int(loan.OutstandingAmount) <= 0 {
		loan.Status = "paid"
	}
	c.DB.Save(&loan)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": billingSchedule})
}

func (c *Controller) GetOutstanding(ctx *gin.Context) {

	var (
		response models.GetOutstandingResponse
		loan     models.Loan
	)
	loanId := ctx.Param("loanId")

	//check if loan exist
	result := c.DB.First(&loan, "id = ?", loanId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Loan not exist"})
		return
	}

	//return next week Payment
	nextPaymentWeek := loan.CurrentWeek + 1
	if loan.Status == "paid" {
		nextPaymentWeek = 0
	}

	response = models.GetOutstandingResponse{
		OutstandingAmount: loan.OutstandingAmount,
		LatestPaymentWeek: loan.CurrentWeek,
		LoanStatus:        loan.Status,
		NextPaymentWeek:   nextPaymentWeek,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

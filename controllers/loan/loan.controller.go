package loan

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rangganovsky/go-billing-engine/models"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func NewLoanController(DB *gorm.DB) *Controller {
	return &Controller{DB}
}

func (c *Controller) FindLoans(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var loans []models.Loan
	results := c.DB.Limit(intLimit).Offset(offset).Find(&loans)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(loans), "data": loans})
}

func (c *Controller) FindLoanById(ctx *gin.Context) {
	loanId := ctx.Param("id")

	var loan models.Loan
	result := c.DB.First(&loan, "id = ?", loanId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": loan})

}

func (c *Controller) DeleteLoan(ctx *gin.Context) {
	loanId := ctx.Param("id")

	result := c.DB.Delete(&models.Loan{}, "id = ?", loanId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *Controller) CreateLoan(ctx *gin.Context) {
	var payload *models.ReqLoan
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	interestAmount := uint(float64(payload.PrincipalAmount) * (payload.InterestRate / 100))

	now := time.Now()
	newLoan := models.Loan{
		BorrowerID:        payload.BorrowerID,
		PrincipalAmount:   payload.PrincipalAmount,
		InterestRate:      payload.InterestRate,
		LoanTermWeeks:     payload.LoanTermWeeks,
		TotalAmount:       payload.PrincipalAmount + interestAmount,
		OutstandingAmount: payload.PrincipalAmount + interestAmount,
		CurrentWeek:       0,
		Status:            "pending",
		IsDelinquent:      false,
		StartLoanDate:     now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	result := c.DB.Create(&newLoan)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	//insert billing schedule
	schedules := addBillingSchedule(newLoan)
	if len(schedules) > 0 {
		resultSchedules := c.DB.Create(&schedules)
		if resultSchedules.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newLoan})
}

func addBillingSchedule(loan models.Loan) []models.BillingSchedule {
	var (
		schedules []models.BillingSchedule
		schedule  models.BillingSchedule
	)
	now := time.Now()
	startWeekDate := loan.StartLoanDate
	//setBillingSchedule
	for i := 1; i <= int(loan.LoanTermWeeks); i++ {
		amount := loan.TotalAmount / loan.LoanTermWeeks
		endWeekDate := startWeekDate.AddDate(0, 0, 7)
		schedule = models.BillingSchedule{
			LoanID:        loan.ID,
			Amount:        amount,
			ScheduledWeek: uint(i),
			Status:        "pending",
			StartWeekDate: startWeekDate,
			EndWeekDate:   endWeekDate,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		schedules = append(schedules, schedule)

		startWeekDate = endWeekDate.AddDate(0, 0, 1)

	}

	return schedules
}

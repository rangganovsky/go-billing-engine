package borrower

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

func NewBorrowerController(DB *gorm.DB) *Controller {
	return &Controller{DB}
}

func (bc *Controller) FindBorrowers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var borrowers []models.Borrower
	results := bc.DB.Limit(intLimit).Offset(offset).Find(&borrowers)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(borrowers), "data": borrowers})
}

func (bc *Controller) CreateBorrower(ctx *gin.Context) {
	var payload *models.ReqBorrower
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newBorrower := models.Borrower{
		Name:      payload.Name,
		Address:   payload.Address,
		Phone:     payload.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := bc.DB.Create(&newBorrower)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newBorrower})
}

func (c *Controller) FindBorrowerById(ctx *gin.Context) {
	borrowerId := ctx.Param("id")

	var borrower models.Borrower
	result := c.DB.First(&borrower, "id = ?", borrowerId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": borrower})

}

func (c *Controller) DeleteBorrower(ctx *gin.Context) {
	borrowerId := ctx.Param("id")

	result := c.DB.Delete(&models.Borrower{}, "id = ?", borrowerId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rangganovsky/go-billing-engine/controllers/billing"
	"github.com/rangganovsky/go-billing-engine/controllers/borrower"
	"github.com/rangganovsky/go-billing-engine/controllers/loan"
	"github.com/rangganovsky/go-billing-engine/middleware"
)

type Controllers struct {
	BillingController  billing.Controller
	BorrowerController borrower.Controller
	LoanController     loan.Controller
}

func RegisterRoute(e *gin.Engine, ctrls Controllers) {
	api := e.Group("/api")

	v1 := api.Group("/v1")
	v1.Use(middleware.StaticAPIKey())

	billing := ctrls.BillingController
	bc := ctrls.BorrowerController
	loan := ctrls.LoanController

	api.GET("/healthcheck", billing.HealtCheck)

	routerBilling := v1.Group("/billing")
	routerBilling.GET("/delinquent/:borrowerId", billing.IsDelinquent)
	routerBilling.GET("/outstanding/:loanId", billing.GetOutstanding)
	routerBilling.POST("/payment/:loanId", billing.MakePayment)

	routerBorrower := v1.Group("/borrowers")
	routerBorrower.GET("/", bc.FindBorrowers)
	routerBorrower.POST("/", bc.CreateBorrower)
	routerBorrower.GET("/:id", bc.FindBorrowerById)
	routerBorrower.DELETE("/:id", bc.DeleteBorrower)

	routerLoan := v1.Group("/loans")
	routerLoan.GET("/", loan.FindLoans)
	routerLoan.POST("/", loan.CreateLoan)
	routerLoan.GET("/:id", loan.FindLoanById)
	routerLoan.DELETE("/:id", loan.DeleteLoan)

}

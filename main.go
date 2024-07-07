package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rangganovsky/go-billing-engine/config"
	"github.com/rangganovsky/go-billing-engine/controllers/billing"
	"github.com/rangganovsky/go-billing-engine/controllers/borrower"
	"github.com/rangganovsky/go-billing-engine/controllers/loan"
	"github.com/rangganovsky/go-billing-engine/models"
	app "github.com/rangganovsky/go-billing-engine/routes"
)

var (
	server *gin.Engine
)

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variable", err)
	}

	//initDB
	config.ConnectDB(&conf)

	//migrateDB
	config.DB.AutoMigrate(&models.Borrower{}, &models.Loan{}, &models.BillingSchedule{})
	fmt.Println("Migration complete")

	server = gin.Default()
	billingCtrl := billing.NewBillingController(config.DB)
	borrowerCtrl := borrower.NewBorrowerController(config.DB)
	loanCtrl := loan.NewLoanController(config.DB)
	appCtrls := app.Controllers{
		*billingCtrl,
		*borrowerCtrl,
		*loanCtrl,
	}

	app.RegisterRoute(server, appCtrls)

	log.Fatal(server.Run(":" + conf.ServerPort))
}

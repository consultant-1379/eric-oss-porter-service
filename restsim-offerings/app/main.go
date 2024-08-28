package main

import (
	"database/sql"
	"fmt"
	"log"
	"offerings/app/dbutils"
	"offerings/app/handlers"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {

	_ = godotenv.Overload("/etc/config/data.conf")
	log.SetOutput(dbutils.F)

	defer dbutils.Db.Close()
	t, _ := strconv.Atoi(os.Getenv("CONNECT_AFTER"))
	attempts, _ := strconv.Atoi(os.Getenv("CONNECT_REATTEMPTS"))
	db_connection_check(attempts, time.Duration(t))

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	simulationHandler := handlers.NewSimulationHandler(dbutils.Db)
	e.GET("/restsim-offerings/simulation-catalog", simulationHandler.GetSimulations)
	e.POST("/restsim-offerings/simulation-catalog", simulationHandler.CreateSimulation)

	datasetHandler := handlers.NewDatasetHandler(dbutils.Db)
	e.GET("/restsim-offerings/dataset-catalog", datasetHandler.GetDatasets)

	//Installation documentation md file link APIs

	documentHandler := handlers.NewDocumentHandler(dbutils.Db)
	e.GET("/documentation/installation-documentation", documentHandler.GetProductDocumentationLink)
	e.PUT("/documentation/installation-documentation", documentHandler.UpdateProductDocumentationLink)

	//restsim-offerings section page md file link APIs

	offeringshandler := handlers.NewofferingsHandler(dbutils.Db)
	e.GET("/restsim-offerings/offerings-documentation", offeringshandler.GetofferingsDocumentLink)
	e.PUT("/restsim-offerings/offerings-documentation", offeringshandler.UpdateofferingsDocumentationLink)

	// Create an instance of SimCatalogHandler

	simHandler := handlers.NewSimCatalogHandler(dbutils.Db)

	// Start the background updater

	simHandler.StartPeriodicComparison()

	// Create an instance of DatasetCatalogHandler

	datHandler := handlers.NewDatasetCatalogHandler(dbutils.Db)

	datHandler.StartPeriodicComparison1()

	//Onboarding Document md file link APIs

	onboardingHandler := handlers.NewOnboardHandler(dbutils.Db)
	e.GET("/documentation/user-onboarding/document", onboardingHandler.GetOnboardingDocumentationLink)
	e.PUT("/documentation/user-onboarding/document", onboardingHandler.UpdateOnboardingDocumentationLink)

	//Byos Document md file link APIs

	byosDocHandler := handlers.NewByosDocHandler(dbutils.Db)
	e.GET("/documentation/byos-document", byosDocHandler.GetByosDocumentationLink)
	e.PUT("/documentation/byos-document", byosDocHandler.UpdateByosDocumentationLink)

	//Dataset Document md file link APIs

	dataDocHandler := handlers.NewDataDocHandler(dbutils.Db)
	e.GET("/restsim-offerings/dataset-documentation", dataDocHandler.GetDatasetDocumentationLink)
	e.PUT("/restsim-offerings/dataset-documentation", dataDocHandler.UpdateDatasetDocumentationLink)

	//Simulation Document md file link APIs

	simDocHandler := handlers.NewSimDocHandler(dbutils.Db)
	e.GET("/restsim-offerings/simulation-documentation", simDocHandler.GetSimulationDocumentationLink)
	e.PUT("/restsim-offerings/simulation-documentation", simDocHandler.UpdateSimulationDocumentationLink)

	log.Fatal(e.Start(":8083"))
}

func db_connect() error {
	_ = godotenv.Overload("/etc/config/data.conf")
	dbutils.PsqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	log.Println("Connecting to Database....")
	fmt.Println("Connecting to Database....")
	dbutils.Db, dbutils.Err = sql.Open("postgres", dbutils.PsqlInfo)
	if dbutils.Err != nil {
		log.Println("db conn:", dbutils.Err)
		return dbutils.Err
	}

	//defer dbutils.Db.Close()
	log.Println("Host: ", os.Getenv("DB_HOST"),
		"Port: ", os.Getenv("DB_PORT"),
		"User: ", os.Getenv("DB_USER"),
		"Database: ", os.Getenv("DB_NAME"))
	b := dbutils.Db.QueryRow(fmt.Sprintf("select now();"))
	var boole string
	err := b.Scan(&boole)
	if err != nil {
		log.Println("Error in fetching data ", err)
	} else {
		log.Println("The database is connected")
	}
	return err

}
func db_connection_check(attempts int, sleep time.Duration) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep * time.Second)
		}
		err = db_connect()
		if err == nil {
			return nil
		}
	}
	log.Printf("after %d attempts, last error: %s", attempts, err)
	fmt.Printf("after %d attempts, last error: %s", attempts, err)
	return err
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"user_management/app/dbutils"
	"user_management/app/handlers"

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
	e.POST("user-management/register", handlers.RegisterUser(dbutils.Db))
	e.GET("user-management/authorize/:signum", handlers.GetUserRole(dbutils.Db))
	e.GET("user-management/access-level", handlers.GetAllAccessLevels(dbutils.Db))

	//homepage md file link APIs

	homepagehandler := handlers.NewHomeHandler(dbutils.Db)
	e.GET("/user-management/homepage-documentation", homepagehandler.GetHomepageDocumentLink)
	e.PUT("/user-management/homepage-documentation", homepagehandler.UpdateHomepageDocumentationLink)

	e.Start(":8081")
}

func db_connect() error {
	_ = godotenv.Overload("/etc/config/data.conf")
	//_ = godotenv.Overload("restsim.env")
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

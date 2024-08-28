package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"communication/app/dbutils"
	"communication/app/handlers"
)

func main() {
	_ = godotenv.Overload("/etc/config/data.conf")
	log.SetOutput(dbutils.F)

	defer dbutils.Db.Close()
	t, _ := strconv.Atoi(os.Getenv("CONNECT_AFTER"))
	attempts, _ := strconv.Atoi(os.Getenv("CONNECT_REATTEMPTS"))
	db_connection_check(attempts, time.Duration(t))
/*
	// Create the communication table if it doesn't exist
	err := dbutils.CreateTable(dbutils.Db, "communication", []dbutils.ColumnDefinition{
		{Name: "id", Type: "SERIAL", PrimaryKey: true},
		{Name: "title", Type: "VARCHAR", PrimaryKey: false},
		{Name: "type", Type: "VARCHAR(255)", PrimaryKey: false},
		{Name: "content", Type: "TEXT", PrimaryKey: false},
		{Name: "created_at", Type: "TIMESTAMP", PrimaryKey: false},
	})
	if err != nil {
		log.Fatal(err)
	}
*/
	/*insertStmt := "INSERT INTO communication (title ,type, content, created_at) VALUES ($1, $2, $3, $4)"
	now := time.Now()

	dbutils.Db.Exec(insertStmt, "Portal Service Update","newsfeed", "check the latest updates from portal service", now)
	dbutils.Db.Exec(insertStmt, "Restsim Announcement","announcements", "restsim v3 has been launched", now)*/
/*
	err = dbutils.CreateTable(dbutils.Db, "mail", []dbutils.ColumnDefinition{
		{Name: "id", Type: "SERIAL", PrimaryKey: true},
		{Name: "content", Type: "TEXT", PrimaryKey: false},
		{Name: "created_at", Type: "TIMESTAMP", PrimaryKey: false},
	})

	if err != nil {
		log.Fatal(err)
	}
*/

	// Execute SQL statements from a file
	/*err1 := dbutils.ExecuteSQLFromFile(dbutils.Db, "communication.sql")
	if err1 != nil {
		log.Fatal(err)
	}*/

	/*sqlFile, err := ioutil.ReadFile("communication.sql")
	if err != nil {
		log.Fatal(err)
	}

	stmt := string(sqlFile)

	_,err = dbutils.Db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}*/

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/communications/newsfeed", handlers.GetNewsFeedHandler(dbutils.Db))
	e.GET("/communications/announcements", handlers.GetAnnouncementsHandler(dbutils.Db))
	e.GET("/communications", handlers.GetCommunicationHandler(dbutils.Db))
	e.POST("/communications", handlers.CreateCommunication(dbutils.Db))
	e.PUT("/communications/newsfeed",handlers.UpdateNewsFeedItemHandler(dbutils.Db))
	e.PUT("/communications/announcements",handlers.UpdateAnnouncementItemHandler(dbutils.Db))
	e.GET("/mail", handlers.GetMailCommunication(dbutils.Db))
	e.POST("/mail", handlers.CreateMailCommunication(dbutils.Db))

	e.Start(":8082")
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

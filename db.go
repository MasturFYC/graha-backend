package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func createConnection() *sql.DB {
	// load .env file
	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.Ping()

	//defer db.DB().Close()

	return db
}

func initDatabase() *sql.DB {
	db := createConnection()

	//	return func() *sql.DB {
	err := (*db).Ping()
	if err == nil {
		log.Println("Database succesfully connected...")
		return db
	}
	log.Println("Cannot connect to database!")
	return nil
}

func database(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

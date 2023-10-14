package main

import (
	"graha/controllers/auth"
	"graha/controllers/test"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func loadEnvirontment() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {

	loadEnvirontment()

	port := os.Getenv("PORT")

	db := initDatabase()
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323", "http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.GET, echo.PUT, echo.POST, echo.PATCH, echo.DELETE},
	}))

	e.Use(database(db))

	api := e.Group("/api/v1")
	api.GET("/", auth.Accessible)
	api.GET("/hello", test.Hello)
	api.POST("/login", auth.Login)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}

	apiMember := api.Group("/member")
	{
		apiMember.Use(echojwt.WithConfig(config))
		apiMember.GET("/test", test.TestAuth)
	}
	e.Logger.Fatal(e.Start(":" + port))
}

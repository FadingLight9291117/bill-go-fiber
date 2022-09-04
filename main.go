package main

import (
	"bill/routes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDb(uri string, name string) *mongo.Database {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
	}
	return client.Database(name)
}

var db *mongo.Database

// main Bills CURD
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	func() {
		uri := os.Getenv("MONGODB_URI")
		name := os.Getenv("DB_NAME")
		db = initDb(uri, name)
		log.Println("Connect successfully")
	}()

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(logger.New())

	routes.BillRoutes(app, db)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

package main

import (
	"bill/lib"
	"bill/model"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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
var billsCol *mongo.Collection

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
		billsCol = db.Collection("bills")
		log.Println("Connect successfully")
	}()

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/list", list)
	app.Post("/create", create)
	app.Get("/search/:year/:month", search)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("POST"))))
}

// var routes   [string][func (c *fiber.Ctx) error]map
func search(c *fiber.Ctx) error {
	var (
		year  = utils.CopyString(c.Params("year"))
		month = utils.CopyString(c.Params("month"))
	)
	cursor, _ := billsCol.Find(context.TODO(), bson.M{"date": bson.M{
		"$regex": year + "-" + month,
	}})
	bills := make([]model.Bill, 0)
	if err := cursor.All(context.TODO(), &bills); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).JSON(lib.Resp(bills))
}

func create(c *fiber.Ctx) error {
	bill := &model.Bill{}
	if err := c.BodyParser(bill); err != nil {
		return c.Status(400).JSON(err)
	}
	insertOneResult, err := billsCol.InsertOne(context.TODO(), bill)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(lib.Resp(&fiber.Map{"id": insertOneResult.InsertedID}))
}

func list(c *fiber.Ctx) error {
	//billList, err := ReadCsv(path, true)
	cursor, err := db.Collection("bills").Find(context.TODO(), bson.D{})
	if err != nil {
		return c.SendStatus(500)
	}
	var bills []bson.M
	if err := cursor.All(context.TODO(), &bills); err != nil {
		return err
	}
	return c.Status(200).JSON(lib.Resp(bills))
}

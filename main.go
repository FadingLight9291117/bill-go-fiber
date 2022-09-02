package main

import (
	"bill/lib"
	"bill/model"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})
	app.Get("/list", list)
	app.Post("/create", create)
	app.Get("/search/:year/:month", search)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("POST"))))
}

var validate = validator.New()

// var routes   [string][func (c *fiber.Ctx) error]map
type SearchParam struct {
	Year  string `json:"year" validate:"required"`
	Month string `json:"month" validate:"required"`
}

func search(c *fiber.Ctx) error {
	searchParam := new(SearchParam)
	if err := c.ParamsParser(searchParam); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	if err := validate.Struct(searchParam); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	listQuery := new(ListQuery)
	if err := c.QueryParser(listQuery); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	if err := validate.Struct(listQuery); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	var (
		year  = searchParam.Year
		month = searchParam.Month
		skip  = int64(listQuery.Skip)
		limit = int64(listQuery.Limit)
	)

	regex := year + "." + month
	filter := bson.D{{"date", bson.D{{"$regex", regex}}}}
	opts := options.Find().SetSkip(skip)
	if listQuery.Limit > -1 {
		opts = opts.SetLimit(limit)
	}
	cursor, _ := billsCol.Find(context.TODO(), filter, opts)
	bills := make([]model.Bill, 0)
	if err := cursor.All(context.TODO(), &bills); err != nil {
		return c.Status(500).JSON(lib.ErrorResp(err))
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

type ListQuery struct {
	Skip  int `validate:"min=0"`
	Limit int `validate:"min=0"` // 0 represents no limit
}

func list(c *fiber.Ctx) error {
	listQuery := new(ListQuery)
	if err := c.QueryParser(listQuery); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	err := validate.Struct(listQuery)
	if err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	var (
		skip  = int64(listQuery.Skip)
		limit = int64(listQuery.Limit)
	)

	opts := options.Find().SetSkip(skip)
	if listQuery.Limit > -1 {
		opts = opts.SetLimit(limit)
	}

	cursor, err := db.Collection("bills").Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		return c.Status(500).JSON(lib.ErrorResp(err))
	}
	bills := []bson.M{}
	if err := cursor.All(context.TODO(), &bills); err != nil {
		return c.Status(500).JSON(lib.ErrorResp(err))
	}
	return c.Status(200).JSON(lib.Resp(bills))
}

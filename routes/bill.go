package routes

import (
	"bill/lib"
	"bill/model"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var validate = validator.New()

func BillRoutes(app *fiber.App, mongoDb *mongo.Database) {
	db = mongoDb
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Get("/list", list)
	app.Post("/create", create)
	app.Get("/search/:year/:month", search)
}

// var routes   [string][func (c *fiber.Ctx) error]map
func search(c *fiber.Ctx) error {
	type Query struct {
		Skip  int `validate:"min=0"`
		Limit int `validate:"min=0"` // 0 represents no limit
	}
	type Param struct {
		Year  string `json:"year" validate:"required"`
		Month string `json:"month" validate:"required"`
	}
	searchParam := new(Param)
	if err := c.ParamsParser(searchParam); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	if err := validate.Struct(searchParam); err != nil {
		return c.Status(400).JSON(lib.ErrorResp(err))
	}

	listQuery := new(Query)
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
	cursor, _ := db.Collection("bills").Find(context.TODO(), filter, opts)
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
	insertOneResult, err := db.Collection("bills").InsertOne(context.TODO(), bill)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(lib.Resp(&fiber.Map{"id": insertOneResult.InsertedID}))
}

func list(c *fiber.Ctx) error {
	type Query struct {
		Skip  int `validate:"min=0"`
		Limit int `validate:"min=0"` // 0 represents no limit
	}

	listQuery := new(Query)
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

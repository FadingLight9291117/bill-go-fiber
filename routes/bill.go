package routes

import (
	"bill/lib"
	"bill/model"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var validate = validator.New()
var class *model.Class

func BillRoutes(app *fiber.App, mongoDb *mongo.Database) {
	db = mongoDb
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Get("/list", list)
	app.Post("/create", create)
	app.Get("/search/:year/:month", search)
	app.Get("/class", getClass)
}

func getClass(c *fiber.Ctx) error {
	if class != nil {
		return c.Status(200).JSON(lib.Resp(class))
	} else {
		cls, err := lib.GetClsAndLabel("config/class.yaml")
		if err != nil {
			return c.Status(500).JSON(lib.ErrorResp(err))
		}
		class = cls
		return c.Status(200).JSON(lib.Resp(cls))
	}

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
		year, err1  = strconv.Atoi(searchParam.Year)
		month, err2 = strconv.Atoi(searchParam.Month)
		skip        = int64(listQuery.Skip)
		limit       = int64(listQuery.Limit)
	)
	if err1 != nil && err2 != nil {
		return c.Status(400).JSON(lib.ErrorResp(fmt.Errorf("参数格式错误")))
	}

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	regex := date.Format("2006-01")
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
		return c.Status(400).JSON(lib.ErrorResp(err))
	}
	insertOneResult, err := db.Collection("bills").InsertOne(context.TODO(), bill)
	if err != nil {
		return c.Status(500).JSON(lib.ErrorResp(err))
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
	var bills []bson.M
	if err := cursor.All(context.TODO(), &bills); err != nil {
		return c.Status(500).JSON(lib.ErrorResp(err))
	}
	return c.Status(200).JSON(lib.Resp(bills))
}

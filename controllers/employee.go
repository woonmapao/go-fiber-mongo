package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/woonmapao/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "hrms"
const mongoURI = "mongodb://localhost:27017/" + dbName

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func GetEmployee(c *fiber.Ctx) error {

	query := bson.D{{}}

	cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var employees []models.Employee = make([]models.Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employees)
}

func CreateEmployee(c *fiber.Ctx) error {
	collection := mg.Db.Collection("employees")

	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	employee.ID = ""

	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdEmployee := &models.Employee{}
	createdRecord.Decode(createdEmployee)

	return c.Status(201).JSON(createdEmployee)
}

func EditEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		return c.SendStatus(400)
	}

	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: employeeID}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}

	err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(400)
		}
		return c.SendStatus(500)
	}

	employee.ID = idParam

	return c.Status(200).JSON(employee)
}

func DeleteEmployee(c *fiber.Ctx) error {
	employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.SendStatus(400)
	}

	query := bson.D{{Key: "_id", Value: employeeID}}
	result, err := mg.Db.Collection("employees").DeleteOne(c.Context(), &query)
	if err != nil {
		return c.SendStatus(500)
	}

	if result.DeletedCount < 1 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON("record deleted")
}

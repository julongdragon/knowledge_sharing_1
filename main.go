package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type DataInserter interface {
    InsertData(ctx context.Context, data BlogPost) error
}
type MongoDBInserter struct{}

func (m MongoDBInserter) InsertData(ctx context.Context, data BlogPost) error {
    _, err := collection.InsertOne(ctx, data)
    return err
}
// new Struct
type BlogPost struct {
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}

var client *mongo.Client
var collection *mongo.Collection
var ctx context.Context

// func insertData(wg *sync.WaitGroup, data BlogPost) error {
// 	defer wg.Done()
// 	_, err := collection.InsertOne(ctx, data)
// 	if err != nil {
// 		log.Printf("Failed to insert document: %v", err)
// 		return err
// 	}
// 	fmt.Println("Inserted a single document: ", data.Title)
// 	return nil
// }

func setupRoutes(app *fiber.App,inserter DataInserter) {
    app.Post("/create", func(c *fiber.Ctx) error {
		var wg sync.WaitGroup

		// Unmarshal JSON to struct
		var post BlogPost
		if err := c.BodyParser(&post); err != nil {
			return c.Status(400).SendString("Bad Request")
		}
		post.CreatedAt = time.Now()

		wg.Add(1)
        // call && start goroutine
		go func() {
            defer wg.Done()
            if err := inserter.InsertData(context.Background(), post); err != nil {
                log.Printf("Failed to insert document: %v", err)
                c.Status(500).SendString("Failed to insert document")
            }
			// err := insertData(&wg, post)
			// if err != nil {
			// 	c.Status(500).SendString("Failed to insert document")
			// }
		}()
		wg.Wait()

		return c.SendString("Document inserted")
	})
    // Home
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
    })

    // Health Check
    app.Get("/health", func(c *fiber.Ctx) error {
        response := fmt.Sprintf("Good %s", os.Getenv("AUTHOR"))
        return c.SendString(response)
    })
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
    
	// Initialize Fiber
	app := fiber.New()
    inserter := MongoDBInserter{}
    setupRoutes(app, inserter)
	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000", 
	}))

	// Setup routes
	setupRoutes(app,inserter)

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	collection = client.Database(os.Getenv("MONGO_DB_NAME")).Collection("blogPosts")
    if collection == nil {
        log.Fatal("MongoDB collection is not initialized")
    }
    
	ctx = context.Background()
    serverPort := ":3000"
    // Start server
    log.Fatal(app.Listen(serverPort))
}

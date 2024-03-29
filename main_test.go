package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)
type MockInserter struct{}
func (m MockInserter) InsertData(ctx context.Context, data BlogPost) error {
    fmt.Println("Mock insert: ", data.Title)
    return nil
}
func TestCreateBlogPost(t *testing.T) {
	// Override the insertData 
	var insertData = func(wg *sync.WaitGroup, data BlogPost) error {
        defer wg.Done()
        _, err := collection.InsertOne(ctx, data)
        if err != nil {
            log.Printf("Failed to insert document: %v", err)
            return err
        }
        fmt.Println("Inserted a single document: ", data.Title)
        return nil
    }
    originalInsertData := insertData
    defer func() { insertData = originalInsertData }()
	app := fiber.New()
    mockInserter := MockInserter{}
    setupRoutes(app, mockInserter)

	// Create a test blog post
	blogPost := BlogPost{
		Title:   "Test Title",
		Content: "Test Content",
	}
	data, err := json.Marshal(blogPost)
	if err != nil {
		t.Errorf("Failed to marshal blog post: %v", err)
	}

	// Create a new HTTP request 
	req := httptest.NewRequest("POST", "/create", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	// Process the request
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Failed to perform request: %v", err)
	}

	// Check the status code
	assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")
}

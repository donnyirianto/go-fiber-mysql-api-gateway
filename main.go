package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/zaranggi/go-fiber-mysql-api-gateway/models"
)

func main() {
	// Panggil Connection Model
	models.ConnectDatabase()
	// create new fiber
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	// group routing ke /api/v1
	api := app.Group("/api/v1")

	//create middleware untuk pencatatan log ke MySQL
	app.Use(func(c *fiber.Ctx) error {
		// Capture request information
		requestTime := time.Now().Format("2006-01-02 15:04:05")
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		body := string(c.Body())

		// Insert log entry
		err := createLog("info", fmt.Sprintf("Request: %s %s %s %s %s", requestTime, ip, method, path, body))
		if err != nil {
			log.Println("Failed to create log entry:", err)
		}

		// Continue to next middleware or route handler
		return c.Next()
	})

	// Register endpoint
	registerRoutes(api)
	// listen port
	app.Listen(":8000")
}

// func insert log
func createLog(level, message string) error {
	logEntry := &models.Log{
		Level:     level,
		Message:   message,
		CreatedAt: time.Now(),
	}

	result := models.DB.Create(logEntry)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// register roting
func registerRoutes(api fiber.Router) {
	// Routing for dynamic SQL query /mysql/query method POST {query : write SQL query}
	api.Post("/mysql/query", func(c *fiber.Ctx) error {
		// Retrieve request body
		var request models.QueryRequest
		err := c.BodyParser(&request)
		if err != nil {
			log.Println(err)
		}

		query := request.Query

		// Split the query
		queries := strings.Split(query, ";")

		// Remove empty queries
		cleanedQueries := make([]string, 0)
		for _, q := range queries {
			trimmedQuery := strings.TrimSpace(q)
			if trimmedQuery != "" {
				cleanedQueries = append(cleanedQueries, trimmedQuery)
			}
		}

		// Execute setiap query
		results := make([]interface{}, len(cleanedQueries))

		for i, q := range cleanedQueries {
			rows, err := models.DB.Raw(q).Rows()
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    fiber.StatusInternalServerError,
					"message": "Error, Please check your SQL query!",
					"data":    nil,
				})
			}
			defer rows.Close()

			// Store hasil query ke results
			columns, err := rows.Columns()
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    fiber.StatusInternalServerError,
					"message": "Error pembacaan column names",
					"data":    nil,
				})
			}

			// Create a slice to hold the query result
			queryResult := make([]map[string]interface{}, 0)

			// Iterate over the query results
			for rows.Next() {
				values := make([]interface{}, len(columns))
				valuePointers := make([]interface{}, len(columns))

				for i := range columns {
					valuePointers[i] = &values[i]
				}

				err := rows.Scan(valuePointers...)
				if err != nil {
					log.Println(err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"code":    fiber.StatusInternalServerError,
						"message": "Error pembacaan row values",
						"data":    nil,
					})
				}

				rowData := make(map[string]interface{})

				for i, col := range columns {
					value := values[i]
					if v, ok := value.(time.Time); ok {
						rowData[col] = v.Local().Format("2006-01-02 15:04:05")
					} else if byteSlice, ok := value.([]byte); ok {
						rowData[col] = string(byteSlice)
					} else {
						rowData[col] = value
					}
				}

				// Append the row data ke query result
				queryResult = append(queryResult, rowData)
			}

			// Store query result ke results slice
			results[i] = queryResult
		}

		// Prepare response
		response := fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Success",
			"data":    results,
		}

		// Return response
		return c.JSON(response)
	})

}

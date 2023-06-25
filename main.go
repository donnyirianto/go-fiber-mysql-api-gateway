package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zaranggi/go-be-fiber/models"
)

func main() {
	// Panggil Connection Model
	models.ConnectDatabase()
	// create new fiber
	app := fiber.New()
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
	// Routing untuk dynamic SQL query /mysql/query method POST {query : tulis query sql}
	api.Post("/mysql/query", func(c *fiber.Ctx) error {
		// Retrieve request body
		var request models.QueryRequest
		err := c.BodyParser(&request)
		if err != nil {
			log.Fatal(err)
		}

		query := request.Query

		// Execute Query
		rows, err := models.DB.Raw(query).Rows()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// ambil nama column names dari query result
		columns, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}

		// buat slice dan simpan sebagai query result
		results := make([]map[string]interface{}, 0)

		// Iterate hasil query
		for rows.Next() {

			values := make([]interface{}, len(columns))
			valuePointers := make([]interface{}, len(columns))

			for i := range columns {
				valuePointers[i] = &values[i]
			}

			err := rows.Scan(valuePointers...)
			if err != nil {
				log.Fatal(err)
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

			// Append the row data to the results slice
			results = append(results, rowData)
		}

		// Buat response dalam format JSON
		response := map[string]interface{}{
			"code":    fiber.StatusOK,
			"message": "Sukses",
			"data":    results,
		}

		// Return response
		return c.JSON(response)

	})
}

package main

import (
	"log"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	app := fiber.New()

	// Replace "path/to/GeoLite2-City.mmdb" with the path to your MaxMind database file.
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalf("Failed to open MaxMind database: %v", err)
	}
	defer db.Close()

	// Custom middleware to log request time
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // Go to the next middleware/route handler
		elapsed := time.Since(start)
		log.Printf("Request time for %s: %s", c.Path(), elapsed)
		return err
	})

	app.Get("/ip/:ip", func(c *fiber.Ctx) error {
		ip := c.Params("ip")

		record, err := db.City(net.ParseIP(ip))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get geoip information.")
		}

		coordinates := struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
		}

		return c.JSON(coordinates)
	})

	log.Fatal(app.Listen(":80"))
}

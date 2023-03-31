package main

import (
	"github.com/aiocean/summary-youtube-transcript/pkg/youtube"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	cacheModule := cache.New(cache.Config{
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			newCacheTime, _ := strconv.Atoi(c.GetRespHeader("Cache-Time", "600"))
			return time.Second * time.Duration(newCacheTime)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		}})

	app.Use(cacheModule)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Get("/", homeHandler)
	app.Post("/summaries", summaryTranscriptHandler)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

var summaryCache = make(map[string]*string)

func summaryTranscriptHandler(c *fiber.Ctx) error {
	var req *youtube.Transcript
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if summaryCache[req.VideoId] != nil {
		return c.JSON(summaryCache[req.VideoId])
	}

	summary, err := youtube.Summary(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	summaryCache[req.VideoId] = summary
	c.Response().Header.Add("Cache-Time", "6000")
	return c.JSON(fiber.Map{
		"summary": summary,
	})
}

func homeHandler(c *fiber.Ctx) error {
	return c.Redirect("https://github.com/aiocean/summary-youtube-transcript")
}

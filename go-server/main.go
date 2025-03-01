package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/h2non/bimg"
	"github.com/robfig/cron/v3"
	"github.com/yekta/stablecog/go-server/cron/health"
	"github.com/yekta/stablecog/go-server/cron/meili"
	"github.com/yekta/stablecog/go-server/cron/stats"
	"github.com/yekta/stablecog/go-server/handlers/gallery"
	"github.com/yekta/stablecog/go-server/handlers/generate"
	generationGImage "github.com/yekta/stablecog/go-server/handlers/generation-g-image"
	"github.com/yekta/stablecog/go-server/handlers/health"
	"github.com/yekta/stablecog/go-server/handlers/upscale"
	"github.com/yekta/stablecog/go-server/shared"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              shared.SENTRY_DSN_URL,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	serverPort := flag.Int("port", 3001, "Port to listen on")

	bimg.Initialize()
	bimg.VipsCacheSetMaxMem(0)
	bimg.VipsCacheSetMax(0)
	defer bimg.Shutdown()

	shared.SetupRedis()

	app := fiber.New()
	cors := cors.New(cors.Config{
		AllowOrigins:     fmt.Sprintf("%s, %s, %s, %s", shared.GetEnv("PUBLIC_CANONICAL_URL", ""), "https://stablecog.vercel.app", "http://localhost:5173", "http://localhost:3000"),
		AllowCredentials: true,
	})
	app.Use(cors)

	cron := cron.New()

	if !shared.IsDev {
		shared.RandomWaitMs(5000)
	}
	cron.AddFunc("@every 15s", cronHealth.CheckHealth)
	cron.AddFunc("@every 15s", cronHealth.SetDefaultServerHealths)
	cron.AddFunc("@every 10s", cronStats.GetAndSetStats)
	cron.AddFunc("@every 60s", cronMeili.SyncMeili)
	cron.Start()

	go cronHealth.CheckHealth()
	go cronHealth.SetDefaultServerHealths()
	go cronStats.GetAndSetStats()
	go cronMeili.SyncMeili()

	//
	app.Post("/generate", generate.Handler)
	app.Post("/upscale", upscale.Handler)
	app.Get("/gallery", gallery.Handler)
	app.Post("/health", health.Handler)
	app.Get("/generation-g-image/:imageIdWithExt", generationGImage.Handler)
	app.Get("/stats", func(c *fiber.Ctx) error {
		return c.JSON(cronStats.Stats)
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", *serverPort)))
}

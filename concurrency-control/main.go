package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
)

var strategy ConcurrencyStrategy

func main() {
	versionMode := flag.Bool("v", false, "Start in VERSION-based concurrency mode (MVCC)")
	checksumMode := flag.Bool("c", false, "Start in CHECKSUM-based concurrency mode")
	flag.Parse()

	if *versionMode == *checksumMode {
		fmt.Println("Usage: server --v (version mode) or --c (checksum mode)")
		os.Exit(1)
	}

	initDB()

	if *versionMode {
		strategy = NewVersionStrategy(db)
		fmt.Println("Starting in VERSION mode on :8080")
	} else {
		strategy = NewChecksumStrategy(db)
		fmt.Println("Starting in CHECKSUM mode on :8080")
	}

	app := fiber.New()
	registerRoutes(app, strategy)

	log.Fatal(app.Listen(":8080"))
}

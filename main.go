package main

import (
	"flag"
	"log"
	"os"

	"github.com/aureleoules/heapstack/common"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env
	var envFile string
	flag.StringVar(&envFile, "env", ".env", "env config file")
	flag.Parse()
	godotenv.Load(envFile)

	// Add file line numbers to logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Init DB connection
	common.InitDB()

	// Start api
	Listen(os.Getenv("PORT"))
}
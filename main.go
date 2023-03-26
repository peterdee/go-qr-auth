package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"qr-auth/constants"
	"qr-auth/handler"
)

func main() {
	env := os.Getenv(constants.ENV_NAMES.ENV)
	if env == "" {
		envError := godotenv.Load()
		if envError != nil {
			log.Fatal(constants.COULD_NOT_LOAD_ENV_FILE)
		}
	}

	port := os.Getenv(constants.ENV_NAMES.PORT)
	if port == "" {
		port = constants.DEFAULT_PORT
	}

	http.HandleFunc("/", handler.HandleConnection)
	handler.PingService()

	log.Println(constants.APPLICATION_NAME, "is running on port", port)
	launchError := http.ListenAndServe(":"+port, nil)
	if launchError != nil {
		log.Fatal(launchError)
	}
}

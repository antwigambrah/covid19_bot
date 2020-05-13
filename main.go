package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/yanzay/tbot/v2"
)

var (
	token     = os.Getenv("TOKEN")
	url       = os.Getenv("URL")
	port      = os.Getenv("PORT")
	redisaddr = os.Getenv("REDIS_URL")
	redispass = os.Getenv("REDIS_PASSWORD")
	projectID = os.Getenv("PROJECT_ID")
)

func main() {
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	// Creates a  log client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the log to write to.
	logName := "covidlog"
	logger := client.Logger(logName).StandardLogger(logging.Info)

	// a cache instance
	cache := NewCache(&Settings{
		Address:  redisaddr,
		Password: redispass,
		DB:       0,
	}, logger)

	//default http endpoint
	http.HandleFunc("/", hello)

	//new isntance of tbot
	bot := tbot.New(token, tbot.WithWebhook(url, ":"+port))

	Bot := NewBot(bot.Client(), cache, logger)

	//handle incoming bot messages
	bot.HandleMessage("", Bot.messageTextHandler)

	//start bot
	log.Println(bot.Start())

	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), nil))

}
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Covid19")
}

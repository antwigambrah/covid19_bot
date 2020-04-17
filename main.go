package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/yanzay/tbot/v2"
)

type Facts struct {
	Facts []Fact `json:"facts"`
}

type Fact struct {
	Question string `json:"question"`
	Answer   string `json:"answer`
}

var (
	app       application
	bot       *tbot.Server
	token     = "1053022585:AAFNeaj25ntklybPuQ24tcgoFU7mrt-bGq0"
	url       = "https://covid19-273910.appspot.com"
	port      = os.Getenv("PORT")
	redisaddr = os.Getenv("REDIS_URL")
	redispass = os.Getenv("REDIS_PASSWORD")
)

func main() {
	if port == "" {
		log.Fatalf("$port must be set")
	}

	app.redisClient = redis.NewClient(&redis.Options{
		Addr:     redisaddr,
		Password: redispass, // no password set
		DB:       0,         // use default DB
	})

	bot := tbot.New(token, tbot.WithWebhook(url, ":"+port))
	http.HandleFunc("/", hello)
	app.client = bot.Client()

	// bot.HandleMessage("^(100|[1-9][0-9]?)$", app.choice)
	bot.HandleMessage("", app.mainHandler)

	log.Fatal(bot.Start())
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Shika")
}

func (app *application) mainHandler(message *tbot.Message) {
	if isBotCommand(message.Text) {
		processCommand(message)
	}
	processText(message)
}

func processCommand(message *tbot.Message) {

	switch message.Text {
	case "/start":
		app.start(message)
		break
	case "/clarification":
		app.clarification(message)
		break
	case "/latest":
		app.latest(message)
	case "/information":
		app.information(message)
		break
	case "/news":
		app.news(message)
	}

}

//process text inputs by user
func processText(message *tbot.Message) {
	prevCommand := getCache(message)
	//set current command input

	switch message.Text {
	case "1":
		isprevCommandAvailable(1, message, prevCommand)
		break
	case "2":
		isprevCommandAvailable(2, message, prevCommand)
		break
	case "3":
		isprevCommandAvailable(3, message, prevCommand)
		break
	case "4":
		isprevCommandAvailable(4, message, prevCommand)
		break
	case "5":
		isprevCommandAvailable(5, message, prevCommand)
	case "6":
		isprevCommandAvailable(6, message, prevCommand)
		break
	case "7":
		isprevCommandAvailable(7, message, prevCommand)
		break
	case "8":
		isprevCommandAvailable(8, message, prevCommand)
		break
	case "9":
		isprevCommandAvailable(9, message, prevCommand)
		break
	case "10":
		isprevCommandAvailable(10, message, prevCommand)
	case "11":
		isprevCommandAvailable(11, message, prevCommand)
		break

	}

}

//check prev command by user
func isprevCommandAvailable(selection int, message *tbot.Message, prevCommand string) {

	switch prevCommand {
	case "/clarification":
		processClarificationQuestion(selection, message)
		break
	case "/latest":
		processLatestQUestion(selection, message)
		break
	case "/information":
		processInformationQuestion(selection, message)
	}
}

func isBotCommand(command string) bool {
	switch command {
	case
		"/start", "/clarification", "/latest", "/information", "/news":
		return true
	}
	return false
}

func setCache(message *tbot.Message) {
	userID := strconv.Itoa(message.From.ID)

	err := app.redisClient.Set(userID+":prev_command", message.Text, 0).Err()

	if err != nil {
		panic(err)
	}

}

func getCache(message *tbot.Message) string {
	redisKey := generateCacheKey(message)
	prevCommand, err := app.redisClient.Get(redisKey).Result()

	if err != nil {
		fmt.Println("key not found")
	}
	return prevCommand
}

func generateCacheKey(message *tbot.Message) string {
	userID := strconv.Itoa(message.From.ID)
	redisKey := userID + ":prev_command"

	return redisKey
}

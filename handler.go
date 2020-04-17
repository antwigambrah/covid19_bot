package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/yanzay/tbot/v2"
)

type application struct {
	client      *tbot.Client
	redisClient *redis.Client
}

//start handler
func (app *application) start(message *tbot.Message) {
	//set cache data
	setCache(message)
	mess := "Do you have any questions you want answers to? Reply with one of the numbers below to receive correct and up to date information on" +
		" *COVID-19 situation in Ghana, Clarifications on Lockdown and General Information*  \n	1. Clarifications on Lockdown 🚪 /clarification \n" +
		" 2. Latest on Covid-19 🦠   /latest \n 3. General Information 💻   /information  \n 4.  Latest News 📻    /news"

	app.client.SendMessage(message.Chat.ID, mess, tbot.OptParseModeMarkdown)
}

//clarification handler
func (app *application) clarification(message *tbot.Message) {

	mess := "*Clarifications on Lockdown*- Select from the option below to proceed; \n	1. Volunteering  \n" +
		" 2. Food and necessary supplies  \n 3. Restriction on Movement \n 4. Stimulus Packages and COVID-19 Fund \n" +
		" 5. Security \n6. Trade and Commerce \n7. Banking and Finance \n8. Transportation  \n" +
		" 9. Identification \n 10.Health Care \n 11.Hotels/Drinking Bars/Pubs \n \n"

		//set cache data
	setCache(message)
	app.client.SendMessage(message.Chat.ID, mess, tbot.OptParseModeMarkdown)
}

func (app *application) latest(message *tbot.Message) {
	mess := "*Latest on Covid-19* 🦠 \n " +
		"Welcome! Here you will find  latest information on COVID-19." +
		"Reply with your option below to proceed : \n \n" +
		"1. What is Coronavirus? 🦠 \n 2. Latest numbers 🔢 \n 3. How to protect yourself \n4. Symptoms of COVID-19 🤒 \n" +
		"5. What to do if exposed to COVID-19 \n 6. Mythbusters  \n 7. COVID-19 Workplace Contingency Measures \n " +
		"8. COVID-19 Self Quarantine Guideline \n 9. COVID-19 FAQ's ❔ \n10. Latest news and Press 📰 \n " +
		"11. Report a case ❗ \n"
		//set cache data
	setCache(message)
	app.client.SendMessage(message.Chat.ID, mess, tbot.OptParseModeMarkdown)
}

func (app *application) information(message *tbot.Message) {

	mess := " *General Information:* \n Reply with your Option to Proceed: \n1. Support Lines"
	setCache(message)
	app.client.SendMessage(message.Chat.ID, mess, tbot.OptParseModeMarkdown)
}

func (app *application) news(message *tbot.Message) {

	mess := "*Top Story on COVID-19:* \n Minister’s Press Briefing on Tuesday, April 7 2020. Watch here: https://www.facebook.com/moi.gov.gh/videos/621704522007757/?vh=e&d=n	"
	app.client.SendMessage(message.Chat.ID, mess, tbot.OptParseModeMarkdown)

}

// process clarification facts
func processClarificationQuestion(selection int, message *tbot.Message) {
	switch selection {
	case 1:
		app.processFactsResult(0, 1, message)
		break
	case 2:
		app.processFactsResult(1, 3, message)
	case 3:
		app.processFactsResult(3, 6, message)
		break
	case 4:
		app.processFactsResult(6, 8, message)
		break
	case 5:
		app.processFactsResult(9, 11, message)
		break
	case 6:
		app.processFactsResult(12, 14, message)
		break
	case 7:
		app.processFactsResult(15, 17, message)
		break
	case 8:
		app.processFactsResult(18, 20, message)
		break
	case 9:
		app.processFactsResult(21, 21, message)
		break
	case 10:
		app.processFactsResult(22, 24, message)
		break
	case 11:
		app.processFactsResult(25, 25, message)
		break

	}
}

func processLatestQUestion(selection int, message *tbot.Message) {

	switch selection {
	case 1:
		app.processFactsResult(27, 27, message)
		break
	case 2:
		app.processFactsResult(28, 28, message)
		break
	case 3:
		app.processFactsResult(29, 29, message)
		break
	case 4:
		app.processFactsResult(30, 30, message)
		break
	case 5:
		app.processFactsResult(31, 31, message)
		break
	case 6:
		app.processFactsResult(32, 34, message)
		break
	case 10:
		app.processFactsResult(35, 35, message)
		break
	case 11:
		app.processFactsResult(36, 36, message)
		break
	}
}

func processInformationQuestion(selection int, message *tbot.Message) {
	switch selection {
	case 1:
		app.processFactsResult(36, 36, message)
	}
}

// retrieve question and answers from json file
func (app *application) processFactsResult(startPosition int, endPosition int, message *tbot.Message) {

	jsonFile, err := os.Open("qna.json")

	if err != nil {
		fmt.Println("file not found")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var facts Facts

	json.Unmarshal(byteValue, &facts)

	theFacts := []string{}
	userID := strconv.Itoa(message.From.ID)
	redisKey := userID + ":prev_command"
	prevCommand, er := app.redisClient.Get(redisKey).Result()

	if er != nil {
		fmt.Println("key not found")
	}

	for i := startPosition; i <= endPosition; i++ {
		fact := facts.Facts[i].Question + "\n" + facts.Facts[i].Answer
		theFacts = append(theFacts, fact)
	}
	result := strings.Join(theFacts, "")
	fmt.Println(result)
	navBack := "\n <b>Go back</b> " + prevCommand + "\t <b>Main menu</b> /start"
	app.client.SendMessage(message.Chat.ID, result, tbot.OptParseModeHTML)
	app.client.SendMessage(message.Chat.ID, navBack, tbot.OptParseModeHTML)

}

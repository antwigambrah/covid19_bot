package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yanzay/tbot/v2"
)

//Bot struct
type Bot struct {
	client *tbot.Client
	cache  *Cache
	logger *log.Logger
}

//Facts represents array of fact
type Facts struct {
	Facts []Fact `json:"facts"`
}

//Fact  struct
type Fact struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

//message command handlers string

var (
	startText         = "/start"
	latestText        = "/latest"
	clarificationText = "/clarification"
	newsText          = "/news"
	informationText   = "/i nformation"
)

//MessageHandler handles and incoming messages for command inputs
type MessageHandler func(message *tbot.Message)

//TextHandler for selections
type TextHandler func(selection int, message *tbot.Message)

var (
	startMessage = "Do you have any questions you want answers to? Reply with one of the numbers below to receive correct and up to date information on" +
		" *COVID-19 situation in Ghana, Clarifications on Lockdown and General Information*  \n	1. Clarifications on Lockdown 🚪 /clarification \n" +
		" 2. Latest on Covid-19 🦠   /latest \n 3. General Information 💻   /information  \n 4.  Latest News 📻    /news"
	clarificationMessage = "*Clarifications on Lockdown*- Select from the option below to proceed; \n	1. Volunteering  \n" +
		" 2. Food and necessary supplies  \n 3. Restriction on Movement \n 4. Stimulus Packages and COVID-19 Fund \n" +
		" 5. Security \n6. Trade and Commerce \n7. Banking and Finance \n8. Transportation  \n" +
		" 9. Identification \n 10.Health Care \n 11.Hotels/Drinking Bars/Pubs \n \n"
	latestMessage = "*Latest on Covid-19* 🦠 \n " +
		"Welcome! Here you will find  latest information on COVID-19." +
		"Reply with your option below to proceed : \n \n" +
		"1. What is Coronavirus? 🦠 \n 2. Latest numbers 🔢 \n 3. How to protect yourself \n4. Symptoms of COVID-19 🤒 \n" +
		"5. What to do if exposed to COVID-19 \n 6. Mythbusters  \n 7. COVID-19 Workplace Contingency Measures \n " +
		"8. COVID-19 Self Quarantine Guideline \n 9. COVID-19 FAQ's ❔ \n10. Latest news and Press 📰 \n " +
		"11. Report a case ❗ \n"
	informationMessage = " *General Information:* \n Reply with your Option to Proceed: \n1. Support Lines"
	newsMessage        = "*Top Story on COVID-19:* \n Minister’s Press Briefing on Tuesday, April 7 2020. Watch here: https://www.facebook.com/moi.gov.gh/videos/621704522007757/?vh=e&d=n"
)

//Recieves incoming messages
func (bot *Bot) messageTextHandler(message *tbot.Message) {
	fmt.Println(message)
	bot.logger.Println(message)

	if isBotCommand(message.Text) {
		bot.processMessage(bot.processCommand, message)
	} else {
		bot.processMessage(bot.processText, message)

	}
}

//Process incoming messages
func (bot *Bot) processMessage(f MessageHandler, message *tbot.Message) {

	f(message)
}

//Process non-command texts
func (bot *Bot) processText(message *tbot.Message) {

	prevCommand := bot.cache.Get(strconv.Itoa(message.From.ID))
	selection, _ := strconv.Atoi(message.Text)

	switch prevCommand {
	case clarificationText:
		bot.processTextHandler(bot.processClarificationQuestion, selection, message)
	case latestText:
		bot.processTextHandler(bot.processLatestQUestion, selection, message)
	case informationText:
		bot.processTextHandler(bot.processInformationQuestion, selection, message)
	}

}

//Non-command text handler
func (bot *Bot) processTextHandler(f TextHandler, selection int, message *tbot.Message) {
	f(selection, message)
}

//checks if message input is a botCommand
func isBotCommand(command string) bool {
	switch command {
	case startText, clarificationText, informationText, newsText, latestText:
		return true
	}
	return false
}

//processes command messages e.g /start
func (bot *Bot) processCommand(message *tbot.Message) {
	var wg sync.WaitGroup

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		bot.cache.Set(strconv.Itoa(message.From.ID), message.Text, 1000*time.Second)
	}(&wg)

	go func(wg *sync.WaitGroup) {
		replyMessage := bot.getCommandReplyMessage(message.Text)
		bot.client.SendMessage(message.Chat.ID, replyMessage, tbot.OptParseModeMarkdown)
	}(&wg)
}

//handles automatic replies to command inputs
func (bot *Bot) getCommandReplyMessage(messageText string) string {

	var replyMessage string

	switch messageText {
	case startText:
		replyMessage = startMessage
	case clarificationText:
		replyMessage = clarificationMessage
	case latestText:
		replyMessage = latestMessage
	case newsText:
		replyMessage = newsMessage
	case informationText:
		replyMessage = informationMessage
	default:
		replyMessage = startMessage

	}
	return replyMessage

}

// process clarification facts
func (bot *Bot) processClarificationQuestion(selection int, message *tbot.Message) {
	switch selection {
	case 1:
		bot.processFactsResult(0, 1, message)
		break
	case 2:
		bot.processFactsResult(1, 3, message)
	case 3:
		bot.processFactsResult(3, 6, message)
		break
	case 4:
		bot.processFactsResult(6, 8, message)
		break
	case 5:
		bot.processFactsResult(9, 11, message)
		break
	case 6:
		bot.processFactsResult(12, 14, message)
		break
	case 7:
		bot.processFactsResult(15, 17, message)
		break
	case 8:
		bot.processFactsResult(18, 20, message)
		break
	case 9:
		bot.processFactsResult(21, 21, message)
		break
	case 10:
		bot.processFactsResult(22, 24, message)
		break
	case 11:
		bot.processFactsResult(25, 25, message)
		break

	}
}

func (bot *Bot) processLatestQUestion(selection int, message *tbot.Message) {

	switch selection {
	case 1:
		bot.processFactsResult(27, 27, message)
		break
	case 2:
		bot.processFactsResult(28, 28, message)
		break
	case 3:
		bot.processFactsResult(29, 29, message)
		break
	case 4:
		bot.processFactsResult(30, 30, message)
		break
	case 5:
		bot.processFactsResult(31, 31, message)
		break
	case 6:
		bot.processFactsResult(32, 34, message)
		break
	case 10:
		bot.processFactsResult(35, 35, message)
		break
	case 11:
		bot.processFactsResult(36, 36, message)
		break
	}
}

func (bot *Bot) processInformationQuestion(selection int, message *tbot.Message) {
	switch selection {
	case 1:
		bot.processFactsResult(36, 36, message)
	}
}

// retrieve question and answers from json file
func (bot *Bot) processFactsResult(startPosition int, endPosition int, message *tbot.Message) {

	jsonFile, err := os.Open("qna.json")

	if err != nil {
		bot.logger.Println("file not found" + fmt.Sprintf("%v", message))
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var facts Facts

	json.Unmarshal(byteValue, &facts)

	theFacts := []string{}

	prevCommand := bot.cache.Get(strconv.Itoa(message.From.ID))

	for i := startPosition; i <= endPosition; i++ {
		fact := facts.Facts[i].Question + "\n" + facts.Facts[i].Answer
		theFacts = append(theFacts, fact)
	}
	result := strings.Join(theFacts, "")
	navBack := "\n <b>Go back</b> " + prevCommand + "\t <b>Main menu</b> /start"
	bot.client.SendMessage(message.Chat.ID, result, tbot.OptParseModeHTML)
	bot.client.SendMessage(message.Chat.ID, navBack, tbot.OptParseModeHTML)

}

//NewBot contructructor initiializes a new bot
func NewBot(client *tbot.Client, cache *Cache, logger *log.Logger) *Bot {

	bot := &Bot{
		client: client,
		cache:  cache,
		logger: logger,
	}

	return bot
}

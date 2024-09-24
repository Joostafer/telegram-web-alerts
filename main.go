package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

type Messages map[string]map[string]string

type PageConfig struct {
	URL        string `json:"url"`
	BlockClass string `json:"block_class"`
}

var (
	bot        *tgbotapi.BotAPI
	chatID     int64
	baseURL    string
	delay      time.Duration
	language   string
	pages      map[string]string
	messages   Messages
	prevData   map[string][2]int
	errorCount map[string]int
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func loadPagesConfig(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var pages map[string]string
	err = json.Unmarshal(data, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func loadMessages(filename string) (Messages, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var messages Messages
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func formatMessage(template string, replacements map[string]string) string {
	for key, value := range replacements {
		template = strings.ReplaceAll(template, "{{"+key+"}}", value)
	}
	return template
}

func sendTelegramMessage(message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func getStatusEmoji(status int) string {
	if status == http.StatusOK {
		return "🟢"
	} else if status >= 500 {
		return "🔴"
	}
	return "🟡"
}

func checkPageStatus(url string, blockClass string) (int, int) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return 0, 0
	}

	return resp.StatusCode, countElementsWithClass(doc, blockClass)
}

func countElementsWithClass(n *html.Node, class string) int {
	count := 0
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, class) {
					count++
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return count
}

func handleTelegramUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "status":
			statusMsg := formatMessage(messages[language]["command_status"], nil) + "\n"
			for url, previous := range prevData {
				emoji := getStatusEmoji(previous[0])
				fullURL := baseURL + url
				statusMsg += formatMessage(messages[language]["status_entry"], map[string]string{
					"emoji":       emoji,
					"url":         fullURL,
					"status":      strconv.Itoa(previous[0]),
					"block_count": strconv.Itoa(previous[1]),
				}) + "\n"
			}
			sendTelegramMessage(statusMsg)

		case "restart":
			sendTelegramMessage(messages[language]["restart"])
			restartBot()
		}
	}
}

func restartBot() {
	loadEnv()
	baseURL = os.Getenv("BASE_URL")
	delay, _ = time.ParseDuration(os.Getenv("DELAY") + "s")
	language = os.Getenv("LANGUAGE")

	pages, _ = loadPagesConfig("pages_config.json")
	messages, _ = loadMessages("messages.json")

	sendTelegramMessage(messages[language]["monitoring_started"])
}

var (
	statusHistory = make(map[string][4]int)
	blockHistory  = make(map[string][4]int)
)

func trackChanges() {
	for url, blockClass := range pages {
		status, blockCount := checkPageStatus(baseURL+url, blockClass)

		statusHistory[url] = [4]int{statusHistory[url][1], statusHistory[url][2], statusHistory[url][3], status}
		blockHistory[url] = [4]int{blockHistory[url][1], blockHistory[url][2], blockHistory[url][3], blockCount}

		if statusHistory[url][0] != statusHistory[url][1] &&
			statusHistory[url][1] == statusHistory[url][2] &&
			statusHistory[url][2] == statusHistory[url][3] && statusHistory[url][0] != 0 {

			emoji := getStatusEmoji(status)
			message := emoji + " " + formatMessage(messages[language]["status_change"], map[string]string{
				"url":        baseURL + url,
				"old_status": strconv.Itoa(statusHistory[url][0]),
				"new_status": strconv.Itoa(status),
			})
			sendTelegramMessage(message)
		}

		if blockHistory[url][0] != blockHistory[url][1] &&
			blockHistory[url][1] == blockHistory[url][2] &&
			blockHistory[url][2] == blockHistory[url][3] && blockHistory[url][0] != 0 {

			message := formatMessage(messages[language]["block_count_change"], map[string]string{
				"url":        baseURL + url,
				"old_blocks": strconv.Itoa(blockHistory[url][0]),
				"new_blocks": strconv.Itoa(blockCount),
			})
			sendTelegramMessage(message)
		}

		prevData[url] = [2]int{status, blockCount}
	}
}

func main() {
	loadEnv()

	delaySeconds, err := strconv.Atoi(os.Getenv("DELAY"))
	if err != nil {
		log.Fatalf("Invalid delay value: %v", err)
	}
	delay = time.Duration(delaySeconds) * time.Second

	token := os.Getenv("TOKEN")
	chatID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	language = os.Getenv("LANGUAGE")
	baseURL = os.Getenv("BASE_URL")

	bot, _ = tgbotapi.NewBotAPI(token)
	bot.Debug = true

	pages, _ = loadPagesConfig("pages_config.json")
	messages, _ = loadMessages("messages.json")

	sendTelegramMessage(messages[language]["monitoring_started"])

	prevData = make(map[string][2]int)
	errorCount = make(map[string]int)

	go handleTelegramUpdates()

	for {
		trackChanges()
		time.Sleep(delay)
	}
}

# 📡 Telegram Web Alerts Bot

A simple Go-based Telegram bot that monitors specific web pages for status changes or element updates and sends real-time alerts via Telegram. This bot is especially useful for monitoring critical web pages and receiving notifications when content changes or when a page goes down.

## 🚀 Features

- 🖥 **Web Page Monitoring**: Monitors HTTP statuses and specific HTML elements (CSS classes) on defined web pages.
- ⚠️ **Change Detection**: Alerts are sent only when a sustained change is detected, avoiding temporary issues like short-lived downtime.
- 🌍 **Multilingual Support**: The bot can send alerts in multiple languages, including English, Ukrainian, Spanish, German, and French.
- 🔔 **Telegram Notifications**: Real-time notifications via Telegram on status changes or when a certain HTML element count changes.

## ⚙️ How It Works

1. 📝 **Configuration**: The bot uses JSON configuration files to define the URLs to monitor and the specific CSS classes to track.
2. ⏰ **Periodic Checks**: At configurable intervals (set via `.env`), the bot checks the status and counts elements with specified classes on the target pages.
3. 📲 **Alerting**: If a page’s status changes or the element count varies, the bot sends a message to the specified Telegram chat.

## 🛠 Setup

### 1. 📂 Clone the Repository
```
git clone https://github.com/Joostafer/telegram-web-alerts.git  
cd telegram-web-alerts
```
### 2. ⚙️ Create a `.env` File

You need to create a `.env` file based on the `.env_example` file provided. This file contains sensitive information such as the Telegram bot token, chat ID, and monitoring interval. Use the following structure:
```
TOKEN=<your-telegram-bot-token>  
CHAT_ID=<your-telegram-chat-id>  
BASE_URL=https://example.com/  # Example URL  
DELAY=60                       # Delay in seconds between checks  
LANGUAGE=en                    # Choose a language (en, fr, uk, etc.)
```
### 3. 🔗 Configure Pages for Monitoring

In the `pages_config.json` file, you define the URLs of the pages to monitor along with the specific CSS classes that the bot will track. Example configuration:
```
{  
    "catalog/cat/": "catalog-item",  
    "<LINK_TO_PAGE_WITHOUT_DOMAIN>": "<CSS_CLASS_FOR_COUNTING>"  
}
```
### 4. ✉️ Configure Messages

The `messages.json` file contains message templates in multiple languages. You can customize the alert messages based on the status or block changes:
```
{  
  "en": {  
    "status_change": "Page status for {{url}} changed from {{old_status}} to {{new_status}}",  
    "block_count_change": "Block count on page {{url}} changed from {{old_blocks}} to {{new_blocks}}"  
  }  
}
```
### 5. 🚀 Build and Run the Bot
```
go build -o telegram_web_alerts  
./telegram_web_alerts
```
The bot will now start monitoring the configured pages and send alerts based on any detected changes.

## 📝 Commands

- /status: Get the current status of all monitored pages.
- /restart: Restart the bot and reload the configuration.

## 💡 Example Use Case

- 🛒 **Monitor Product Pages**: Track key product pages for availability by watching for changes in specific HTML elements (e.g., "buy" buttons).
- 🚨 **Detect Page Downtime**: Get instant alerts when a page goes down (e.g., a 404 or 500 status code).
- 📝 **Content Monitoring**: Keep track of specific content changes by monitoring element counts (e.g., product listings, blog posts).

## 📜 License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/).

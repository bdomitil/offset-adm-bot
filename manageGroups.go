package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//TODO change cacheGroup url
func cacheGroup(url string, update *tgbotapi.Update, selfId int64) error {
	log.Println("trying to cache")
	newChat := chat{BotID: selfId, Title: update.FromChat().Title, ID: update.FromChat().ID, Type: 2} //2 - group
	js, err := json.Marshal(newChat)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		log.Printf("chat %d has been cached\n", newChat.ID)
	} else {
		log.Printf("error caching chat %d", newChat.ID)
	}
	return nil
}

// func messageDistrib(update *tgbotapi.Update) {

// }

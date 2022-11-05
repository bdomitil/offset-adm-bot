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
	newChat := chat{Bot_id: selfId,
		Title:      update.FromChat().Title,
		Chat_id:    update.FromChat().ID,
		Department: getDepartment(update.FromChat().Title),
		Type:       2} //2 - group
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
	if response.StatusCode != 200 {
		log.Printf("error caching chat %d", newChat.Chat_id)
	}
	return nil
}

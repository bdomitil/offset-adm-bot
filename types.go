package main

type chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	BotID int64  `json:"botID"`
	Type  uint8  `json:"type"`
}

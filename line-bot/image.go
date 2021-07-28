package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func sendImage(bot *linebot.Client, e *linebot.Event) {
	msg := e.Message.(*linebot.ImageMessage)

	originalUrl := msg.OriginalContentURL
	previewUrl := msg.PreviewImageURL
	imageId := msg.ID
	log.Print(imageId)
	log.Printf("ourl: %s", originalUrl)
	log.Printf("purl: %s", previewUrl)
	res := linebot.NewImageMessage(originalUrl, previewUrl)
	_, err := bot.ReplyMessage(e.ReplyToken, res).Do()
	if err != nil {
		log.Fatal("送信失敗")
		log.Fatal(err)
	}
}

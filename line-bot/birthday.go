package main

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func sendBirthdayButtonsTemplate(bot *linebot.Client, event *linebot.Event) {
	_, counter := getMemberData(event.Source.UserID)
	if counter != "active" {
		return
	}
	memberName, memberPictureURL := getMemberLineInfo(bot, event.Source.UserID)
	imageUrl := memberPictureURL
	title := "プレゼントがあります！"
	text := memberName + "さん。\n下の中から一つ選んでください。"
	action1 := linebot.NewPostbackAction(
		"A",
		"a",
		"",
		"Aを選択",
	)
	action2 := linebot.NewPostbackAction(
		"B",
		"b",
		"",
		"Bを選択",
	)
	action3 := linebot.NewPostbackAction(
		"C",
		"c",
		"",
		"Cを選択",
	)
	buttonTemplate := linebot.NewButtonsTemplate(
		imageUrl,
		title,
		text,
		action1,
		action2,
		action3,
	)
	replyTemplate := linebot.NewTemplateMessage(
		"誕生日メッセージ",
		buttonTemplate,
	)
	_, err := bot.ReplyMessage(event.ReplyToken, replyTemplate).Do()
	if err != nil {
		log.Print(err)
	}
	updateCounter(event.Source.UserID, "ready")
}

func receivePostback(bot *linebot.Client, event *linebot.Event) {
	// メンバーのカウンターが1の場合のみポストバックを受理する
	// 受理したらカウンターを2にする
	_, counter := getMemberData(event.Source.UserID)
	if counter != "ready" {
		log.Print(counter)
		log.Print("not 1")
		return
	}

	type presentInfo struct {
		URL     string
		Caption string
	}

	presentMap := map[string]*presentInfo{
		"a": {
			URL:     "https://lh3.googleusercontent.com/pw/AM-JKLWNQpQngugw70MsOX0A0g1yFYCkGPU2A7zUKD2cjX-qIf9nLcmV0a-KwaKWNSbKCggyFi1yHE6kJTLok6aYBXoir2RcS59syO68kfG92zNLQ8vhlm26SWAcgt5vvKpiTyFq7seRz1mMPgOmn_oY6pVR=w978-h1442-no?authuser=0",
			Caption: "Aを選んだ君は輝いてるよ！",
		},
		"b": {
			URL:     "https://lh3.googleusercontent.com/pw/AM-JKLUrW-2q1YsVh0Q3-ziYGp4C_Gd9o4DMmQeOA-Acay2Ejh4wVJech5OKGlVOzFM_rziC7DVN8h2E9CqBUIILg163zuAyGVHjAq90cUuonxXUqp34WSFua78h7U9TJkDXDGvKpWZo7gmhWQqRdEaybHWi=s749-no?authuser=0",
			Caption: "Bはからあげクンとハイボールです",
		},
		"c": {
			URL:     "https://lh3.googleusercontent.com/pw/AM-JKLUZHDs546akVKAUOEdHv1K4vApRJBvVYZUCT6ZAutppnVEnX9uvxA7sdkQColNHHmFJmWws2nKUO92dVe2h2WSd1bIgLzEHt5q0CY88EDSxddYDLheRGV3hPQeV0Ad0_QmgopEBiQI-d8IEgw-w_hBL=s750-no?authuser=0",
			Caption: "Cはアイスです",
		},
	}

	// ex) event.Postback.Data is "a"
	mappedPresent, ok := presentMap[event.Postback.Data]
	if !ok {
		return
	}
	originalUrl, previewUrl := mappedPresent.URL, mappedPresent.URL

	newImageMessage := linebot.NewImageMessage(
		originalUrl,
		previewUrl,
	)
	newTextMessage := linebot.NewTextMessage(mappedPresent.Caption)

	_, err := bot.ReplyMessage(event.ReplyToken, newImageMessage, newTextMessage).Do()
	if err != nil {
		log.Print(err)
	}
	updateCounter(event.Source.UserID, "inactive")
}

func getMemberLineInfo(bot *linebot.Client, memberId string) (string, string) {
	res, err := bot.GetGroupMemberProfile(os.Getenv("T2_GROUP_ID"), memberId).Do()
	if err != nil {
		log.Print(err)
	}
	memberName := res.DisplayName
	memberPictureURL := res.PictureURL

	return memberName, memberPictureURL
}

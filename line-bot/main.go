package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Webhook struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	log.Print("初期化します")
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		// Do something when something bad happened.
		log.Print("clientの初期化に失敗")
		log.Print(err)
	}

	log.Print(request.Headers)
	log.Print(request.Body)

	if !validateSignature(os.Getenv("CHANNEL_SECRET"), request.Headers["x-line-signature"], []byte(request.Body)) {
		log.Print("署名失敗")
		return Response{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", linebot.ErrInvalidSignature.Error()),
		}, nil
	}
	log.Print("success sign")

	var webhook Webhook

	if err := json.Unmarshal([]byte(request.Body), &webhook); err != nil {
		return Response{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
		}, nil
	}

	log.Print("webhook")
	for _, event := range webhook.Events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			log.Print("メッセージはテキスト")
			switch event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				log.Print("text")
				// 誕生日メンバーでなければなにもなし
				if event.Source.UserID != os.Getenv("BD_MEMBER_ID") {
					return Response{
						StatusCode: 200,
						Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", "not BDMmember"),
					}, nil
				}
				sendBirthdayButtonsTemplate(bot, event)
			// メッセージが位置情報の場合
			case *linebot.LocationMessage:
				log.Print("location")
				sendlocation(bot, event)
			// メッセージが画像の場合
			case *linebot.ImageMessage:
				log.Print("image")
				sendImage(bot, event)
				// メッセージがポストバックの場合
			}
			// 他にもスタンプや画像、位置情報など色々受信可能
		} else if event.Type == linebot.EventTypePostback {
			receivePostback(bot, event)
		}
	}

	return Response{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

func validateSignature(channelSecret string, signature string, body []byte) bool {
	log.Print("start validation")
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Print("base64")
		return false
	}

	log.Print("hash")
	hash := hmac.New(sha256.New, []byte(channelSecret))
	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	log.Print("Return")
	return hmac.Equal(decoded, hash.Sum(nil))
}

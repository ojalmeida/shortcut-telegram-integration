package telegram

import (
	"bot/config"
	"bot/crud"
	"bot/log"
	"fmt"
	tlsdk "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	notificationRate = time.Hour * 2
)

var (
	unnamedReviewRequestsMessageTemplates = []string{

		"@{requested.username}, solicitaram seu review na estória #{story.number}\n\nLink: {story.url}",
		"@{requested.username}, por bem ou por mal pediram seu review na estória #{story.number}\n\nLink: {story.url}",
		"@{requested.username}, só passando para avisar que seu review na estória #{story.number} foi solicitado...\n\nLink: {story.url}",
		"Senhor(a) @{requested.username}, venho por meio desta informar que seu review na estória #{story.number} foi solicitado 🧐\n\nLink: {story.url}",
	}
	namedReviewRequestsMessageTemplates = []string{

		"@{requested.username}, @{requester.username} solicitou seu review na estória #{story.number}\n\nLink: {story.url}",
		"@{requested.username}, por bem ou por mal, @{requester.username} pediu seu review na estória #{story.number}\n\nLink: {story.url}",
		"@{requested.username}, só passando para avisar que seu review na estória #{story.number} foi solicitado por @{requester.username}...\n\nLink: {story.url}",
		"Senhor(a) @{requested.username}, venho por meio desta informar que seu review na estória #{story.number} foi solicitado por @{requester.username} 🧐\n\nLink: {story.url}",
	}

	doneReviewMessageTemplates = []string{

		"@{requester.username}, o review na estória #{story.number} foi feito :)\n\nLink: {story.url}",
		"@{requester.username}, noticia boa (ou não, vai saber), o review na estória #{story.number} foi feito :)\n\nLink: {story.url}",
		"Senhor(a) @{requester.username}, venho por meio desta informar que o review na estória #{story.number} foi realizado 🧐\n\nLink: {story.url}",
	}
)

func notifyReviewsRequests() {

	for {

		time.Sleep(time.Minute)

		reviews, err := crud.GetReviews()
		if err != nil {
			log.Error(fmt.Sprintf("unable to retrieve reviews from database: %s", err.Error()))
			continue
		}

		for _, review := range reviews {

			if review.NotifyAt < time.Now().Unix() {

				requested, err := crud.GetUser(review.RequestedID)
				if err != nil {
					log.Error(fmt.Sprintf("unable to get user %s: %s", review.RequestedID, err.Error()))
					continue
				}

				if review.SuppressedByRequested {
					continue
				}

				if requested.NotificationRate != 0 {

					review.NotifyAt = time.Now().
						Add(time.Duration(requested.NotificationRate) * time.Second).
						Unix()

				} else {

					review.NotifyAt = time.Now().
						Add(time.Duration(config.Config.Telegram.NotificationRating) * time.Second).
						Unix()

				}

				_, err = crud.UpdateReview(review.ID, review)
				if err != nil {
					log.Error(fmt.Sprintf("unable to update review notification time: %s", err.Error()))
				}

				err = sendReviewRequestNotification(review)
				if err != nil {
					log.Error(fmt.Sprintf("unable to send review notification: %s", err.Error()))
				}

			}

		}

	}

}

func sendReviewRequestNotification(review crud.Review) (err error) {

	requester, err := crud.GetUser(review.RequesterID)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s requester user: %s", review.ID, err.Error()))
	}

	requested, err := crud.GetUser(review.RequestedID)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s requested user: %s", review.ID, err.Error()))
		return
	}

	if requested.TelegramChatID == 0 {

		return
	}

	var messageText string

	switch {

	// unknown requester
	case requester.ShortcutMentionName == "":
		messageText = unnamedReviewRequestsMessageTemplates[rand.Intn(len(unnamedReviewRequestsMessageTemplates))]

	case requester.ShortcutMentionName != "":

		messageText = namedReviewRequestsMessageTemplates[rand.Intn(len(namedReviewRequestsMessageTemplates))]
	}

	messageText = strings.ReplaceAll(messageText, "{requested.username}", requested.ShortcutMentionName)
	messageText = strings.ReplaceAll(messageText, "{requester.username}", requester.ShortcutMentionName)
	messageText = strings.ReplaceAll(messageText, "{story.number}", strconv.FormatInt(review.Number, 10))
	messageText = strings.ReplaceAll(messageText, "{story.url}", review.URL)

	msg := tlsdk.NewMessage(requested.TelegramChatID, messageText)

	msg.ReplyMarkup = tlsdk.NewInlineKeyboardMarkup(

		tlsdk.NewInlineKeyboardRow(

			tlsdk.NewInlineKeyboardButtonData("Review feito", fmt.Sprintf("%s done", review.ID)),
		),

		tlsdk.NewInlineKeyboardRow(

			tlsdk.NewInlineKeyboardButtonData("Não notificar", fmt.Sprintf("%s suppress", review.ID)),
		),
	)

	_, err = bot.Send(msg)
	if err != nil {
		log.Error(fmt.Sprintf("unable to send review notification: %s", err.Error()))
	}

	return

}

func sendReviewDoneNotification(review crud.Review) (err error) {

	requester, err := crud.GetUser(review.RequesterID)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s requester user: %s", review.ID, err.Error()))
	}

	if requester.TelegramChatID == 0 {

		return
	}

	var messageText string

	messageText = doneReviewMessageTemplates[rand.Intn(len(doneReviewMessageTemplates))]

	storyURL, err := url.Parse(review.URL)
	if err != nil {
		log.Error(fmt.Sprintf("unable to parse story url %s: %s", storyURL, err.Error()))
		return
	}
	storyURL.Fragment = "" //removes "#foobar..." from end of url

	messageText = strings.ReplaceAll(messageText, "{requester.username}", requester.ShortcutMentionName)
	messageText = strings.ReplaceAll(messageText, "{story.number}", strconv.FormatInt(review.Number, 10))
	messageText = strings.ReplaceAll(messageText, "{story.url}", storyURL.String())

	msg := tlsdk.NewMessage(requester.TelegramChatID, messageText)

	_, err = bot.Send(msg)
	if err != nil {
		log.Error(fmt.Sprintf("unable to send review notification: %s", err.Error()))
	}

	return
}

package telegram

import (
	"bot/config"
	"bot/log"
	"fmt"
	tlsdk "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

var (
	bot                      *tlsdk.BotAPI
	subscribeMessageRegexp   = regexp.MustCompile(`^/subscribe @\w+`)
	updateMessageRegexp      = regexp.MustCompile(`^/update @\w+`)
	unsubscribeMessageRegexp = regexp.MustCompile(`^/unsubscribe`)
	startMessageRegexp       = regexp.MustCompile(`^/start`)
	authMessageRegexp        = regexp.MustCompile(`^/auth [0-9a-z]{32}`)
	doneCallbackRegexp       = regexp.MustCompile(`^[0-9a-z-]{36} done`)
	suppressCallbackRegexp   = regexp.MustCompile(`^[0-9a-z-]{36} suppress`)
)

func Start() {
	var err error
	bot, err = tlsdk.NewBotAPI(config.Config.Telegram.Token)

	if err != nil {
		panic(fmt.Sprintf("unable to set telegram bot: %s", err.Error()))
	}

	err = tlsdk.SetLogger(log.Logger)

	if err != nil {
		panic(fmt.Sprintf("unable to set telegram sdk logger: %s", err.Error()))
	}

	go listenForMessages()
	go notifyReviewsRequests()
}

func listenForMessages() {

	u := tlsdk.NewUpdate(0)
	u.Timeout = 60

	updatesChan := bot.GetUpdatesChan(u)

	for update := range updatesChan { // blocks execution

		if update.CallbackQuery != nil {

			err := routeCallback(update.CallbackQuery)

			if err != nil {
				log.Error(fmt.Sprintf("error when handling button click: %s", err.Error()))
			}

			continue

		}

		if update.Message != nil {

			log.Debug(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))

			msg := tlsdk.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			err := routeMessage(msg)
			if err != nil {
				log.Error(fmt.Sprintf("error when handling message %s: %s", msg.Text, err.Error()))
				msg.Text = "Isso é embaraçoso, me perdi pensando na sua resposta, alguma coisa deu errado..."
			}
		}
	}

}

func routeMessage(msg tlsdk.MessageConfig) (err error) {

	switch {

	case subscribeMessageRegexp.MatchString(msg.Text):

		err = subscribeMessageHandler(msg)

	case unsubscribeMessageRegexp.MatchString(msg.Text):

		err = unsubscribeMessageHandler(msg)

	case updateMessageRegexp.MatchString(msg.Text):

		err = updateMessageHandler(msg)

	case startMessageRegexp.MatchString(msg.Text):

		err = startMessageHandler(msg)

	case authMessageRegexp.MatchString(msg.Text):

		err = authMessageHandler(msg)

	default:

		err = unknownMessageHandler(msg)

	}

	return

}

func routeCallback(query *tlsdk.CallbackQuery) (err error) {

	switch {

	case doneCallbackRegexp.MatchString(query.Data):

		err = doneCallbackHandler(*query)

	case suppressCallbackRegexp.MatchString(query.Data):

		err = suppressCallbackHandler(*query)

	}

	return
}

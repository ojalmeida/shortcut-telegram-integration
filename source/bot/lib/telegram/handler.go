package telegram

import (
	"bot/config"
	"bot/crud"
	"bot/lib/shortcut"
	"bot/log"
	"errors"
	"fmt"
	tlsdk "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"strings"
)

func subscribeMessageHandler(msg tlsdk.MessageConfig) (err error) {

	var reply tlsdk.MessageConfig
	reply.ChatID = msg.ChatID
	reply.ReplyToMessageID = msg.ReplyToMessageID

	defer func() {

		_, err = bot.Send(reply)

	}()

	if !isAuthorized(msg.ChatID) {

		reply.Text = "Você não está autorizado a se inscrever"
		return

	}

	toCreate := false
	toUpdate := false

	mentionName := strings.ReplaceAll(strings.Split(msg.Text, " ")[1], "@", "")

	user := crud.User{
		TelegramChatID:      msg.ChatID,
		ShortcutMentionName: mentionName,
	}

	var retrieved crud.User
	retrieved, err = crud.GetUserByMentionName(user.ShortcutMentionName)
	if err != nil {

		if errors.Is(err, crud.ErrUserNotFound) {

			toCreate = true

		} else {

			err = ErrUnknown
			return

		}

	} else {
		toUpdate = true
	}

	if toCreate {

		user.ID, err = shortcut.GetIDByMentionName(user.ShortcutMentionName)

		_, err = crud.CreateUser(user)
		if err != nil {

			if errors.Is(err, crud.ErrMentionNameAlreadyExists) {

				err = crud.ErrMentionNameAlreadyExists

			} else {

				err = ErrUnknown

			}

			return

		}

	} else if toUpdate {

		user.ShortcutMentionName = retrieved.ShortcutMentionName
		user.ID = retrieved.ID

		_, err = crud.UpdateUser(user.ID, user)
		if err != nil {

			if errors.Is(err, crud.ErrMentionNameAlreadyExists) {

				err = crud.ErrMentionNameAlreadyExists

			} else {

				err = ErrUnknown

			}

			return

		}

	}

	reply.Text = "Tudo bem então, anotei seu nome agenda :)"

	return

}

func unsubscribeMessageHandler(msg tlsdk.MessageConfig) (err error) {

	var reply tlsdk.MessageConfig
	reply.ChatID = msg.ChatID
	reply.ReplyToMessageID = msg.ReplyToMessageID

	defer func() {

		_, err = bot.Send(reply)

	}()

	if !isAuthorized(msg.ChatID) {

		reply.Text = "Você não está autorizado a se desinscrever"
		return

	}

	chatID := msg.ChatID

	retrieved, err := crud.GetUserByChatID(strconv.FormatInt(chatID, 10))
	if err != nil {

		if errors.Is(err, crud.ErrUserNotFound) {

			err = crud.ErrUserNotFound

		} else {

			err = ErrUnknown

		}

		return

	}

	err = crud.DeleteUser(retrieved.ID)
	if err != nil {

		if errors.Is(err, crud.ErrUserNotFound) {

			err = crud.ErrUserNotFound

		} else {

			err = ErrUnknown

		}

		return

	}

	reply.Text = "Tudo bem, não te enviarei mais mensagens"

	return
}

func updateMessageHandler(msg tlsdk.MessageConfig) (err error) {

	var reply tlsdk.MessageConfig
	reply.ChatID = msg.ChatID
	reply.ReplyToMessageID = msg.ReplyToMessageID

	defer func() {

		_, err = bot.Send(reply)

	}()

	if !isAuthorized(msg.ChatID) {

		reply.Text = "Você não está autorizado a atualizar suas configurações"
		return

	}

	chatID := msg.ChatID
	username := strings.Split(msg.Text, " ")[1]

	_, err = crud.GetUserByMentionName(username)
	if err == nil {
		err = crud.ErrMentionNameAlreadyExists
		return
	}

	if !errors.Is(err, crud.ErrUserNotFound) {

		err = ErrUnknown
		return

	}

	user := crud.User{
		TelegramChatID:      chatID,
		ShortcutMentionName: username,
	}

	_, err = crud.CreateUser(user)
	if err != nil {

		if errors.Is(err, crud.ErrMentionNameAlreadyExists) {

			err = crud.ErrMentionNameAlreadyExists

		} else {

			err = ErrUnknown

		}

		return

	}

	reply.Text = "Oh, esse é seu novo username então ? Tudo bem!"

	return

}

func startMessageHandler(msg tlsdk.MessageConfig) (err error) {

	var reply tlsdk.PhotoConfig
	reply.ChatID = msg.ChatID

	defer func() {

		_, err = bot.Send(reply)

	}()

	reply.File = tlsdk.FileURL("https://static.wikia.nocookie.net/theoffice/images/7/75/YoungMichaelScott.jpg/revision/latest/scale-to-width-down/680?cb=20200413232331")
	reply.Caption = "Pois não, no que posso ajudar ?\n\n\n" +
		"/auth {token} -- primeiro de tudo, até o momento você é um completo estranho, forneça o token disponibilizado pelo administrador para usar o serviço\n\n" +
		"/subscribe @username (e.g. @jperlin) -- se cadastra para receber notificações de reviews requisitadas para @xpto\n\n" +
		"/unsubscribe -- é triste, mas esse comando retira você da minha agenda (caso você não queira receber mais nenhuma notificação)\n\n" +
		"/update @new_username -- atualiza o seu username (útil quando você resolve trocar seu username no Shortcut)"

	return

}

func authMessageHandler(msg tlsdk.MessageConfig) (err error) {

	var reply tlsdk.MessageConfig
	reply.ChatID = msg.ChatID
	reply.ReplyToMessageID = msg.ReplyToMessageID

	defer func() {

		if err == ErrUnknown {

			reply.Text = "Oh, algo deu errado..."

		}

		_, err = bot.Send(reply)

	}()

	if isAuthorized(msg.ChatID) {

		reply.Text = "Você já se autorizou anteriormente"
		return

	}

	if authMessageRegexp.MatchString(msg.Text) {

		token := strings.Split(msg.Text, " ")[1]

		if token == config.Config.Telegram.AuthorizationToken {

			err = authorize(msg.ChatID)
			if err != nil {
				err = ErrUnknown
				return
			}

			reply.Text = "Boa, agora podemos continuar!"

		} else {

			reply.Text = "Esse não é o token de autorização, contate o administrador"

		}

	}

	return

}

func unknownMessageHandler(msg tlsdk.MessageConfig) (err error) {

	var reply tlsdk.MessageConfig
	reply.ChatID = msg.ChatID
	reply.ReplyToMessageID = msg.ReplyToMessageID

	defer func() {

		_, err = bot.Send(reply)

	}()

	reply.Text = "Desculpe, não entendi..."

	return
}

func doneCallbackHandler(query tlsdk.CallbackQuery) (err error) {

	id := regexp.MustCompile(`[0-9a-z-]{36}`).FindString(query.Data)

	retrieved, err := crud.GetReview(id)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s: %s", id, err.Error()))
		err = ErrUnknown
		return
	}

	err = crud.DeleteReview(id)
	if err != nil {
		log.Error(fmt.Sprintf("unable to delete review %s: %s", id, err.Error()))
		err = ErrUnknown
		return
	}

	err = sendReviewDoneNotification(retrieved)

	return
}

func suppressCallbackHandler(query tlsdk.CallbackQuery) (err error) {

	id := regexp.MustCompile(`[0-9a-z-]{36}`).FindString(query.Data)
	chatID := query.Message.Chat.ID

	review, err := crud.GetReview(id)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s: %s", id, err.Error()))
		err = ErrUnknown
		return
	}

	requester, err := crud.GetUser(review.RequesterID)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s requester user: %s", review.ID, err.Error()))
	}

	requested, err := crud.GetUser(review.RequestedID)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get review %s requested user: %s", review.ID, err.Error()))
		return
	}

	if chatID == requester.TelegramChatID {
		review.SuppressedByRequester = true
	}

	if chatID == requested.TelegramChatID {
		review.SuppressedByRequested = true
	}

	_, err = crud.UpdateReview(id, review)
	if err != nil {
		log.Error(fmt.Sprintf("unable to update review %s: %s", id, err.Error()))
		err = ErrUnknown
		return
	}

	return
}

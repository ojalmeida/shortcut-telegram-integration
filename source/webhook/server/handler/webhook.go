package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"webhook/crud"
	"webhook/lib/shortcut"
	"webhook/log"
)

var (
	reviewRequestRegexp = regexp.MustCompile(`.*\breview\b.*`)
)

func WebhookHandler(ctx *fiber.Ctx) (err error) {

	var (
		status          = http.StatusOK
		reviewsToCreate []crud.Review
	)

	defer func() {

		err = nil // Shortcut's webhook does have non-OK responses threshold
		ctx.Status(status)

	}()

	shortcutReq := shortcutRequest{}

	err = json.Unmarshal(ctx.Body(), &shortcutReq)
	if err != nil {
		log.Error(fmt.Sprintf("unable to unmarshal shortcut request into JSON: %s", err.Error()))
		return
	}

	for _, action := range shortcutReq.Actions {

		// only handle story comments creations
		if action.Action == "create" && action.EntityType == "story-comment" && reviewRequestRegexp.MatchString(action.Text) {

			// url will be used in notification
			parsedURL, parseErr := url.Parse(action.AppUrl)
			if parseErr != nil {
				log.Error(fmt.Sprintf("unable to parse AppUrl: %s", parseErr.Error()))
				return
			}
			parsedURL.Fragment = ""

			// story number is retrieved from url path and will be used in notification
			storyNumber, intParseErr := strconv.ParseInt(strings.Split(parsedURL.Path, "/")[3], 10, 64)
			if intParseErr != nil {
				log.Error(fmt.Sprintf("unable to parse story number to int: %s", err.Error()))
				return
			}

			review := crud.Review{
				URL:      parsedURL.String(),
				Number:   storyNumber,
				NotifyAt: time.Now().Unix(),
			}

			// try to get user of specified id
			requester, getErr := crud.GetUser(action.AuthorId)

			// if there is no user, get mention name via shortcut's API and creates it
			if errors.Is(getErr, crud.ErrUserNotFound) {

				mentionName, err2 := shortcut.GetMentionNameByID(action.AuthorId)
				if err2 != nil {
					log.Error(fmt.Sprintf("unable to get mention name by id %s: %s", action.AuthorId, getErr.Error()))
					return
				}

				user := crud.User{
					ID:                  action.AuthorId,
					ShortcutMentionName: mentionName,
				}

				var createErr error
				requester, createErr = crud.CreateUser(user)
				if createErr != nil {
					log.Error(fmt.Sprintf("unable to create user: %s", createErr.Error()))
					return
				}

			} else if getErr != nil {
				log.Error(fmt.Sprintf("unable to get user %s: %s", action.AuthorId, getErr.Error()))
				return
			}

			review.RequesterID = requester.ID

			// more than one person can be mentioned in comment, create a review entity for each
			for _, mentionedID := range action.MentionIds {

				// try to get user of specified id
				requested, getErr := crud.GetUser(mentionedID)

				if errors.Is(getErr, crud.ErrUserNotFound) {

					continue // only queue reviews when reviewer is subscribed to service

				} else if getErr != nil {
					log.Error(fmt.Sprintf("unable to get user %s: %s", mentionedID, getErr.Error()))
					return
				}

				review.RequestedID = requested.ID
				reviewsToCreate = append(reviewsToCreate, review)

			}

		}

	}

	for _, review := range reviewsToCreate {

		_, err = crud.CreateReview(review)
		if err != nil {
			log.Error(fmt.Sprintf("unable to create review: %s", err.Error()))
		}

	}

	return

}

type shortcutRequest struct {
	Id        string    `json:"id,omitempty"`
	ChangedAt time.Time `json:"changed_at,omitempty"`
	PrimaryId int       `json:"primary_id,omitempty"`
	MemberId  string    `json:"member_id,omitempty"`
	Version   string    `json:"version,omitempty"`
	Actions   []struct {
		Id         int      `json:"id,omitempty"`
		EntityType string   `json:"entity_type,omitempty"`
		Action     string   `json:"action,omitempty"`
		AuthorId   string   `json:"author_id,omitempty"`
		AppUrl     string   `json:"app_url,omitempty"`
		Text       string   `json:"text,omitempty"`
		MentionIds []string `json:"mention_ids,omitempty"`
		Name       string   `json:"name,omitempty"`
		StoryType  string   `json:"story_type,omitempty"`
		Changes    struct {
			Started struct {
				New bool `json:"new,omitempty"`
				Old bool `json:"old,omitempty"`
			} `json:"started,omitempty"`
			WorkflowStateId struct {
				New int `json:"new,omitempty"`
				Old int `json:"old,omitempty"`
			} `json:"workflow_state_id,omitempty"`
			OwnerIds struct {
				Adds []string `json:"adds,omitempty"`
			} `json:"owner_ids,omitempty"`
		} `json:"changes,omitempty"`
	} `json:"actions,omitempty"`
	References []struct {
		Id         int    `json:"id,omitempty"`
		EntityType string `json:"entity_type,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"references,omitempty"`
}

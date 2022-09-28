package shortcut

import (
	"bot/config"
	"bot/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	baseURL            = "https://api.app.shortcut.com/api/v3"
	apiTokenHeaderName = "Shortcut-Token"
	errNonOKResponse   = errors.New("shortcut API non-OK response")
	ErrProfileNotFound = errors.New("profile not found")
)

func GetMentionNameByID(id string) (mentionName string, err error) {

	var shortcutJSONResponse struct {
		CreatedAt  time.Time `json:"created_at,omitempty"`
		Disabled   bool      `json:"disabled,omitempty"`
		EntityType string    `json:"entity_type,omitempty"`
		GroupIds   []string  `json:"group_ids,omitempty"`
		Id         string    `json:"id,omitempty"`
		Profile    struct {
			Deactivated bool `json:"deactivated,omitempty"`
			DisplayIcon struct {
				CreatedAt  time.Time `json:"created_at,omitempty"`
				EntityType string    `json:"entity_type,omitempty"`
				Id         string    `json:"id,omitempty"`
				UpdatedAt  time.Time `json:"updated_at,omitempty"`
				Url        string    `json:"url,omitempty"`
			} `json:"display_icon,omitempty"`
			EmailAddress           string `json:"email_address,omitempty"`
			EntityType             string `json:"entity_type,omitempty"`
			GravatarHash           string `json:"gravatar_hash,omitempty"`
			Id                     string `json:"id,omitempty"`
			MentionName            string `json:"mention_name,omitempty"`
			Name                   string `json:"name,omitempty"`
			TwoFactorAuthActivated bool   `json:"two_factor_auth_activated,omitempty"`
		} `json:"profile,omitempty"`
		Role      string    `json:"role,omitempty"`
		State     string    `json:"state,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/members/%s", baseURL, id),
		nil,
	)

	if err != nil {
		log.Error(fmt.Sprintf("unable to create request: %s", err.Error()))
		return
	}

	req.Header.Set(apiTokenHeaderName, config.Config.Shortcut.Token)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Error(fmt.Sprintf("unable to make request: %s", err.Error()))
		return
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {

		if res.StatusCode == http.StatusNotFound {

			err = ErrProfileNotFound
			return

		}

		log.Error(fmt.Sprintf("shortcut answered with non-OK status code: %d", res.StatusCode))
		err = errNonOKResponse
		return
	}

	err = json.NewDecoder(res.Body).Decode(&shortcutJSONResponse)
	if err != nil {
		log.Error(fmt.Sprintf("unable to unmarshal response into JSON: %s", err.Error()))
	}

	mentionName = shortcutJSONResponse.Profile.MentionName

	return

}

func GetIDByMentionName(mentionName string) (id string, err error) {

	var shortcutJSONResponse []struct {
		CreatedAt  time.Time `json:"created_at,omitempty"`
		Disabled   bool      `json:"disabled,omitempty"`
		EntityType string    `json:"entity_type,omitempty"`
		GroupIds   []string  `json:"group_ids,omitempty"`
		Id         string    `json:"id,omitempty"`
		Profile    struct {
			Deactivated bool `json:"deactivated,omitempty"`
			DisplayIcon struct {
				CreatedAt  time.Time `json:"created_at,omitempty"`
				EntityType string    `json:"entity_type,omitempty"`
				Id         string    `json:"id,omitempty"`
				UpdatedAt  time.Time `json:"updated_at,omitempty"`
				Url        string    `json:"url,omitempty"`
			} `json:"display_icon,omitempty"`
			EmailAddress           string `json:"email_address,omitempty"`
			EntityType             string `json:"entity_type,omitempty"`
			GravatarHash           string `json:"gravatar_hash,omitempty"`
			Id                     string `json:"id,omitempty"`
			MentionName            string `json:"mention_name,omitempty"`
			Name                   string `json:"name,omitempty"`
			TwoFactorAuthActivated bool   `json:"two_factor_auth_activated,omitempty"`
		} `json:"profile,omitempty"`
		Role      string    `json:"role,omitempty"`
		State     string    `json:"state,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/members", baseURL),
		nil,
	)

	if err != nil {
		log.Error(fmt.Sprintf("unable to create request: %s", err.Error()))
		return
	}

	req.Header.Set(apiTokenHeaderName, config.Config.Shortcut.Token)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Error(fmt.Sprintf("unable to make request: %s", err.Error()))
		return
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {

		log.Error(fmt.Sprintf("shortcut answered with non-OK status code: %d", res.StatusCode))
		err = errNonOKResponse
		return
	}

	err = json.NewDecoder(res.Body).Decode(&shortcutJSONResponse)
	if err != nil {
		log.Error(fmt.Sprintf("unable to unmarshal response into JSON: %s", err.Error()))
	}

	found := false
	for _, member := range shortcutJSONResponse {

		if member.Profile.MentionName == mentionName {

			found = true
			id = member.Id
			break

		}

	}

	if !found {

		err = ErrProfileNotFound

	}

	return
}

package userservice

import (
	"log"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/mattermost/mattermost-server/model"
)

func GetUser(clientId string) ApiResponse {
	url, err := lookupService("workstation-user")
	if err != nil {
		log.Fatalf("Error. %s", err)
		return NewErrorApiResponse(500, err.Error())
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/api/userim/" + clientId)
	if err != nil {
		log.Fatalf("Error. %s", err)
		return NewErrorApiResponse(500, err.Error())
	}
	defer resp.Body.Close()

	result := ApiResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error. %s", err)
		return NewErrorApiResponse(500, err.Error())
	}
	return result
}

func SaveUser(user *model.User) error {
	result := ApiResponse{}
	url, err := lookupService("workstation-user")
	if err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	client := &http.Client{}
	userStr, userErr := json.Marshal(user)
	if userErr == nil {
		userTest := string(userStr)
		resp, err := client.Post(url+"/api/userim/add", "application/json", strings.NewReader(userTest))
		if err != nil {
			log.Fatalf("Error. %s", err)
			return err
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Fatalf("Error. %s", err)
			return err
		}
	}
	return nil
}

func updateUser(user *model.User) error {
	result := ApiResponse{}
	url, err := lookupService("workstation-user")
	if err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	client := &http.Client{}
	userStr, userErr := json.Marshal(user)
	if userErr == nil {
		resp, err := client.Post(url+"/api/userim/update", "application/json", strings.NewReader(string(userStr)))
		if err != nil {
			log.Fatalf("Error. %s", err)
			return err
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Fatalf("Error. %s", err)
			return err
		}
	}
	return nil
}

package userservice

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/mattermost/mattermost-server/model"
)

func SaveOAuthApp(app *model.OAuthApp) error {
	result := ApiResponse{}
	url, err := lookupService("user-service")
	if err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/auth")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	return nil
}

func SaveOAuthAccessData(app *model.AccessData) error {
	result := ApiResponse{}
	url, err := lookupService("user-service")
	if err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/auth")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	return nil
}

func SaveOAuthData(app *model.AuthData) error {
	result := ApiResponse{}
	url, err := lookupService("user-service")
	if err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/auth")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error. %s", err)
		return err
	}
	return nil
}

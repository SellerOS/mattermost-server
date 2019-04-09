package userservice

import (
	"github.com/mattermost/mattermost-server/model"
	"log"
	"net/http"
	"encoding/json"
)

func saveUser(user *model.User) (ApiResponse ) {
	result := ApiResponse{}
	url, err := lookupService("user-service")
	if err != nil {
		log.Fatalln("Error. %s", err)
		return NewErrorApiResponse(500, err.Error())
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/products")
	if err != nil {
		return NewErrorApiResponse(500, err.Error())
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalln("Error. %s", err)
		return NewErrorApiResponse(500, err.Error())
	}
	return result
}


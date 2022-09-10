package main

import (
	"fmt"
	"github.com/nicklaw5/helix/v2"
)

func GetUser(userAccessToken string, clientID string, clientSecret string) (*TwitchUser, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		UserAccessToken: userAccessToken,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.GetUsers(&helix.UsersParams{})

	if len(resp.Data.Users) > 0 {
		var apiUser = resp.Data.Users[0]
		fmt.Println("Logged in as: " + apiUser.DisplayName)
		return &TwitchUser{
			Name:           apiUser.Login,
			ReadableName:   apiUser.DisplayName,
			Description:    apiUser.Description,
			ProfilePicture: apiUser.ProfileImageURL,
			Email:          apiUser.Email,
		}, nil
	}

	return nil, fmt.Errorf("no user found")
}

type TwitchUser struct {
	Name           string
	ID             string
	ReadableName   string
	Description    string
	ProfilePicture string
	Email          string
}

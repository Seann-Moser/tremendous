package main

import (
	"context"
	"fmt"
	"github.com/Seann-Moser/tremendous"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	client := tremendous.NewClient(http.DefaultClient)
	defer client.Close()

	// if enabled and refresh token and client_id, client_secret are provided,
	// it will generate a new access_toke
	autoRefresh := false

	/*
		This was done so you can create child oauth config for each sub account you wish to access
	*/
	oauthClient := client.NewClientWithOAuth(tremendous.OauthConfig{
		ClientId:     "",
		ClientSecret: "",
		AccessToken:  "PROD_6aRVAo5SA--NQTlACyH6wZsxAnEWSVihNM0C1WofJ5G",
	}, autoRefresh)
	go func() {
		for oauth := range oauthClient.OauthRefresh() {
			// this channel will emit changes when we generate a new refresh token
			// should be encrypted and stored somewhere
			// save refresh token else were since refresh token can only be used once
			println(oauth.RefreshToken)
		}
	}()
	members, err := oauthClient.ListMembers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range members.Members {
		fmt.Println(m)
	}

}

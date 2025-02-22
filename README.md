# Tremendous Go API Client

## Reference Docs:
- [Tremendous Docs](https://developers.tremendous.com/reference/rewards)
- [Tremendous Oauth](https://developers.tremendous.com/docs/oauth-2)


### Setup

```go
client := tremendous.NewClient(http.DefaultClient)
defer client.Close()
```

#### Using API Key
```go
client := tremendous.NewClient(http.DefaultClient)
defer client.Close()
apiKeyClient := client.NewClientWithAPIKey("{APIKEY}")
apiKeyClient.ListMembers(context.Background())
```

##### Using Oauth
```go
client := tremendous.NewClient(http.DefaultClient)
defer client.Close()

// if enabled and refresh token and client_id, client_secret are provided,
// it will generate a new access_token
autoRefresh := false

/*
	This was done so you can create child oauth config for each sub account you wish to access
*/

oauthClientAccount1 := client.NewClientWithOAuth(tremendous.OauthConfig{
AccessToken:  "account_token1",
}, autoRefresh)


members, err := oauthClientAccount1.ListMembers(context.Background())

autoRefresh := true
oauthClientAccount2 := client.NewClientWithOAuth(tremendous.OauthConfig{
ClientId:     "test1",
ClientSecret: "test1",
RefreshToken: "refresh_token"
AccessToken:  "accountToken2",
}, autoRefresh)

go func() {
    for oauth := range oauthClientAccount2.OauthRefresh() {
    // this channel will emit changes when we generate a new refresh token
    // should be encrypted and stored somewhere
    // save refresh token else were since refresh token can only be used once
    println(oauth.RefreshToken)
    }
}()

members, err =oauthClientAccount2.ListMembers(context.Background())

```


```go
client := tremendous.NewClient(http.DefaultClient)
defer client.Close()

// Should retrieve CODE from url param on tremendous redirect
token, err := client.SendOauthRequest(ctx, &AccessTokenRequest{
    ClientId:     c.clientID,
    ClientSecret: c.clientSecret,
    GrantType:    GrantTypeAuthorizationCode,
    RedirectUri:"{OAUTH_URI}",
    Code: "{CODE}"
})

// persist token
// create oauth client using client.NewClientWithOAuth()


func handle(w http.ResponseWriter, r *http.Request) {
  client := tremendous.NewClient(http.DefaultClient)
  defer client.Close()
  
  // Should retrieve CODE from url param on tremendous redirect
  token, err := client.SendOauthRequest(ctx, &AccessTokenRequest{
    ClientId:     c.clientID,
    ClientSecret: c.clientSecret,
    GrantType:    GrantTypeAuthorizationCode,
    RedirectUri:  "{OAUTH_URI}",
    Code:         r.URL.Query().Get("code"),
  })
    // persist token
}
```
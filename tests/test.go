package main

import (
	"code.google.com/p/goauth2/oauth"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"fmt"
	// "github.com/eliothedeman/drivefs/fs"
	"log"
)

func main() {

	clientId := "266884823682-qfd1cvgpe2vucokeuhv0qr0bihn87ge4.apps.googleusercontent.com"
	clientSecret := "6SGgwuVL8KMyFsigaPn2MGBK"
	scope := "https://www.googleapis.com/auth/drive"
	redirectURL := "urn:ietf:wg:oauth:2.0:oob"
	authURL := "https://accounts.google.com/o/oauth2/auth"
	tokenURL := "https://accounts.google.com/o/oauth2/token"
	code := "4/1KWWgmhjszhUjor_khaQnI2eMZM1.AvMtkVA0uv0dmmS0T3UFEsPKgYAYjQI"
	cachefile := "cache.json"

	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		TokenCache:   oauth.CacheFile(cachefile),
	}
	transport := &oauth.Transport{Config: config}
	token, err := config.TokenCache.Token()
	if err != nil {
		if code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			return
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}
	transport.Token = token
	client := transport.Client()
	if err != nil {
		log.Fatal("Get", err)
	}
	d, err := drive.New(client)
	if err != nil {
		log.Fatal("Client", err)
	}
	f, e := d.Files.List().Do()
	if e != nil {
		fmt.Println(err)
	}
	fileList := f.Items

	for i := 0; i < len(fileList); i++ {

		fmt.Println(len(fileList[i].Parents), fileList[i].Title)

	}

}

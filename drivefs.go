package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"github.com/eliothedeman/drivefs/drivefs"
	"os"
)

// Collect command line arguments for OAUTH and mounting options
var (
	// OAUTH options
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	scope        = flag.String("scope", "https://www.googleapis.com/auth/drive", "OAuth scope")
	redirectURL  = flag.String("redirect_url", "oob", "Redirect URL")
	authURL      = flag.String("auth_url", "https://accounts.google.com/o/oauth2/auth", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://accounts.google.com/o/oauth2/token", "Token URL")
	requestURL   = flag.String("request_url", "https://www.googleapis.com/oauth2/v1/userinfo", "API request")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
)

func main() {
	flag.Parse()
	if _, err := os.Stat(*cachefile); err != nil {
		f, err := os.Create(*cachefile)
		if err == nil {
			f.Close()
		}
	}
	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		Scope:        *scope,
		AuthURL:      *authURL,
		TokenURL:     *tokenURL,
		TokenCache:   oauth.CacheFile(*cachefile),
	}
	drivefs.NewServer(config, *code)

}

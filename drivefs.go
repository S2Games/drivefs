package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"github.com/eliothedeman/drivefs/drivefs"
	"log"
	"os"
	"time"
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

	// Filesystem options
	mountpoint = flag.String("mount", "", "Mount point for drivefs")
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
	// fail if the server can't auth
	server, err := drivefs.NewServer(config, *code)
	if err != nil {
		log.Fatal(err)
	}
	// Attempt to mount the filesystem, fail if mountpoint is not given
	if *mountpoint == "" {
		log.Fatal("Must provide mountpoint via argument -mount")
	}
	err = server.Mount(*mountpoint)
	if err != nil {
		log.Println(err)
	}
	go server.Serve()
	time.Sleep(10 * time.Second)
	server.Unmount(*mountpoint, 3)

}

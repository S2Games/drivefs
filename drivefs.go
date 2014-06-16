package main

import (
	"bazil.org/fuse/fs/fstestutil"
	"code.google.com/p/goauth2/oauth"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"flag"
	"github.com/eliothedeman/drivefs/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Collect command line arguments for OAUTH and mounting options
var (
	// OAUTH options
	clientId     = flag.String("id", "266884823682-qfd1cvgpe2vucokeuhv0qr0bihn87ge4.apps.googleusercontent.com", "Client ID")
	clientSecret = flag.String("secret", "6SGgwuVL8KMyFsigaPn2MGBK", "Client Secret")
	scope        = flag.String("scope", "https://www.googleapis.com/auth/drive", "OAuth scope")
	redirectURL  = flag.String("redirect_url", "oob", "Redirect URL")
	authURL      = flag.String("auth_url", "https://accounts.google.com/o/oauth2/auth", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://accounts.google.com/o/oauth2/token", "Token URL")
	requestURL   = flag.String("request_url", "https://www.googleapis.com/oauth2/v1/userinfo", "API request")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
	// Filesystem options
	mountpoint = flag.String("mount", "", "Mount point for drivefs")
	refresh    = flag.Int("refresh", 10, "Rate at which to refresh if local file system has not changed.")
)

// Exists checks if a file or directory exists on disk
func Exists(fileName string) bool {
	if a, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	} else {
		log.Println(a)
		return true
	}
}
func main() {
	fstestutil.DebugByDefault()
	flag.Parse()
	// if the cache file does not exists, create it
	if !Exists(*cachefile) {
		f, err := os.Create(*cachefile)
		if err == nil {
			f.Close()
		} else {
			log.Fatal("CacheFile: ", err)
		}
	}
	// if the mountpoint given does not exist, make it
	if *mountpoint == "" {
		log.Fatal("Must provide mountpoint via argument -mount")
	}
	if !Exists(*mountpoint) {
		log.Println(*mountpoint)
		err := os.Mkdir(*mountpoint, 0777)
		if err != nil {
			log.Fatal("Mountpoint: ", err)
		}
	}
	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		Scope:        drive.DriveScope,
		AuthURL:      *authURL,
		TokenURL:     *tokenURL,
		TokenCache:   oauth.CacheFile(*cachefile),
	}
	// fail if the server can't auth
	server, err := lib.NewServer(config, *code)
	if err != nil {
		log.Fatal(err)
	}
	// Attempt to mount the filesystem, fail if mountpoint is not given
	err = server.Mount(*mountpoint)
	if err != nil {
		log.Fatalln(err)
	}
	// Start the server
	go server.Serve(*refresh)

	// wait for a termination before exit
	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Interrupt)
	signal.Notify(killChan, syscall.SIGTERM)
	for sig := range killChan {
		log.Println("drivefs: stopping due to ", sig)
		break
	}
	server.Unmount(*mountpoint, 3)

}

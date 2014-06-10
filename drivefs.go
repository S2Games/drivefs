package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api/go-client/drive/v2"
	"flag"
	drivefs "github.com/eliothedeman/drivefs/fs"
)

var (
	clientId     = flag.String("id", "", "OAUTH client id")
	clientSecret = flag.String("secret", "", "OAUTH client secret")
	scope        = flag.String("scope", "", "https://www.googleapis.com/auth/buzz")
)

func main() {

}

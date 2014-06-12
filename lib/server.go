package lib

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"code.google.com/p/goauth2/oauth"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"log"
	"time"
)

// DriveFs is a struct which holds the FUSE filesystem
type Server struct {
	Config *oauth.Config
	conn   *fuse.Conn
}

// initialize and return a new DriveFS
func NewServer(config *oauth.Config, code string) (*Server, error) {
	d := &Server{Config: config}
	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if d.Config.ClientId == "" || d.Config.ClientSecret == "" {
			log.Fatal("Both Client Id and Client Secret must be given")
		}
		if code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := d.Config.AuthCodeURL("")
			log.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			log.Fatalln(url)
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		log.Printf("Token is cached in %v\n", d.Config.TokenCache)
	}
	transport.Token = token

	// initialize global vars
	client = transport.Client()
	nameToDir = make(map[string]DriveDir)
	nameToFile = make(map[string]DriveFile)
	idToDir = make(map[string]DriveDir)
	idToFile = make(map[string]DriveFile)
	service, err = drive.New(client)
	return d, err
}

// Mount attmpts to mount the filesystem
func (s *Server) Mount(mountPoint string) (err error) {
	s.conn, err = fuse.Mount(mountPoint)
	return
}

// Serve attempts to serve the filesystem
func (s *Server) Serve(refreshRate int) {
	refreshFileIndex()
	refreshChildIndex()
	refreshNameToFile()
	refreshNameToDir()
	// refresh the fileIndex and DirIndex every ten seconds
	go func() {
		for {
			<-time.After(time.Duration(refreshRate) * time.Second)
			log.Println("Refreshing indexes.")
			refreshFileIndex()
			go refreshChildIndex()
			go refreshNameToFile()
			go refreshNameToDir()
		}
	}()
	if err := fs.Serve(s.conn, Root{}); err != nil {
		log.Fatalf("Could not serve drivefs %s", err.Error())
	}
}

// Unmount attempts to unmount the filesystem
func (s *Server) Unmount(mountPoint string, timeout int) (err error) {
	for i := 0; i < timeout; i++ {
		err = fuse.Unmount(mountPoint)
		if err == nil {
			err = s.conn.Close()
			if err != nil {
				log.Fatalln(err)
			}
			return nil
		}
		log.Println(err)
	}
	err = s.conn.Close()
	return err

}

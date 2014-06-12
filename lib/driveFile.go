package lib

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"io/ioutil"
	"log"
	"time"
)

// DriveFile represents a file in google drive
type DriveFile struct {
	File     *drive.File
	Modified time.Time
	Created  time.Time
	Root     bool
}

// Attr returns the file attributes
func (d DriveFile) Attr() fuse.Attr {
	return fuse.Attr{
		Mode:  0777,
		Mtime: time.Now(),
		Size:  uint64(d.File.FileSize),
	}
}

// ReadAll reads an entire file from google drive and returns the resulting bytes
func (d DriveFile) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	byteChan := make(chan *[]byte)
	errChan := make(chan error)
	defer func() {
		close(byteChan)
		close(errChan)
	}()
	// launch read goroutine
	go func() {
		// grab file from google drive api
		log.Println(d.File.Title)
		b, err := client.Get(d.File.DownloadUrl)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}
		// read the data from body and close connection
		c, err := ioutil.ReadAll(b.Body)
		if err != nil {
			log.Println(err)
			errChan <- err
		}
		byteChan <- &c
	}()
	// wait for real to be done, or for file system interupt and return values
	select {
	case tmp := <-byteChan:
		return *tmp, nil
	case <-intr:
		return nil, nil
	}

}

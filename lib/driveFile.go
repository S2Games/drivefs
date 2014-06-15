package lib

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// DriveFile represents a file in google drive
type DriveFile struct {
	File     *drive.File
	TmpFile  *os.File
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
	// wait for read to be done, or for file system interupt and return values
	select {
	case tmp := <-byteChan:
		return *tmp, nil
	case <-intr:
		return nil, nil
	}

}

// Setattr sets file attributs
func (d *DriveFile) Setattr(req *fuse.SetattrRequest, resp *fuse.SetattrResponse, intr fs.Intr) fuse.Error {
	valid := req.Valid
	if valid.Size() {
		return nil
	}
	return nil
}

// Write writes bytes to a tmp file, which are then synced by FSync
func (d *DriveFile) Write(req *fuse.WriteRequest, resp *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	log.Println("hereererrrrrrrr\n\n\n")
	var size int
	// check if d already has a tmp file
	if path, ok := idToTmpFile[d.File.Id]; ok {
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			return err
		}
		size, err = f.Write(req.Data)
		if err != nil {
			log.Println(err)
			return err
		}
		return f.Close()
	}
	// If d does not have a tmp file, create one and write to it
	path := "/tmp/" + d.File.Title
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}
	size, err = f.Write(req.Data)
	if err != nil {
		log.Println(err)
		return err
	}
	// add d's tmp file to the lookup map
	idToTmpFile[d.File.Id] = path
	resp.Size = size
	return f.Close()

}

// Open a file or directory
func (d *DriveFile) Open(req *fuse.OpenRequest, resp *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	log.Println("hereererrrrrrrr\n\n\n")
	// If d does not have a tmp file, create one and write to it
	path := "/tmp/" + d.File.Id
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// add d's tmp file to the lookup map
	idToTmpFile[d.File.Id] = path
	resp.Flags &^= fuse.OpenDirectIO
	return d, f.Close()
}

// Flush a file tmp to google drive
func (d *DriveFile) Flush(req *fuse.FlushRequest, intr fs.Intr) fuse.Error {
	log.Println("Flushing: ", d.TmpFile.Name())
	// sync file to disk
	err := d.TmpFile.Sync()
	if err != nil {
		log.Println(err)
		return err
	}
	// upload file to google drive
	d.File, err = service.Files.Update(d.File.Id, d.File).Media(d.TmpFile).Do()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

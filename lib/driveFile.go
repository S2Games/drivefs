package lib

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// DriveFile represents a file in google drive
type DriveFile struct {
	File     *drive.File
	TmpFile  *os.File
	Modified time.Time
	Created  time.Time
	Root     bool
	*sync.Mutex
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
	d.Lock()
	defer d.Unlock()
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
	d.Lock()
	defer d.Unlock()
	// check if d already has a tmp file
	// If d does not have a tmp file, create one and write to it
	size, err := d.TmpFile.Write(req.Data)
	if err != nil {
		log.Println(err)
		return err
	}
	resp.Size = size
	return nil

}

// Open a file or directory
func (d *DriveFile) Open(req *fuse.OpenRequest, resp *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	d.Lock()
	defer d.Unlock()
	d.TmpFile.Close()
	f, err := os.Create("/tmp/drivefs-" + d.File.Id)
	if err != nil {
		return nil, err
	}
	d.TmpFile = f

	resp.Flags &^= fuse.OpenDirectIO
	return d, nil
}

// Flush a file tmp to google drive
func (d *DriveFile) Flush(req *fuse.FlushRequest, intr fs.Intr) fuse.Error {
	d.Lock()
	defer d.Unlock()
	// sync file to disk
	err := d.TmpFile.Sync()
	if err != nil {
		log.Println(err)
		return err
	}
	d.TmpFile.Close()
	d.TmpFile, err = os.Open(d.TmpFile.Name())
	if err != nil {
		log.Println(err)
	}
	// upload file to google drive
	// this is done in another go routine to catch interupts
	errChan := make(chan error)
	fileChan := make(chan *drive.File)
	go func() {
		f, err := service.Files.Update(d.File.Id, d.File).Media(d.TmpFile).Do()
		if err != nil {
			errChan <- err
			return
		}
		fileChan <- f
		return
	}()
	// wait for interupt while uploading file
	select {
	// if all goes well, set new file and return
	case f := <-fileChan:
		d.File = f
		if d.TmpFile != nil && exists(d.TmpFile.Name()) {
			return os.Remove(d.TmpFile.Name())
		}
		return nil
	case err := <-errChan:
		return err
	// catch interupt
	case <-intr:
		return fuse.EINTR
	}

}

// exists checks if a file or directory exists on disk
func exists(fileName string) bool {
	if a, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	} else {
		log.Println(a)
		return true
	}
}

// Lock locks a *DriveFile's mutex
func (d *DriveFile) Lock() {
	d.Mutex.Lock()
}

// Unlock unlocks a *DriveFile's mutex
func (d *DriveFile) Unlock() {
	d.Mutex.Unlock()
}

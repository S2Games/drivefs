package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// DriveDir represents a directory in google drive
type DriveDir struct {
	Dir      *drive.File
	Modified time.Time
	Created  time.Time
	Root     bool
}

// Attr returns the file attributes
func (d *DriveDir) Attr() fuse.Attr {
	return fuse.Attr{
		Mode: 0644,
	}
}

// TODO implement create function to actually create file
func (DriveDir) Create(req *fuse.CreateRequest, res *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
	return nil, nil, fuse.Errno(syscall.EROFS)
}

// TODO implement fsync function to actually perform an fsync
func (DriveDir) Fsync(req *fuse.FsyncRequest, intr fs.Intr) fuse.Error {
	return fuse.Errno(syscall.EROFS)
}

// TODO implement link function to actually perform a link
func (DriveDir) Link(req *fuse.LinkRequest, node fs.Node, intr fs.Intr) (fs.Node, fuse.Error) {
	return nil, fuse.Errno(syscall.EROFS)
}

// Lookup scans the current directory for matching files or directories
func (d *DriveDir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	// Lookup dir by name
	if dir, ok := nameToDir[name]; ok {
		return dir, nil
	}

	// Lookup file by name
	if file, ok := nameToDir[name]; ok {
		return file, nil
	}

	// File not found
	return nil, fuse.ENOENT
}

// ReadDir return a slice of directory entries
func (d *DriveDir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	// List of directories to return
	var dirs []fuse.Dirent

	// get all new list of files
	f, err := service.Files.List().Do()
	if err != nil {
		return nil, err
	}
	fileList := f.Items
	// Populate idToFile with new ids
	for i := range fileList {
		idToFile[fileList[i].Id] = fileList[i]
	}
	// get list of children
	c, err := service.Children.List(d.Dir.Id).Do()
	// Get children of this folder
	children := c.Items

	dirs = make([]fuse.Dirent, len(children))

	// populate dirs with children
	for i := range children {
		// pull out a child temporarally
		tmp := idToFile[children[i].Id]
		// If child is a folder/directory create a DirveDir else create a DriveFile
		if strings.Contains(tmp.File.MimeType, "folder") {
			dirs[i] = fuse.Dirent{
				Name: tmp.File.Title,
				Type: fuse.DT_Dir,
			}

		} else {
			dirs[i] = fuse.Dirent{
				Name: tmp.File.Title,
				Type: fuse.DT_File,
			}
		}

	}

}

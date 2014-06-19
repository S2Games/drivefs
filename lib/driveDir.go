package lib

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
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
func (DriveDir) Attr() fuse.Attr {
	return fuse.Attr{
		Mode: os.ModeDir | 0555,
	}
}

// Create creates an empty file inside of d
func (d *DriveDir) Create(req *fuse.CreateRequest, res *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
	newFile := &drive.File{}
	newFile.Title = req.Name
	p := &drive.ParentReference{Id: d.Dir.Id}
	newFile.Parents = []*drive.ParentReference{p}
	// create temporary file to serve as the cache until the data is uploaded
	path := "/tmp/drivefs-" + req.Name
	tmpFile, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	createdFile, err1 := service.Files.Insert(newFile).Media(tmpFile).Do()
	if err1 != nil {
		log.Println(err1)
		return nil, nil, err1
	}
	os.Remove(path)
	tmpFile, err = os.Create("/tmp/drivefs-" + createdFile.Id)
	f := &DriveFile{File: createdFile, Root: false, TmpFile: tmpFile, Mutex: new(sync.Mutex)}
	// add the new file to the cach/index
	nameToFile[f.File.Title] = f
	idToFile[f.File.Id] = f
	fileIndex[f.File.Id] = f.File

	return f, f, nil
}

// Lookup scans the current directory for matching files or directories
func (d *DriveDir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	// Lookup dir by name
	if dir, ok := nameToDir[name]; ok {
		return dir, nil
	}

	// Lookup file by name
	if file, ok := nameToFile[name]; ok {
		return file, nil
	}
	// This comes up as the node id for first access, so just show the root folder
	if name == ".xdg-volume-info" {
		if dir, ok := nameToDir["root"]; ok {
			return dir, nil
		}
	}
	// File not found
	return nil, fuse.ENOENT
}

// ReadDir return a slice of directory entries
func (d *DriveDir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	dirChan := make(chan *[]fuse.Dirent)
	errChan := make(chan error)
	defer func() {
		close(dirChan)
		close(errChan)
	}()
	go func() {
		// List of directories to return
		var dirs []fuse.Dirent
		tmpDirMap := make(map[string]fuse.Dirent)
		// get list of children
		// If d is at root, fetch the root children, else fetch this file's children
		var c *drive.ChildList
		if d.Root {
			c = childIndex["root"]
		} else {
			c = childIndex[d.Dir.Id]
		}
		// Get children of this folder
		children := c.Items
		// populate dirs with children
		for i := range children {
			// pull out a child temporarally
			tmp, ok := idToFile[children[i].Id]
			// only add if you have access to the file
			if ok {
				// If child is a folder/directory create a DirveDir else create a DriveFile
				if strings.Contains(tmp.File.MimeType, "folder") {
					tmpDirMap[tmp.File.Id] = fuse.Dirent{
						Name: tmp.File.Title,
						Type: fuse.DT_Dir,
					}
				} else {
					tmpDirMap[tmp.File.Id] = fuse.Dirent{
						Name: tmp.File.Title,
						Type: fuse.DT_File,
					}
				}
			}

		}
		// collaps map to a slice and return
		dirs = make([]fuse.Dirent, len(tmpDirMap))
		i := 0
		for _, v := range tmpDirMap {
			dirs[i] = v
			i += 1
		}
		dirChan <- &dirs
	}()
	// Wait for the lookups to be done, or die if interupt happens
	select {
	case tmp := <-dirChan:
		return *tmp, nil
	case err := <-errChan:
		return nil, err
	case <-intr:
		return nil, fuse.EINTR
	}

}

// Mkdir registers a new directory
func (d *DriveDir) Mkdir(req *fuse.MkdirRequest, intr fs.Intr) (fs.Node, fuse.Error) {
	f := &drive.File{Title: req.Name, MimeType: "application/vnd.google-apps.folder"}
	newDir, err := service.Files.Insert(f).Do()
	if err != nil {
		log.Println(err)
		return nil, fuse.Errno(syscall.EROFS)
	}
	return DriveDir{Dir: newDir, Root: false}, nil
}

// Remove deletes a fild or folder from google drive
func (d *DriveDir) Remove(req *fuse.RemoveRequest, intr fs.Intr) fuse.Error {
	if file, ok := nameToFile[req.Name]; ok {
		err := service.Files.Delete(file.File.Id).Do()
		if err != nil {
			log.Println(err)
		}
		return err
	}
	return fuse.ENODATA
}

// Rename a file in d
func (d *DriveDir) Rename(req *fuse.RenameRequest, node fs.Node, intr fs.Intr) fuse.Error {
	// copy the file on google drive to the new name
	_, err := service.Files.Copy(nameToFile[req.OldName].File.Id, &drive.File{Title: req.NewName}).Do()
	if err != nil {
		log.Println(err)
		return err
	}
	err = service.Files.Delete(nameToFile[req.OldName].File.Id).Do()
	if err != nil {
		log.Println(err)
		return err
	}

	go refreshAll()
	return nil
}

// FSync is a place holder and does nothing but satisfyes the FSyncer interface
func (d *DriveDir) Fsync(req *fuse.FsyncRequest, intr fs.Intr) fuse.Error {
	return nil
}

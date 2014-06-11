package drivefs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"log"
	"strings"
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
	dirChan := make(chan *[]fuse.Dirent)
	go func() {
		// List of directories to return
		var dirs []fuse.Dirent
		// get all new list of files
		f, err := service.Files.List().Do()
		if err != nil {
			log.Println(err)
		}
		fileList := f.Items
		// Populate idToFile with new ids
		for i := range fileList {
			idToFile[fileList[i].Id] = DriveFile{File: fileList[i]}
		}
		// get list of children
		// If d is at root, fetch the root children, else fetch this file's children
		if d.Root {
			c, err := service.Children.List("root").Do()
		} else {
			c, err := service.Children.List(d.Dir.Id).Do()
		}

		// Get children of this folder
		children := c.Items\

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
		dirChan <- &dirs
	}()
	// Wait for the lookups to be done, or die if interupt happens
	select {
	case tmp := <-dirChan:
		return *tmp, nil
	case <-intr:
		return nil, fuse.EINTR
	}

}

// Mkdir does nothing, because drivefs is read-only
func (d *DriveDir) Mkdir(req *fuse.MkdirRequest, intr fs.Intr) (fs.Node, fuse.Error) {
	return nil, fuse.Errno(syscall.EROFS)
}

// Mknod does nothing, because drivefs is read-only
func (d *DriveDir) Mknod(req *fuse.MknodRequest, intr fs.Intr) (fs.Node, fuse.Error) {
	return nil, fuse.Errno(syscall.EROFS)
}

// Remove does nothing, because drivefs is read-only
func (d *DriveDir) Remove(req *fuse.RemoveRequest, intr fs.Intr) fuse.Error {
	return fuse.Errno(syscall.EROFS)
}

// Removexattr does nothing, because drivefs is read-only
func (d *DriveDir) Removexattr(req *fuse.RemovexattrRequest, intr fs.Intr) fuse.Error {
	return fuse.Errno(syscall.EROFS)
}

// Rename does nothing, because drivefs is read-only
func (d *DriveDir) Rename(req *fuse.RenameRequest, node fs.Node, intr fs.Intr) fuse.Error {
	return fuse.Errno(syscall.EROFS)
}

// Setattr does nothing, because drivefs is read-only
func (d *DriveDir) Setattr(req *fuse.SetattrRequest, res *fuse.SetattrResponse, intr fs.Intr) fuse.Error {
	return fuse.Errno(syscall.EROFS)
}

// Setxattr does nothing, because drivefs is read-only
func (d *DriveDir) Setxattr(req *fuse.SetxattrRequest, intr fs.Intr) fuse.Error {
	return fuse.Errno(syscall.EROFS)
}

// Symlink does nothing, because drivefs is read-only
func (d *DriveDir) Symlink(req *fuse.SymlinkRequest, intr fs.Intr) (fs.Node, fuse.Error) {
	return nil, fuse.Errno(syscall.EROFS)
}

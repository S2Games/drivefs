package lib

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"log"
)

// Root represents the root of the filesystem
type Root struct{}

// Root is called ot get the root directory node of this filesystem
func (r Root) Root() (fs.Node, fuse.Error) {
	dir, err := service.Files.Get("root").Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &DriveDir{Dir: dir, Root: true}, nil
}

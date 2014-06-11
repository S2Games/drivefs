package drivefs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// Root represents the root of the filesystem
type Root struct{}

// Root is called ot get the root directory node of this filesystem
func (r Root) Root() (fs.Node, fuse.Error) {
	return &DriveDir{Root: true}, nil
}

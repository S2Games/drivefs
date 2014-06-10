package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"log"
)

func main() {
	mount := "/tmp/test"
	c, err := fuse.Mount(mount)
	if err != nil {
		log.Fatal(err)
	}
}

package lib

import (
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"log"
	"strings"
)

// fileIndexCache holds, and periodically retrieves the file list
type fileIndexCache map[string]*DriveFile

// refresh refreshes the id -> DriveFile pairs it retrieves from the drive api
func (i *fileIndexCache) refresh() {
	// get the file list from the google api
	f, err := service.Files.List().Do()
	if err != nil {
		log.Println(err)
		return
	}
	// wipe the indexCache
	for k, _ := range i {
		delete(i[k])
	}
	list := f.Items
	for j := range list {
		i[list[j].Id] = list[j].Id
	}
}

// dirIndexCache holds, and periodically retrives the
type dirIndexCache map[string]*drive.ChildList

// refresh refreshes the id -> pairs it retrives from the drive api
func (d *dirIndexCache) refresh() {
	// get the file list from the google api
	f, err := service.Files.List().Do()
	if err != nil {
		log.Println(err)
		return
	}
	// while the indexCache
	for k, _ := range d {
		delete(d[k])
	}
	list := f.Items
	parents := make(map[string]*drive.File)
	for i := range list {
		if strings.Contains(list[i].MimeType, "folder") {
			parents[list[i].Id] = list[i]
		}
	}
	// collect the children
	var c *drive.ChildList
	var cErr error
	for _, v := range parents {
		c, cErr = service.Children.List(v.Id).Do()
		if cErr != nil {
			log.Println(cErr)
		} else {
			d[v.Id] = c
		}
	}
	// collect the children for the root directory
	c, cErr = service.Children.List("root").Do()
	if cErr != nil {
		log.Println(err)
	} else {
		d["root"] = c
	}
}

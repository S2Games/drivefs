package lib

import (
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"log"
	"strings"
)

// refresh refreshes the id -> DriveFile pairs it retrieves from the drive api
func refreshFileIndex() {
	// create tmp map to replace fileIndex
	tmpFileIndex := make(map[string]*drive.File)
	// get the file list from the google api
	f, err := service.Files.List().Do()
	if err != nil {
		log.Println(err)
		return
	}
	list := f.Items
	for i := range list {
		tmpFileIndex[list[i].Id] = list[i]
	}
	fileIndex = tmpFileIndex
}

// refresh refreshes the id -> pairs it retrives from the drive api
func refreshDirIndex() {
	// get the file list from the google api
	f, err := service.Files.List().Do()
	if err != nil {
		log.Println(err)
		return
	}
	list := f.Items

	// make new tmp maps
	tmpChildIndex := make(map[string]*drive.ChildList)
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
			tmpChildIndex[v.Id] = c
		}
	}
	// collect the children for the root directory
	c, cErr = service.Children.List("root").Do()
	if cErr != nil {
		log.Println(err)
	} else {
		tmpChildIndex["root"] = c
	}

	// replace old index with new
	childIndex = tmpChildIndex
}

package lib

import (
	drive "code.google.com/p/google-api-go-client/drive/v2"
	"net/http"
)

var (

	// client is a google drive *&http.client
	client *http.Client

	// nameToDir maps a directory name to its DriveDir
	nameToDir map[string]*DriveDir

	// nameToFile maps a file name to its DriveFile
	nameToFile map[string]*DriveFile

	// idToDir maps a directory name to its DriveDir
	idToDir map[string]*DriveDir

	// idToFile maps a file name to its DriveFile
	idToFile map[string]*DriveFile

	// service is a google drive service
	service *drive.Service

	// childCache  is a cache of the the current directories, and their children
	// stored by directory id
	childIndex map[string]*drive.ChildList

	// fileCache is a cache of the current files stored by id
	fileIndex map[string]*drive.File
)

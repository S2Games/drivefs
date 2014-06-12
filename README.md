drivefs
=======

An experimental **FUSE** filesystem for **Google Drive** written in pure go.



## Installation

```bash
go get github.com/eliothedeman/drivefs
cd $GOPATH/src/github.com/eliothedeman/drivefs
go build
```

## Usage
```bash
./drivefs -mount /path/to/mountpoint -code yourGoogleDriveAuthCode
```

## How to get your Google Drive auth code

Simply running ``./drivefs`` will give you a link to open in a browser which will provide you with your personal code.


### Special Thanks
* [Matt Layher](https://github.com/mdlayher) For his [subfs](https://github.com/mdlayher/subfs) project, which this project was based on.
* The Creators of [Bazil](http://bazil.org) and [Russ Cox](http://swtch.com/~rsc/) for their FUSE libraries. 
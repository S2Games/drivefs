drivefs
=======

An experimental **FUSE** filesystem for **Google Drive** written in pure go.
## Status
drivefs is currently functional, but experimental. Bugs still exist so please report them!

### FUSE functions supported.
|Function|drifevs|Google Drive|
|:-:|:-:|:-:|
|Create|X|X|
|Flush|X|X|
|FSync|-|X|
|GetAttr|X|X|
|GetXAttr|-|X|
|Link|-|-|
|Mkdir|X|X|
|Read|-|X|
|ReadAll|X|X|
|ReadDir|X|X|
|Remove|X|X|
|Rename|X|X|
|RmDir|X|X|
|SetAttr|X|X|
|SetXAttr|-|-|
|Symlink|-|-|
|Update|X|X|
|Write|X|X|


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
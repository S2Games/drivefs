MD=drivefs
# Full go import path of the project
MDPath=github.com/eliothedeman/${MD}

# Build the binary for the current platform
make:
	go build -o bin/${MD}

# Remove the bin folder
clean:
	rm -rf bin/


all:
	make darwin_amd64
	make linux_amd64
	make linux_arm


darwin_amd64:
	GOOS="darwin" GOARCH="amd64" go build -o bin/${MD}_darwin_amd64

linux_amd64:
	GOOS="linux" GOARCH="amd64" go build -o bin/${MD}_linux_amd64

linux_arm:
	GOOS="linux" GOARCH="arm" go build -o bin/${MD}_linux_arm
GO15VENDOREXPERIMENT=1

build: bindata
	mkdir -p bin/$(GOOS)/$(GOARCH); cd bin/$(GOOS)/$(GOARCH); go build -v github.com/michaeldwan/static; cd ../../..

build-all:
	GOOS=darwin ARCH=amd64 make build
	GOOS=darwin ARCH=386 make build
	GOOS=linux ARCH=amd64 make build
	GOOS=linux ARCH=386 make build

install: bindata
	go install -v

bindata:
	go-bindata -pkg staticlib -prefix data/ -o ./staticlib/bindata.go data/...

test: bindata
	go test ./staticlib -v

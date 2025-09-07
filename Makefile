.PHONY: build-win

build-win:
	go test -v ./...
	cd app/cmd/hdr_switcher && goversioninfo versioninfo.json
	go build -ldflags -H=windowsgui -o hdr_switcher.exe ./app/cmd/hdr_switcher
	rm -f app/cmd/hdr_switcher/resource.syso
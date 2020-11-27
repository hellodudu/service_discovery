.PHONY: build_linux
build_linux: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/exporter_linux main.go

.PHONY: build_mac
build_mac: 
	env CGO_ENABLED=0 go build -o ../bin/exporter_mac main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../bin/exporter_win.exe main.go

.PHONY: build_all
build_all: build_linux build_mac build_win
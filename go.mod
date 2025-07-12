module github.com/javanhut/CarrionLSP

go 1.23.0

toolchain go1.24.5

require (
	github.com/javanhut/TheCarrionLanguage v0.1.7
	github.com/sourcegraph/jsonrpc2 v0.2.0
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/peterh/liner v1.2.2 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

replace github.com/javanhut/TheCarrionLanguage => ../TheCarrionLanguage

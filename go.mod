module senec-monitor

go 1.20

require github.com/lib/pq v1.10.9

require github.com/google/uuid v1.5.0 // indirect

require (
	github.com/gorilla/websocket v1.5.1 // direct
	golang.org/x/net v0.17.0 // indirect
)

require (
	github.com/fatih/color v1.16.0 // direct
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.14.0 // indirect
)

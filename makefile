build:
	go build -o soab cmd/sokker-org-auto-bidder/main.go

format:
	go fmt ./...

test:
	go test ./...

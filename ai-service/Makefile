run:
	@go run main.go

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" \
		-o ai-service \
		main.go
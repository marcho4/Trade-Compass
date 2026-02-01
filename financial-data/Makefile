run:
	@go run cmd/main.go

build:
	@go build -o financial-data cmd/main.go

build-prod:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" \
		-o financial-data \
		cmd/main.go
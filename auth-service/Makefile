run:
	@go run cmd/main.go

build:
	@go build -o auth-service cmd/main.go

build-prod:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" \
		-o auth-service \
		cmd/main.go
build:
	env GOOS=js GOARCH=wasm go build -o bin/astroapp.wasm cmd/main.go

run:
	go run cmd/main.go
build:
	env GOOS=js GOARCH=wasm go build -o html/astroapp.wasm cmd/main.go

run:
	go run cmd/main.go
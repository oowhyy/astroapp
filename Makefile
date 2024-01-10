build:
	env GOOS=js GOARCH=wasm go build -o html/astroapp.wasm cmd/astroapp/main.go

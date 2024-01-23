include .env

build:
	env GOOS=js GOARCH=wasm go build -ldflags "-X main.refreshToken=${REFRESH_TOKEN} -X main.appAuth=${APP_AUTH}" -o html/astroapp.wasm cmd/astroapp/main.go


dev:
	godotenv -f ./dev.env air -c air.toml

build-dev: $(shell find ./src -name "*.go")
	mkdir -p ./build
	go build -o ./build/mint-dev ./src

clean-dev:
	rm -rf ./tmp && rm -rf dev.sqlite && rm -rf mint-dev

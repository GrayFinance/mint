dev:
	godotenv -f ./dev.env air -c air.toml

build-dev: $(shell find ./src -name "*.go")
	mkdir -p ./build
	go build -o ./build/mint-dev ./src

test-dev: $(shell find ./src -name "*.go")
	godotenv -f ./tests/test.env go test ./tests -v

clean-dev:
	rm -rf ./tmp && rm -rf dev.sqlite
	rm -rf mint-dev && rm -rf ./tests/*.sqlite

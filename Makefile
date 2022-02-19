./bin/jwk_server: $(shell find -name "*.go")
	go build -o ./bin/jwk_server ./cmd 

build: ./bin/jwk_server
	
run: ./bin/jwk_server
	./bin/jwk_server

test:     
	go test ./...


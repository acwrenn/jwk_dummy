GPG_KEY?=dummy
REGION?=us-west-2
ARCH?=amd64
BUCKET?=dummy-apt
VERSION?=1.0.0

./bin/dummy_jwk: $(shell find -name "*.go")
	go build -o ./bin/dummy_jwk ./cmd 

build: ./bin/dummy_jwk
	
run: ./bin/dummy_jwk
	./bin/jwk_server

test:     
	go test ./...

publish: ./dummy_jwk-$(ARCH).deb
	apt-s3 -region $(REGION) -bucket $(BUCKET) -deb ./dummy_jwk-$(ARCH).deb -key $(GPG_KEY)

./dummy_jwk-$(ARCH).deb: ./bin/dummy_jwk
	mkdir -p ./out/DEBIAN
	cp files/control ./out/DEBIAN/control
	mkdir -p ./out/usr/bin
	cp ./bin/dummy_jwk ./out/usr/bin/dummy_jwk
	dpkg-deb --build ./out ./dummy_jwk-$(ARCH).deb

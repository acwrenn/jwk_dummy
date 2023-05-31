GPG_KEY?=dummy
REGION?=us-west-2
ARCH?=amd64
BUCKET?=dummy-apt
VERSION?=1.0.0

./bin/jwk_dummy: $(shell find -name "*.go")
	go build -o ./bin/jwk_dummy ./cmd 

build: ./bin/jwk_dummy
	
run: ./bin/jwk_dummy
	./bin/jwk_server

test:     
	go test ./...

publish: ./jwk_dummy-$(ARCH).deb
	apt-s3 -region $(REGION) -bucket $(BUCKET) -deb ./jwk_dummy-$(ARCH).deb -key $(GPG_KEY)

./jwk_dummy-$(ARCH).deb: ./bin/jwk_dummy
	mkdir -p ./out/DEBIAN
	cp files/control ./out/DEBIAN/control
	mkdir -p ./out/usr/bin
	cp ./bin/jwk_dummy ./out/usr/bin/jwk_dummy
	dpkg-deb --build ./out ./jwk_dummy-$(ARCH).deb

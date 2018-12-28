all:
	$(MAKE) deps
	$(MAKE) build

package_linux:
	$(MAKE) build_linux
	$(MAKE) package

clean:
	rm -rf build/
	rm -rf bin/
	rm -rf out/
	rm -rf tmp/
	rm -rf pkg/rusted
	rm -rf pkg/rusted.tgz

deps:
	dep ensure

build:
	mkdir -p bin
	go build -ldflags="-s -w" -o bin/rusted pkg/*.go

run:
	$(MAKE) build
	bin/rusted

build_linux:
	mkdir -p bin
	env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o bin/rusted pkg/*.go

package:
	mkdir -p out
	rm -rf tmp
	mkdir -p tmp
	cp bin/rusted tmp/
	tar czf out/rusted.tgz -C tmp/ rusted
	rm -rf tmp

.PHONY: clean deps package
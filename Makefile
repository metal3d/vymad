VERSION:=$(shell git describe --tags | sed 's/v//g')
OPTS:=-ldflags '-X main.VERSION=$(VERSION)'

all: clean dist linux freebsd darwin

clean:
	rm -f vymad vymd
	rm -rf dist

dist:
	mkdir dist

linux:
	GOOS=linux go build $(OPTS)
	tar cfz vymad-linux-x64.tgz vymad
	mv vymad-linux-x64.tgz dist/
	rm vymad

freebsd:
	GOOS=freebsd go build $(OPTS)
	tar cfz vymad-freebsd-x64.tgz vymad
	mv vymad-freebsd-x64.tgz dist/
	rm vymad

darwin:
	GOOS=darwin go build $(OPTS)
	tar cfz vymad-osx-x64.tgz vymad
	mv vymad-osx-x64.tgz dist/
	rm vymad

install:
	go build $(OPTS) -o vymad
	install -Dm755 vymad ~/.local/bin/vymad



VERSION:=$(shell git describe --tags | sed 's/v//g')
OPTS:=-ldflags '-X main.VERSION $(VERSION)'

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

freebsd:
	GOOS=freebsd go build $(OPTS)
	tar cfz vymad-freebsd-x64.tgz vymad
	mv vymad-freebsd-x64.tgz dist/

darwin:
	GOOS=darwin go build $(OPTS)
	tar cfz vymad-osx-x64.tgz vymad
	mv vymad-osx-x64.tgz dist/


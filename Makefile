PROGRAMS = checkcts watchcts
INSTALL = install

BUILD_DATE=$(shell date -Iseconds)
BUILD_VERSION=$(shell git rev-parse --short HEAD)

LDFLAGS=-ldflags "-X=serialtools/version.BuildVersion=$(BUILD_VERSION) -X=serialtools/version.BuildDate=$(BUILD_DATE)"

prefix=/usr
bindir=$(prefix)/bin

%: cmd/%/main.go version/version.go
	go build $(LDFLAGS) -o $@ ./$(dir $<)

all: $(PROGRAMS)

install: all
	$(INSTALL) -d -m 755 $(DESTDIR)$(bindir)
	$(INSTALL) -m 755 checkcts watchcts $(DESTDIR)$(bindir)/

clean:
	rm -f $(PROGRAMS)

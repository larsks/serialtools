PROGRAMS = checkcts watchcts
INSTALL = install
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"


prefix=/usr
bindir=$(prefix)/bin

%: cmd/%/main.go
	go build $(LDFLAGS) -o $@ ./$(dir $<)

all: $(PROGRAMS)

install: all
	$(INSTALL) -d -m 755 $(DESTDIR)$(bindir)
	$(INSTALL) -m 755 checkcts watchcts $(DESTDIR)$(bindir)/

clean:
	rm -f $(PROGRAMS)

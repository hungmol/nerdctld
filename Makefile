
GO ?= go

PREFIX ?= /usr/local
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

all: binaries

VERSION = 0.3.1

nerdctld: main.go go.mod
	$(GO) build -o $(ROOT_DIR)/bin/$@ $(BUILDFLAGS)

.PHONY: binaries
binaries: nerdctld

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fix
fix:
	golangci-lint run --fix

.PHONY: install
install: nerdctld
	install -D -m 755 $(ROOT_DIR)/bin/nerdctld $(DESTDIR)$(PREFIX)/bin/nerdctld
	install -D -m 755 $(ROOT_DIR)/systemd/nerdctl.service $(DESTDIR)$(PREFIX)/lib/systemd/system/nerdctl.service
	install -D -m 755 $(ROOT_DIR)/systemd/nerdctl.socket $(DESTDIR)$(PREFIX)/lib/systemd/system/nerdctl.socket
	install -D -m 755 $(ROOT_DIR)/systemd/10-group.conf /etc/systemd/system/nerdctl.socket.d/10-group.conf
	# install -D -m 755 $(ROOT_DIR)/systemd/nerdctl.service $(DESTDIR)$(PREFIX)/lib/systemd/user/nerdctl.service
	# install -D -m 755 $(ROOT_DIR)/systemd/nerdctl.socket $(DESTDIR)$(PREFIX)/lib/systemd/user/nerdctl.socket

.PHONY: artifacts
artifacts:
	$(RM) nerdctld
	GOOS=linux GOARCH=amd64 \
	GO111MODULE=on CGO_ENABLED=0 $(MAKE) binaries \
	BUILDFLAGS="-ldflags '-s -w'"
	GOOS=linux GOARCH=amd64 VERSION=$(VERSION) nfpm pkg --packager deb
	GOOS=linux GOARCH=amd64 VERSION=$(VERSION) nfpm pkg --packager rpm
	tar --owner=0 --group=0 -czvf nerdctld-$(VERSION)-linux-amd64.tar.gz nerdctld docker.sh
	$(RM) nerdctld
	GOOS=linux GOARCH=arm64 \
	GO111MODULE=on CGO_ENABLED=0 $(MAKE) binaries \
	BUILDFLAGS="-ldflags '-s -w'"
	GOOS=linux GOARCH=arm64 VERSION=$(VERSION) nfpm pkg --packager deb
	GOOS=linux GOARCH=arm64 VERSION=$(VERSION) nfpm pkg --packager rpm
	tar --owner=0 --group=0 -czvf nerdctld-$(VERSION)-linux-arm64.tar.gz nerdctld docker.sh
	$(RM) nerdctld

.PHONY: clean
clean:
	$(RM) $(ROOT_DIR)/bin/nerdctld
	# sudo systemctl stop nerdctl.service
	# sudo systemctl disable nerdctl.service
	# sudo systemctl stop nerdctl.socket
	# sudo systemctl disable nerdctl.socket
	# $(RM) $(DESTDIR)$(PREFIX)/lib/systemd/system/nerdctl.service
	# $(RM) $(DESTDIR)$(PREFIX)/lib/systemd/system/nerdctl.socket
	# $(RM) $(DESTDIR)$(PREFIX)/lib/systemd/user/nerdctl.service
	# $(RM) $(DESTDIR)$(PREFIX)/lib/systemd/user/nerdctl.socket
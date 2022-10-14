# Go parameters
COMPILE_VER=1.0
COMPILE_REV= $(shell git rev-parse --short HEAD)
COMPILE_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

BUILD=build
INSTALL=/usr/local/medical-sys
PWD=$(shell pwd)
GOBUILD=/usr/local/go/bin/go build
TAR_MAIN=build/medical-sys

.PHONY: build

build:
	export GOPROXY=https://goproxy.cn,direct
	$(GOBUILD) -o $(TAR_MAIN) -v -ldflags "-X 'main._BUILD_TIME_=$(COMPILE_TIME)' -X 'main._BUILD_REV_=$(COMPILE_REV)' -X 'main._BUILD_VER_=$(COMPILE_VER)'" main.go
	$(GOBUILD) -o build/medical-tools main_tools.go

clean:
	rm $(TAR_MAIN) build/medical-tools

install: $(TAR_MAIN)
	mkdir -p $(INSTALL)
	ln -s $(PWD)/application $(INSTALL)/application
	cp $(TAR_MAIN) $(INSTALL)
	cp build/medical-tools $(INSTALL)
	cp build/medical-sys.service /lib/systemd/system/
	systemctl daemon-reload
	echo systemctl enable medical-sys

upgrade: $(TAR_MAIN)
	mv $(INSTALL)/medical-sys $(INSTALL)/medical-sys.del
	cp $(TAR_MAIN) $(INSTALL)
	systemctl restart medical-sys
	rm $(INSTALL)/medical-sys.del


run:$(TAR_MAIN)
	$(TAR_MAIN)
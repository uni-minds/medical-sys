COMPILE_VER=2.1.2

GO := /usr/local/go/bin/go
GOBUILD = ${GO} build
FLAGS = "-X 'main._BUILD_TIME_=$(shell date +"%Y-%m-%d %H:%M:%S")' -X 'main._BUILD_VER_=$(COMPILE_VER)' -X 'main._BUILD_REV_=$(shell git rev-parse --short HEAD)'"

clean:
	rm -rf build/

build/medical_sys: main.go
	${GOBUILD} -o $@ -ldflags ${FLAGS} $^

build/medical_sys_tools: main_tools/main.go
	${GOBUILD} -o $@ $^

run:build/medical_sys
	$^ -v -debug

tools:build/medical_sys_tools

core:build/medical_sys

build:core tools

install: build
	mkdir -p /usr/local/uni-ledger/medical-sys
	if [ ! -d /usr/local/uni-ledger/medical-sys/application ];\
	then ln -s $(shell pwd)/application /usr/local/uni-ledger/medical-sys/;\
	fi
	cp build/* /usr/bin
	cp install/medical-sys-base/lib/systemd/system/medical-sys.service /lib/systemd/system/
	systemctl daemon-reload
	systemctl enable medical-sys

upgrade: build
	systemctl stop medical-sys
	cp build/* /usr/bin
	systemctl start medical-sys

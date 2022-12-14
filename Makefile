COMPILE_VER=3.0.1

PREFIX := env GIT_TERMINAL_PROMPT=1 GOPRIVATE=gitee.com
FLAGS := -X 'main._BUILD_TIME_=$(shell date +"%Y-%m-%d %H:%M:%S")'
FLAGS += -X 'main._BUILD_VER_=$(COMPILE_VER)'
FLAGS += -X 'main._BUILD_REV_=$(shell git rev-parse --short HEAD)'

ifneq (, $(shell which go))
	GOBUILD := ${PREFIX} go build
else
	GOBUILD := ${PREFIX} /usr/local/go/bin/go build
endif

clean:
	rm -rf build/

build/medical_sys: core_main.go
	${GOBUILD} -o $@ -ldflags "${FLAGS}" $^

build/medical_sys_tools: loader/core_tools.go
	${GOBUILD} -o $@ $^

run:build/medical_sys
	$^ -p -c ./config.yaml

debug:build/medical_sys
	$^ -v -d -p -c ./config.yaml

run_tools:build/medical_sys_tools
	$^

tools:build/medical_sys_tools

core:build/medical_sys

build:core tools

docker:core tools
	docker build -t medisys:latest -f build_docker/Dockerfile .

install: build
	mkdir -p /usr/local/uni-ledger/medical-sys
	if [ ! -d /usr/local/uni-ledger/medical-sys/application ];\
	then ln -s $(shell pwd)/application /usr/local/uni-ledger/medical-sys/;\
	fi
	cp build/* /usr/bin
	cp build_install/medical-sys-base/lib/systemd/system/medical-sys.service /lib/systemd/system/
	systemctl daemon-reload
	systemctl enable medical-sys

upgrade: build
	systemctl stop medical-sys
	cp build/* /usr/bin
	systemctl start medical-sys

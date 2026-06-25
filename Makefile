version := 0.1.11

CURDIR := $(shell pwd)
PKG_AMD64 := build/pkg-amd64
PKG_ARM64V6 := build/pkg-arm64v6
PKG_ARM64V7 := build/pkg-arm64v7
PKG_ARM64 := build/pkg-arm64
PKG_ARMHF := build/pkg-armhf
PKG_I386 := build/pkg-i386
PKG_OUT := build/packages

SYSTEMD_UNITS := packaging/systemd/glance-agent.service

.PHONY: all
all: clean glance-agent

.PHONY: glance-agent
glance-agent: clean
	go mod download
	mkdir build
	GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.x86_64
	GOOS=linux GOARCH=arm64 GOARM=6 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.aarch64-v6
	GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.aarch64-v7
	GOOS=linux GOARCH=arm64 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.aarch64
	GOOS=linux GOARCH=arm go build -ldflags="-extldflags=-static" -o ./build/glance-agent.arm
	GOOS=linux GOARCH=386 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.i386
	GOOS=windows GOARCH=amd64 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.x86_64.exe
	GOOS=windows GOARCH=arm64 GOARM=6 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.aarch64-v6.exe
	GOOS=windows GOARCH=arm64 GOARM=7 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.aarch64-v7.exe
	GOOS=windows GOARCH=arm64 go build -ldflags="-extldflags=-static" -o ./build/glance-agent.aarch64.exe

	upx --best ./build/glance-agent.x86_64
	upx --best ./build/glance-agent.aarch64-v6
	upx --best ./build/glance-agent.aarch64-v7
	upx --best ./build/glance-agent.aarch64
	upx --best ./build/glance-agent.arm
	upx --best ./build/glance-agent.i386
	upx --best ./build/glance-agent.x86_64.exe

	# Stage pkg roots per arch with binaries and systemd units
	mkdir -p $(PKG_AMD64)/usr/bin $(PKG_AMD64)/usr/lib/systemd/system $(PKG_AMD64)/usr/lib/glance-agent
	cp ./build/glance-agent.x86_64 $(PKG_AMD64)/usr/bin/glance-agent
	cp $(SYSTEMD_UNITS) $(PKG_AMD64)/usr/lib/systemd/system/
	cp ./.env.example $(PKG_AMD64)/usr/lib/glance-agent/config.env.example

	mkdir -p $(PKG_ARM64V6)/usr/bin $(PKG_ARM64V6)/usr/lib/systemd/system $(PKG_ARM64V6)/usr/lib/glance-agent
	cp ./build/glance-agent.aarch64-v6 $(PKG_ARM64V6)/usr/bin/glance-agent
	cp $(SYSTEMD_UNITS) $(PKG_ARM64V6)/usr/lib/systemd/system/
	cp ./.env.example $(PKG_ARM64V6)/usr/lib/glance-agent/config.env.example

	mkdir -p $(PKG_ARM64V7)/usr/bin $(PKG_ARM64V7)/usr/lib/systemd/system $(PKG_ARM64V7)/usr/lib/glance-agent
	cp ./build/glance-agent.aarch64-v7 $(PKG_ARM64V7)/usr/bin/glance-agent
	cp $(SYSTEMD_UNITS) $(PKG_ARM64V7)/usr/lib/systemd/system/
	cp ./.env.example $(PKG_ARM64V7)/usr/lib/glance-agent/config.env.example

	mkdir -p $(PKG_ARM64)/usr/bin $(PKG_ARM64)/usr/lib/systemd/system $(PKG_ARM64)/usr/lib/glance-agent
	cp ./build/glance-agent.aarch64 $(PKG_ARM64)/usr/bin/glance-agent
	cp $(SYSTEMD_UNITS) $(PKG_ARM64)/usr/lib/systemd/system/
	cp ./.env.example $(PKG_ARM64)/usr/lib/glance-agent/config.env.example

	mkdir -p $(PKG_ARMHF)/usr/bin $(PKG_ARMHF)/usr/lib/systemd/system $(PKG_ARMHF)/usr/lib/glance-agent
	cp ./build/glance-agent.arm $(PKG_ARMHF)/usr/bin/glance-agent
	cp $(SYSTEMD_UNITS) $(PKG_ARMHF)/usr/lib/systemd/system/
	cp ./.env.example $(PKG_ARMHF)/usr/lib/glance-agent/config.env.example

	mkdir -p $(PKG_I386)/usr/bin $(PKG_I386)/usr/lib/systemd/system $(PKG_I386)/usr/lib/glance-agent
	cp ./build/glance-agent.i386 $(PKG_I386)/usr/bin/glance-agent
	cp $(SYSTEMD_UNITS) $(PKG_I386)/usr/lib/systemd/system/
	cp ./.env.example $(PKG_I386)/usr/lib/glance-agent/config.env.example

.PHONY: deb
deb:
	# Debian packages with staged roots and post-install script
	mkdir -p $(PKG_OUT)
	fpm -s dir -t deb -v $(version) \
		--architecture amd64 \
		-C $(PKG_AMD64) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t deb -v $(version) \
		--architecture arm64v6 \
		-C $(PKG_ARM64V6) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t deb -v $(version) \
		--architecture arm64v7 \
		-C $(PKG_ARM64V7) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t deb -v $(version) \
		--architecture arm64 \
		-C $(PKG_ARM64) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t deb -v $(version) \
		--architecture armhf \
		-C $(PKG_ARMHF) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t deb -v $(version) \
		--architecture i386 \
		-C $(PKG_I386) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

.PHONY: rpm
rpm:
	# RPM packages with staged roots and post-install script
	mkdir -p $(PKG_OUT)
	fpm -s dir -t rpm -v $(version) \
		--architecture amd64 \
		-C $(PKG_AMD64) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t rpm -v $(version) \
		--architecture arm64v6 \
		-C $(PKG_ARM64V6) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t rpm -v $(version) \
		--architecture arm64v7 \
		-C $(PKG_ARM64V7) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t rpm -v $(version) \
		--architecture arm64 \
		-C $(PKG_ARM64) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t rpm -v $(version) \
		--architecture armhf \
		-C $(PKG_ARMHF) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

	fpm -s dir -t rpm -v $(version) \
		--architecture i386 \
		-C $(PKG_I386) usr/bin/glance-agent usr/lib/systemd/system/glance-agent.service usr/lib/glance-agent/config.env.example

.PHONY: clean
clean:
	if [ -d build ]; then rm -r build; fi

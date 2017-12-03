.PHONY: all galaxy-s3-gateway release-version

all: clean dist

version := $(shell cat VERSION)

#util/version.go:
release-version:
	git rev-parse HEAD|awk 'BEGIN {print "package util"} {print "const BuildGitVersion=\""$$0"\""} END{}' > util/version.go
	date +'%Y%m%d%H'| awk 'BEGIN{} {print "const BuildGitDate=\""$$0"\""} END{}' >> util/version.go

dist: galaxy-s3-gateway
	mkdir -p build/galaxy-s3-gateway-$(value version)
	cp script/run.sh build/galaxy-s3-gateway-$(value version)/bin
	cp script/galaxy-s3-gateway.cfg build/galaxy-s3-gateway-$(value version)/bin
	cd build && tar cvzf galaxy-s3-gateway.tar.gz galaxy-s3-gateway-$(value version)
	#rm -r build/galaxy-s3-gateway-$(value version)

galaxy-s3-gateway: release-version
	mkdir -p build/galaxy-s3-gateway-$(value version)/bin
	go build ${BUILD_FLAGS} -o build/galaxy-s3-gateway-$(value version)/bin/galaxy-s3-gateway

clean:
	rm -rf build


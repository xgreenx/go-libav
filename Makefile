PACKAGES=$(shell find * -name *.go -print0 | xargs -0 -n1 dirname | sort --unique)
TEST_PACKAGES=$(shell find * -name *_test.go -print0 | xargs -0 -n1 dirname | sort --unique)

FIXTURES=sample_iPod.m4v        \
         sample_iTunes.mov      \
         sample_mpeg4.mp4

.PHONY: all gofmt golint govet test clean

all: gofmt golint govet test cover

gofmt:
	@for dir in $(PACKAGES); do gofmt -s=true -d=true -l=true $${dir}; done

golint:
	@for dir in $(PACKAGES); do golint $${dir}; done

govet:
	@for dir in $(PACKAGES); do go tool vet -all $${dir}; done

FIXTURE_TARGETS=$(addprefix fixtures/,$(FIXTURES))

$(FIXTURE_TARGETS):
	mkdir -p "$(dir $@)"
	rm -f "$@.zip" "$@"
	cd "$(dir $@)" && curl -L "https://bintray.com/imkira/go-libav/download_file?file_path=$(notdir $@)" -o "$(notdir $@)"
	rm -f "$@.zip"

fixtures: $(FIXTURE_TARGETS)

cover-test:
	rm -f coverage.*
	@for dir in $(TEST_PACKAGES); do (cd $${dir} && go test -v -tags $(FFMPEG_TAG) -race -cpu=1,2,4 -coverprofile=coverage.txt -covermode=atomic || touch $(PWD)/coverage.failed); done
	@for dir in $(TEST_PACKAGES); do (if [ -f coverage.txt ]; then cat $${dir}/coverage.txt | tail -n +2 >> coverage.txt; else cp $${dir}/coverage.txt .; fi); done
	@test ! -f coverage.failed || (echo Tests failed; exit 1)

cover:
	go tool cover -html=coverage.txt -o coverage.html

clean:
	rm -rf fixtures

ffmpeg33:
	sudo apt-get install -y yasm
	sudo apt-get install -y yasm libvpx. libx264.
#	wget -O x264.tar.bz2 ftp://ftp.videolan.org/pub/x264/snapshots/x264-snapshot-20170607-2245-stable.tar.bz2
	wget -O ffmpeg34.tar.bz2 http://ffmpeg.org/releases/ffmpeg-3.4.2.tar.bz2
#	mkdir x264 && tar xf x264.tar.bz2 -C x264 --strip-components=1 && cd x264/ && ./configure --prefix=/usr/local --enable-static && make lib-static -j 4 && sudo make install
	mkdir ffmpeg34 && tar xf ffmpeg34.tar.bz2 -C ffmpeg34 --strip-components=1 && cd ffmpeg34 && ./configure --prefix=/usr/local --disable-debug --enable-pthreads --enable-nonfree --enable-gpl --disable-indev=jack --enable-libx264 --enable-shared --disable-static && make -j 4 && sudo make install
#	mkdir ffmpeg34 && tar xf ffmpeg34.tar.bz2 -C ffmpeg34 --strip-components=1 && cd ffmpeg34 && ./configure --prefix=/usr/local --disable-debug --enable-pthreads --enable-nonfree --enable-gpl --disable-indev=jack --enable-libx264 --enable-static --disable-shared --extra-libs=--static && make -j 4 && sudo make install

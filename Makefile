
analyze:
	find . -name '*.su3' -exec file {} \;

gen:
	#go run -tags generate gen.go

fmt: clean
	gofmt -w -s *.go

setup: fmt
	rsync -rav ~/i2p/certificates/ certificates

build: fmt
	go build -o reseed-monitor/reseed-monitor ./reseed-monitor

run: build
	./reseed-monitor/reseed-monitor

clean: gen
	rm -rf reseed-monitoring reseed-monitor/reseed-monitor \
		i2pseed.creativecowpat.net:8443/ \
		reseed.i2p2.no/ \
		reseed.i2pgit.org/ \
		reseed.memcpy.io/ \
		reseed.onion.im/ \
		reseed.i2p-projekt.de/ \
		i2p.novg.net/ \
		i2p.mooo.com/ \
		data-dir*/


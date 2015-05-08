

test:
	cd tmp && go test github.com/chonglou/ksana/utils


clean:
	-rm -r tmp

fmt:
	gofmt -w *

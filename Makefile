

test:
	cd tmp && go test github.com/chonglou/ksana/utils


clean:
	-rm -r tmp
	-rm -r utils/config

fmt:
	gofmt -w gails gake

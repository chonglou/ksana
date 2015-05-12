
test:
	go test 


install:
	go install



clean:
	-rm -r tmp

fmt:
	gofmt -w *.go utils sessions




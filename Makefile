
test:
	go test


install:
	go install



clean:	
	-rm tmp/migrate/*.json




fmt:
	gofmt -w *.go examples/*.go auth/*.go

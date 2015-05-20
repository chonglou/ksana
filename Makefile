
test:
	go test


install:
	go install



clean:	
	-rm -r /tmp/migrate/*.sql




fmt:
	gofmt -w *.go orm/*.go examples/*.go #auth/*.go

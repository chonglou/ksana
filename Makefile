
test:
	cd utils && go test
	cd orm && go test
	cd web && go test
	cd redis && go test
	#go test


install:
	go install



clean:	
	-rm -r /tmp/migrate/*.sql




fmt:
	gofmt -w *.go orm/*.go web/*.go utils/*.go redis/*.go examples/*.go #auth/*.go


test:
	cd utils && go test
	cd orm && go test
	cd web && go test
	cd redis && go test
	cd i18n && go test
	go test


install:
	go install



clean:	
	-rm -r /tmp/migrate/*.sql




fmt:
	gofmt -w *.go examples/*.go

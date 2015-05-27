
test:
	go test


install:
	go install



clean:
	-rm tmp/migrate/*.json


vet:
	go vet *.go
	go vet examples/*.go
	go vet auth/*.go


fmt:
	go fmt *.go
	go fmt examples/*.go
	go fmt auth/*.go
	go fmt platform/*.go

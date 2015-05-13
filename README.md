Ksana(A golang web framework)
--------------------------------

## Install

    bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer) # If you are using zsh just change bash with zsh
    # Restart your terminal session
    gvm listall # List all Go versions available for download
    gvm install go1.4.2 # Install go
    gvm list
    gvm use go1.4.2 --default # Using go1.4.2
    go get github.com/chonglou/ksana

## Examples

 * [app.go](examples/app.go)
 * [context.xml](examples/context.xml)

    go run app.go -h
    go run app.go -migrate # need create your database first
    go run app.go -server





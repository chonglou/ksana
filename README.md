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

### Files

 * [app.go](examples/app.go)
 * [config.json](examples/config.json)

### Usage
    cd examples/
    go run app.go -h # show options
    go run app.go -r migrate  # database migrate
    go run app.go -r rollback # database rollback
    go run app.go -r server   # run server
    go run app.go -r routes   # show http routes
    go run app.go -r db       # connect to database



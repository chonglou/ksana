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

## Getting Started
After installing Go and setting up your GOPATH, create your first .go file. We'll call it app.go.

    package main
    import "github.com/chonglou/ksana"
    func main() {
      app = ksana.Application{}
      app.Get("/", func()string{
        return "Hello, Ksona!"
      })
      app.Run(8080)
    }


And then run your server:

    go run app.go

## Some 3rd lib
    go get github.com/lib/pq
    go get github.com/fzzy/radix/redis

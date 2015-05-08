Ksana(A golang web framework)
--------------------------------

## Install
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

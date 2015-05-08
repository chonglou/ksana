package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//"github.com/chonglou/ksana/utils"
)

const version = "v20150507"

func do_usage() {
	fmt.Println("Usage: ")
	fmt.Println("  gails new APP_NAME\t\t\t\t# create a new app")
	fmt.Println("  gails -e ENVIRONMENT -p PORT server\t\t# start web server")
	fmt.Println("  gails -e ENVIRONMENT db\t\t\t# start web server")
	os.Exit(0)
}

func do_server() {
	env := flag.String("e", "development", "Specifies the environment to run this server under (test/development/production).")
	port := flag.Int("p", 8080, "Runs Gails on the specified port.")
	flag.Parse()
	log.Printf("=> Booting Ksana %s", version)
	log.Printf("=> Application starting in %s on http://0.0.0.0:%v\n", *env, *port)
	log.Println("=> Run `gails -h server` for more startup options")
	log.Println("=> Ctrl-C to shutdown server")
}

func do_db() {
	env := flag.String("e", "development", "Specifies the environment to run this server under (test/development/production).")
	flag.Parse()
	log.Printf("=> Booting Ksana %s", version)
	log.Printf("=> Application starting in %s\n", *env)
	log.Println("=> Run `gails -h db` for more startup options")
}

func do_new(name string) {
	log.Printf("Using Ksana %s", version)
	if _, err := os.Stat(name); err == nil {
		log.Fatalf("File [%s] exists!!!", name)
	}

	dirs := []string{
		"app/jobs",
		"app/mailers",
		"app/models",
		"app/controllers",
		"config/environments",
		"config/initializers",
		"db/migrate",
		"log",
		"lib",
		"tmp/pids",
		"tmp/sockets",
		"public/assets",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(name+"/"+d, 0755); err != nil {
			log.Fatal(err)
		}
	}

	files := make(map[string]string)
	files[".gitignore"] = ""
	files["config/settings.yml"] = ""
	files["config/environment.go"] = ""
	files["config/environments/test.go"] = ""
	files["config/environments/development.go"] = ""
	files["config/environments/production.go"] = ""
	files["db/seed.go"] = ""
	files["db/migrate/20150508055759_init.go"] = ""
	files["public/favicon.ico"] = ""
	files["public/robots.txt"] = ""
	files["public/404.html"] = ""
	files["public/422.html"] = ""
	files["public/500.html"] = ""

	for k, v := range files {
		if err := ioutil.WriteFile(name+"/"+k, []byte(v), 0644); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Create application %s success!", name)
}

func main() {
	size := len(os.Args)

	switch {
	case size == 1:
		do_usage()
	case os.Args[size-1] == "server":
		do_server()
	case os.Args[size-1] == "db":
		do_db()
	case os.Args[size-2] == "new":
		do_new(os.Args[size-1])
	default:
		do_usage()
	}
}

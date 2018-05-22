package main

//go:generate rice embed-go

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mmcloughlin/professor"

	"git.mills.io/prologic/je"
)

func main() {
	var (
		version bool
		debug   bool

		dbpath  string
		bind    string
		workers int
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug logging")

	flag.StringVar(&dbpath, "dbpath", "je.db", "Database path")
	flag.StringVar(&bind, "bind", "0.0.0.0:8000", "[int]:<port> to bind to")
	flag.IntVar(&workers, "workers", 16, "worker pool size")

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if version {
		fmt.Printf("je v%s", je.FullVersion())
		os.Exit(0)
	}

	if debug {
		go professor.Launch(":6060")
	}

	opts := &je.Options{
		Workers: workers,
	}

	db := je.InitDB(dbpath)
	defer db.Close()

	log.Infof("je %s listening on %s", je.FullVersion(), bind)
	je.NewServer(bind, opts).ListenAndServe()
}
package main

import (
	"log"
	"github.com/eyedeekay/reseed-monitoring"
	"io/ioutil"
)

func main() {
	config, err := monitor.SortedMap("config.json")
	if err != nil {
		log.Fatal(err)
	}
	errs := monitor.SortedMonitor(config)
	if errs != nil {
		if len(errs) > 0 {
			log.Println(errs)
		}
	}
	index, err := monitor.GeneratePage()
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("index.html", []byte(index), 0644)
}

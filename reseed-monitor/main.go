package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/eyedeekay/reseed-monitoring"
)

func main() {
	go loop()
	if err := http.ListenAndServe("0.0.0.0:7672", &monitor.MonitorServer{}); err != nil {
		log.Fatal(err)
	}
}

func loop() {
	for {
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
		err = ioutil.WriteFile("index.html", []byte(index), 0644)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(24 * time.Hour))
	}
}

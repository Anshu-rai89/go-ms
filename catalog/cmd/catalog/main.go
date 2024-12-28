package main

import (
	"log"
	"time"

	"github.com/Anshu-rai89/go-ms/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	ElasticSearchURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("catalog %v", config)
	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(config.ElasticSearchURL)

		if err != nil {
			log.Println(err)
		}
		return
	})

	defer r.Close()
	log.Println("Listening on port:8080")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}

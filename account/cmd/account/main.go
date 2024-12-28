package main

import (
	"log"
	"time"

	"github.com/Anshu-rai89/go-ms/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("catalog %v", config)
	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(config.DatabaseURL)

		if err != nil {
			log.Println(err)
		}
		return
	})

	defer r.Close()
	log.Println("Listening on port:8080")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}

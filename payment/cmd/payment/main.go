package main

import (
	"log"
	"time"

	"github.com/Anshu-rai89/go-ms/payment"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type AppConfig struct {
	Url string `envconfig:"DATABASE_URL"`
}

func main() {
	var config AppConfig

	err := envconfig.Process("", &config)

	if err != nil {
		log.Fatal(err)
	}

	var r payment.Repository
	retry.ForeverSleep(time.Second*2, func(_ int) (err error) {
		r, err = payment.NewPostgresRepository(config.Url)

		if err != nil {
			log.Println(err)
			return err
		}

		return
	})

	defer r.Close()

	s := payment.NewPaymentService(r)
	log.Fatal(payment.ListenGRPC(s, 8080))
}

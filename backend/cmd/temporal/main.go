package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"laxo.vn/laxo/processing"
)

func main() {
  log.Println("Starting workers...")

  c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "product", worker.Options{})

	w.RegisterWorkflow(processing.ProcessLazadaProducts)

	activities := &processing.Activities{RedisClient: "Cool Redis Client"}
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

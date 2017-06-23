package main

import (
	"flag"
	"fmt"
	"github.com/itmarketplace/go-queue"
	"github.com/nsqio/go-nsq"
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

var start = time.Now()
var ops uint64 = 0
var numbPtr = flag.Int("msg", 5, "number of messages (default: 10000)")
var ipnsqlookupd = flag.String("ipnsqlookupd", "", "IP address of ipnsqlookupd")

func main() {

	flag.Parse()

	c := queue.NewConsumer("India", "ch")

	c.Set("nsqlookupd", "127.0.0.1:4161")
	c.Set("concurrency", runtime.GOMAXPROCS(runtime.NumCPU()))
	c.Set("max_attempts", 10)
	c.Set("max_in_flight", 150)
	c.Set("default_requeue_delay", "15s")
	c.Set("max_backoff_duration", "1s")
	c.Set("max_attempts", 2)

	c.Start(nsq.HandlerFunc(func(msg *nsq.Message) error {

		log.Println(string(msg.Body))

		atomic.AddUint64(&ops, 1)
		if ops == uint64(*numbPtr) {
			elapsed := time.Since(start)
			log.Printf("Time took %s", elapsed)
		}

		/*fmt.Println("Reque/backoff test")
		msg.Requeue(time.Second * 10)*/

		return nil
	}))

	fmt.Scanln()

	c.Stop()
}

//cerror.go

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/itmarketplace/go-queue"
	"github.com/nsqio/go-nsq"
	"runtime"
	"strconv"
)

var numbPtr = flag.Int("msg", 100, "number of messages (default: 100)")
var lkp = flag.String("lkp", "", "IP address of nsqlookupd")

func FuncSimulateError(msg *nsq.Message) error {

	str := string(msg.Body)

	//http://play.golang.org/p/MFboCiikYW
	t := 0 //number of vowels in str
	for _, value := range str {
		switch value {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
			t++
		}
	}

	if t%2 == 0 {
		return errors.New(strconv.Itoa(t))
	} else {
		return nil
	}
}

func main() {

	flag.Parse()

	c := queue.NewConsumer("India", "ch")

	c.Set("nsqlookupd", *lkp+":4161")
	c.Set("concurrency", runtime.GOMAXPROCS(runtime.NumCPU()))
	c.Set("max_attempts", 10)
	c.Set("max_in_flight", 150)
	c.Set("default_requeue_delay", "15s")

	c.Start(nsq.HandlerFunc(FuncSimulateError))

	fmt.Scanln()

	c.Stop()
}


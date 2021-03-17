package redis_mutex

import (
	"log"
	"testing"
	"time"
)

func TestMu(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	mutex, err := New(SetUri("127.0.0.1", 6379))
	if err != nil {
		log.Fatalln(err)
	}

	look, err := mutex.Look()
	if err != nil {
		log.Fatalln(err)
	}

	if !look {
		log.Fatalln("Lock Failure look")
	}

	time.Sleep(time.Second * 10)

	unLook, err := mutex.UnLook()
	if err != nil {
		log.Fatalln(err)
	}

	if !unLook {
		log.Fatalln("Lock Failure unLook")
	}

	time.Sleep(time.Second * 10)
}

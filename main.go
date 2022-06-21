package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"


	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	seedBrokers = flag.String("brokers", "127.0.0.1:54789", "comma delimited list of seed brokers")
	topic       = flag.String("topic", "cwxstat", "topic to consume from")
	style       = flag.String("commit-style", "autocommit", "commit style (which consume & commit is chosen); autocommit|records|uncommitted")
	group       = flag.String("group", "group", "group to consume within")
	logger      = flag.Bool("logger", false, "if true, enable an info level logger")
)

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	seeds := []string{"127.0.0.1:54455", "127.0.0.1:54461", "127.0.0.1:54460"}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		//kgo.ConsumerGroup("group"),
		kgo.ConsumeTopics("cwxstat"),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()

	// 1.) Producing a message
	// All record production goes through Produce, and the callback can be used
	// to allow for synchronous or asynchronous production.
	var wg sync.WaitGroup
	wg.Add(1)
	record := &kgo.Record{Topic: "cwxstat", Value: []byte("bar..wow worked")}
	cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			fmt.Printf("record had a produce error: %v\n", err)
		}

	})
	wg.Wait()

	// Alternatively, ProduceSync exists to synchronously produce a batch of records.
	if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
		fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
	}

	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		// We can iterate through a record iterator...
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Println(string(record.Value), "from an iterator!")
		}

		// or a callback function.
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				fmt.Println(string(record.Value), "from range inside a callback!")
			}

			// We can even use a second callback!
			p.EachRecord(func(record *kgo.Record) {
				fmt.Println(string(record.Value), "from a second callback!")
			})
			fmt.Println("topic:", p.Topic)
		})
	}

}

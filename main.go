package main

import (
	"fmt"
	"os"
	"time"
	"flag"
	kinesis "github.com/sendgridlabs/go-kinesis"
)

var streamName *string
var region *string

func init() {
	streamName = flag.String("stream", "", "stream name")
	region = flag.String("region", "", "AWS region eg. us-west-2")
	flag.Parse()
}

func main() {
	auth, err := kinesis.NewAuthFromEnv()

	if err != nil {
		fmt.Printf("Unable to retrieve authentication credentials from the environment: %v", err)
		os.Exit(1)
	}

	awsRegion := *region
	stream := *streamName

	fmt.Println("authenticating")
	ksis := kinesis.New(auth, awsRegion)

	fmt.Println("finding shards")
	args := kinesis.NewArgs()
	args.Add("StreamName", stream)
	resp3, _ := ksis.DescribeStream(args)

	exampleShardID := resp3.StreamDescription.Shards[0].ShardId

	fmt.Println("connecting to shard", exampleShardID)

	args = kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", exampleShardID)
	args.Add("ShardIteratorType", "TRIM_HORIZON") // AT_SEQUENCE_NUMBER | AFTER_SEQUENCE_NUMBER | TRIM_HORIZON | LATEST
	resp, _ := ksis.GetShardIterator(args)

	shardIterator := resp.ShardIterator

	receptions := 0
	receivedActivity := 0

	for {
		args = kinesis.NewArgs()
		args.Add("ShardIterator", shardIterator)

		receptions += 1

		records, err := ksis.GetRecords(args)

		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, d := range records.Records {
			fmt.Printf("GetRecords Data: %v\n", string(d.GetData()))
			receivedActivity += 1
		}

		fmt.Println("received", receivedActivity, "records in ", receptions, "requests")

		shardIterator = records.NextShardIterator

		time.Sleep(5000 * time.Millisecond)
	}
}

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

func getShardRecords(client kinesis.KinesisClient, streamName string, shardId string, messages chan []kinesis.GetRecordsRecords) {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", shardId)
	args.Add("ShardIteratorType", "LATEST") // AT_SEQUENCE_NUMBER | AFTER_SEQUENCE_NUMBER | TRIM_HORIZON | LATEST
	resp, _ := client.GetShardIterator(args)

	shardIterator := resp.ShardIterator

	for {
		args = kinesis.NewArgs()
		args.Add("ShardIterator", shardIterator)
		records, err := client.GetRecords(args)

		if err != nil {
			time.Sleep(3000 * time.Millisecond)
			continue
		}

		messages <- records.Records

		shardIterator = records.NextShardIterator
		time.Sleep(2000 * time.Millisecond)
	}
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
	description, _ := ksis.DescribeStream(args)

	messages := make(chan []kinesis.GetRecordsRecords)

	for _, shard := range description.StreamDescription.Shards {
		go getShardRecords(ksis, stream, shard.ShardId, messages)
	}

	for {
		for _, d := range <-messages {
			fmt.Println(string((d.GetData())))
		}
	}
}

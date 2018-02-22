package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func simpleTest(workerIndex int, totalMessages int, wg *sync.WaitGroup) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	// manually rotate partition
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.CompressionLevel = sarama.CompressionLevelDefault
	config.Producer.Compression = sarama.CompressionNone

	brokers := []string{"kafka:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	// kkTopic1 = 1 partition
	// kkTopic2 = 10 partitions

	topic := "kkTopic2"

	i := 0
	bytes := 0

	messageText := fmt.Sprintf("%d%d%d%d%d", (i/10000)%10, (i/1000)%10, (i/100)%10, (i/10)%10, i%10)
	messageText = "ABCDE"
	fmt.Printf("\nWorker %d going to send %d messages, each %d bytes.\n", workerIndex, totalMessages, len(messageText))
	startTime := time.Now()
	partition := int32(0)
	for i < totalMessages {

		messageText = fmt.Sprintf("%d%d%d%d%d", (i/10000)%10, (i/1000)%10, (i/100)%10, (i/10)%10, i%10)
		bytes = bytes + len(messageText)
		partition = int32(i % 10)
		//fmt.Printf("partition %d \n", partition)
		msg := &sarama.ProducerMessage{
			Topic:     topic,
			Value:     sarama.StringEncoder(messageText),
			Partition: partition,
		}
		if i%1000 == 0 {
			fmt.Printf("%d-%d ... ", workerIndex, i)
		}
		//partition, offset, err := producer.SendMessage(msg)
		producer.SendMessage(msg)

		if err != nil {
			panic(err)
		}
		//fmt.Printf("Message (%d) is stored in topic(%s)/partition(%d)/offset(%d)\n", i, topic, partition, offset)
		i++
		//time.Sleep(5 * time.Millisecond)
	}

	t := time.Now()
	elapsed := t.Sub(startTime)
	bytesPerSecond := float64(bytes) / elapsed.Seconds()
	messagesPerSecond := float64(totalMessages) / elapsed.Seconds()
	fmt.Printf("\nWorker %d sent %d messages totaling %d bytes in %s. Bytes/second: %f. Messages/second : %f.\n", workerIndex, i, bytes, elapsed, bytesPerSecond, messagesPerSecond)

	wg.Done()
}

func main() {

	targetMessages := 1000000 // total amount of messages we produce
	bytesPerMessage := 5      // ToDo : now this is just used for calc and msg has always 5 bytes.
	workers := 200            // how many worker threads (goroutines) are created to divide the work.

	messagesPerWorker := targetMessages / workers
	totalBytes := workers * messagesPerWorker * bytesPerMessage
	totalMessages := workers * messagesPerWorker
	i := 1

	var wg sync.WaitGroup
	wg.Add(workers)

	fmt.Printf("Start %d workers each to send %d messages.\n", workers, messagesPerWorker)
	mainStartTime := time.Now()

	for i <= workers {
		go simpleTest(i, messagesPerWorker, &wg)
		i++
	}

	fmt.Println("wait for workers to finish...")
	wg.Wait()

	fmt.Println("\nSUMMARY\n")
	elapsedTotal := time.Now().Sub(mainStartTime)
	fmt.Printf("All workers finished in %s seconds.\n", elapsedTotal)
	fmt.Printf("%d workers sent in total %d messages.\n", workers, totalMessages)
	messagesPerSecond := float64(totalMessages) / elapsedTotal.Seconds()
	fmt.Printf("Workers sent %f messages per second.\n", messagesPerSecond)
	bytesPerSecond := float64(totalBytes) / elapsedTotal.Seconds()
	fmt.Printf("Workers sent %f bytes per second.\n", bytesPerSecond)
}

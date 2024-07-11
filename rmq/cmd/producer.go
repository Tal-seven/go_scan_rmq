package main

import (
	"context"
	"fmt"
	"go_scan_rmq/rmq/internal"
	"go_scan_rmq/services/common/genproto/sensors"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

var i float32 = 0

func sendScanMessages(ctx context.Context, done <-chan int, stream <-chan int, client internal.RabbitClient) {
	for {
		select {
		case <-done:
			return
		case <-stream:
			scanInstance := sensors.Scan{
				AngleMin:       0.15,
				AngleMax:       1.74,
				Ranges:         make([]float32, 540, 540),
				AngleIncrement: 0.01,
				ScanTime:       i,
			}
			i += 1
			wireMesg, err := proto.Marshal(&scanInstance)

			if err != nil {
				fmt.Println("Failed to serialize scan message")
			}

			if err := client.Send(ctx, "robot1", "sensors.lidar.scan", amqp091.Publishing{
				ContentType:  "text/plain",
				DeliveryMode: amqp091.Persistent,
				Body:         []byte(wireMesg),
			}); err != nil {
				panic(err)
			}
			fmt.Printf("Message Scan Sent %v\n", time.Now())
		}
	}
}

func main() {
	conn, err := internal.ConnectRabbitMQ("seven-admin", "123456", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	if err := client.NewExchangeDeclare("robot1", "topic", true, false); err != nil {
		panic(err)
	}

	if err := client.NewQueueDeclare("lidarscan", true, false); err != nil {
		panic(err)
	}
	// if err := client.NewQueueDeclare("customers_test", false, false); err != nil {
	// 	panic(err)
	// }

	if err := client.CreateBinding("lidarscan", "sensors.lidar.*", "robot1"); err != nil {
		panic(err)
	}

	// if err := client.CreateBinding("customers_test", "customers.*", "customer_events"); err != nil {
	// 	panic(err)
	// }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	done := make(chan int)
	defer close(done)

	stream := make(chan int)

	go func() {
		for {
			select {
			case <-done:
				return
			case stream <- 1:
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	go sendScanMessages(ctx, done, stream, client)

	time.Sleep(time.Hour * 10)
	log.Println(client)
}

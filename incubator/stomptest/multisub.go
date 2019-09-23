package testmulti

import (
	"fmt"
	"sync"

	"github.com/go-stomp/stomp"
)

var channels []<-chan *stomp.Message

func handle(channel <-chan *stomp.Message) {
	for {
		m := <-channel
		fmt.Printf("Destination %v, Message: %v\n", m.Destination, string(m.Body))
	}
}

func combine(inputs []<-chan *stomp.Message, output chan<- *stomp.Message) {
	var group sync.WaitGroup
	for i := range inputs {
		group.Add(1)
		go func(input <-chan *stomp.Message) {
			for val := range input {
				output <- val
			}
			group.Done()
		}(inputs[i])
	}
	go func() {
		group.Wait()
		close(output)
	}()
}

func main() {

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.HeartBeat(0, 0),
	}

	fmt.Printf("Connecting to Stomp server\n")
	conn, err := stomp.Dial("tcp", "3dexp.18xfd05.ds:61613", options...)
	if err != nil {
		fmt.Printf("Error Connecting to Stomp server: %v\n", err)
		return
	}
	sub, err := conn.Subscribe("/queue/testQueue", stomp.AckClient, stomp.SubscribeOpt.Header("id", "Sub1"))
	if err != nil {
		fmt.Printf("Error Subscribing to queue: %v\n", err)
		return
	}
	channels = append(channels, sub.C)
	sub, err = conn.Subscribe("/topic/testTopic", stomp.AckClient, stomp.SubscribeOpt.Header("id", "Sub2"))
	if err != nil {
		fmt.Printf("Error Subscribing to topic: %v\n", err)
		return
	}
	channels = append(channels, sub.C)

	outchan := make(chan *stomp.Message)
	combine(channels, outchan)

	go handle(outchan)

	for {
	}
}

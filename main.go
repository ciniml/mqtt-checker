// MQTT connection checker
// Author: Kenta IDA (fuga@fugafuga.org)
// License: Boost Software License 1.0

package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const testServer = "test.mosquitto.org:1883"
const timeoutSeconds = 10

func innerMain() error {
	messageCh := make(chan mqtt.Message)
	var handler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
		messageCh <- message
	}

	options := mqtt.NewClientOptions()
	options.AddBroker("tcp://" + testServer)
	client := mqtt.NewClient(options)

	fmt.Printf("Connecting to %s... \n", testServer)
	if token := client.Connect(); !token.WaitTimeout(time.Second * timeoutSeconds) {
		return fmt.Errorf("Connection failed: Timed out")
	} else if err := token.Error(); err != nil {
		return fmt.Errorf("Connection failed: %s", err)
	}

	subscription := client.Subscribe("test", 0, handler)
	if !subscription.WaitTimeout(time.Second * timeoutSeconds) {
		return fmt.Errorf("Subscription failed: Timed out")
	} else if err := subscription.Error(); err != nil {
		return fmt.Errorf("Subscription failed: %s", err)

	}

	fmt.Println("Publishing some data... ")
	if token := client.Publish("test", 0, false, "test"); !token.WaitTimeout(time.Second * timeoutSeconds) {
		return fmt.Errorf("Publish failed: Timed out")
	} else if err := token.Error(); err != nil {
		return fmt.Errorf("Publish failed: %s", err)
	}

	fmt.Println("Waiting for receiving some data... ")
	timer := time.NewTimer(time.Second * timeoutSeconds)
	select {
	case message := <-messageCh:
		fmt.Printf("Subscribed data received. topic: %s\n", message.Topic())
	case <-timer.C:
		return fmt.Errorf("Subscription data timed out")
	}

	fmt.Println("MQTT connection is OK.")
	return nil
}

func main() {
	if err := innerMain(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Press Enter to exit")
	fmt.Scanln()
}

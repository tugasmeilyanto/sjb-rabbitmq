package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type User struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Age         int    `json:"age"`
	PhoneNumber string `json:"phone_number"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msg {
			var user User
			if err := json.Unmarshal(d.Body, &user); err != nil {
				log.Printf("error decoding JSON: %s", err)
				continue
			}
			SendEmail(user)
		}
	}()

	fmt.Println("Que consumed")

	log.Printf("Waiting for event.")
	<-forever
}

func SendEmail(user User) {
	log.Printf("Email sent to %s", user.Email)
	//smtp
}

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type User struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Age         int    `json:"age"`
	PhoneNumber string `json:"phone_number"`
}

func main() {
	// REST API
	// HTTP METHOD : POST, PUT, GET, DELETE
	// ROUTER / ENDPOINT : /users , /user/1
	// Web Server : host:port -> localhost:8080
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
	e := echo.New()
	e.POST("/register", func(c echo.Context) error {
		user := new(User)
		err := c.Bind(user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
		}
		// logic submit data to database
		if err := sendEmailNotif(ch, q.Name, user); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email notification"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User registered successfully"})
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func sendEmailNotif(ch *amqp.Channel, queName string, user *User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		queName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	log.Printf("sent email notification for user: %s", user.Username)
	return nil
}

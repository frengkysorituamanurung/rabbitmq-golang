package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Koneksi ke RabbitMQ
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	failOnError(err, "Gagal terhubung ke RabbitMQ")
	defer conn.Close()

	// Membuka channel
	ch, err := conn.Channel()
	failOnError(err, "Gagal membuka channel")
	defer ch.Close()

	// Mendeklarasikan queue (harus sama dengan publisher)
	q, err := ch.QueueDeclare(
		"test_queue", // Nama queue
		false,        // Durable
		false,        // Auto-delete
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	failOnError(err, "Gagal mendeklarasikan queue")

	// Menerima pesan dari queue
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-acknowledge (jika false, harus manual ACK)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	failOnError(err, "Gagal mendaftarkan consumer")

	fmt.Println(" [*] Menunggu pesan. Tekan CTRL+C untuk keluar.")

	// Loop membaca pesan
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf(" [x] Diterima: %s\n", d.Body)
		}
	}()
	<-forever // Keep running
}

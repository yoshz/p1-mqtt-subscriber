package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	databaseUrl  string = getEnv("DATABASE_URL", "postgres://metrics:metrics@localhost/metrics?sslmode=disable")
	mqttBroker   string = getEnv("MQTT_BROKER", "tcp://localhost:1883")
	mqttUsername string = getEnv("MQTT_USERNAME", "")
	mqttPassword string = getEnv("MQTT_PASSWORD", "")
	mqttTopic    string = getEnv("MQTT_TOPIC", "energy/meters")
	message      EnergyMeterMessage
	db           *sql.DB
)

type EnergyMeterMessage struct {
	Time        time.Time `json:"time"`
	Location    string    `json:"location"`
	PowerDraw   int64     `json:"powerDraw"`
	PowerMeter1 int64     `json:"powerMeter1"`
	PowerMeter2 int64     `json:"powerMeter2"`
	GasMeter    int64     `json:"gasMeter"`
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var err error
	var message EnergyMeterMessage
	err = json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		"INSERT INTO energy_meter (time, location, power_meter1, power_meter2, gas_meter) VALUES ($1, $2, $3, $4, $5)",
		message.Time,
		message.Location,
		message.PowerMeter1,
		message.PowerMeter2,
		message.GasMeter,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to postgres")

	opts := MQTT.NewClientOptions()
	opts.SetClientID("p1-mqtt-subscriber")
	opts.AddBroker(mqttBroker)
	if mqttUsername != "" {
		opts.SetUsername(mqttUsername)
	}
	if mqttPassword != "" {
		opts.SetPassword(mqttPassword)
	}
	opts.OnConnect = func(client MQTT.Client) {
		if token := client.Subscribe(mqttTopic, 0, messagePubHandler); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
		log.Printf("Subscribed to topic %s", mqttTopic)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

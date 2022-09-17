package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var err error
	var message EnergyMeterMessage
	err = json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var err error

	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
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
	opts.CleanSession = false
	opts.OnConnect = func(client MQTT.Client) {
		if token := client.Subscribe(mqttTopic, 1, messagePubHandler); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
		log.Printf("Subscribed to topic %s", mqttTopic)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
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

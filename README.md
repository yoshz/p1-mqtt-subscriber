# P1-MQTT-Subscriber

A small Go application that reads P1 messages from MQTT and saves them to TimescaleDB hypertable.

## Local usage

Start TimescaleDB and Mosquitto
```bash
docker compose up -d
```

Create database user
```bash
cat <<EOF | docker compose exec timescaledb psql -U postgres
CREATE USER metrics WITH PASSWORD 'metrics';
CREATE DATABASE metrics OWNER 'metrics';
EOF
```

Create hypertable
```bash
cat sql/hypertable.sql | docker compose exec timescaledb psql -U postgres metrics
```

Install Go modules
```bash
go get
```

Run the P1-MQTT-Publisher
```bash
go run main.go
```

Publish example message
```bash
cat example/message.json | docker compose exec mosquitto mosquitto_pub -h mosquitto -t energy/meters -l
```

## Environment variables

 * `DATABASE_URL`: url of the timescaledb database (default is "postgres://metrics:metrics@localhost/metrics?sslmode=disable")
 * `MQTT_BROKER`: url of the MQTT broker (default is "tcp://localhost:1883")
 * `MQTT_USERNAME`: Username to authenticate to MQTT broker
 * `MQTT_PASSWORD`: Password to authenticate to MQTT broker
 * `MQTT_TOPIC`: MQTT Topic to publish to (default is "energy/meters")

package main

import "os"

func env() {
	var tmp string

	tmp = os.Getenv("PORT")
	if tmp != "" {
		Port = tmp
	}
	tmp = ""

	tmp = os.Getenv("DEBUG")
	if tmp != "" {
		Debug = tmp
	}
	tmp = ""

	tmp = os.Getenv("CONFIG_PATH")
	if tmp != "" {
		ConfigPath = tmp
	}
	tmp = ""

	tmp = os.Getenv("POSTGRES_CONN")
	if tmp != "" {
		PostgresConn = tmp
	}
	tmp = ""

	tmp = os.Getenv("KAFKA_BROKERS")
	if tmp != "" {
		KafkaBrokers = tmp
	}
	tmp = ""

	tmp = os.Getenv("KAFKA_TOPIC")
	if tmp != "" {
		KafkaTopic = tmp
	}
	tmp = ""

	tmp = os.Getenv("KAFKA_GROUP_ID")
	if tmp != "" {
		KafkaGroupID = tmp
	}
	tmp = ""

	tmp = os.Getenv("TEMPLATES")
	if tmp != "" {
		Templates = tmp
	}
	tmp = ""
}

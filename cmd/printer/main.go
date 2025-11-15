package main

import (
	"encoding/json"
	"log"
	"os"
	"other/go-printer-termal/application"
	"other/go-printer-termal/domain"
	cfg "other/go-printer-termal/infrastructure/config"
	escpos "other/go-printer-termal/infrastructure/escpos"
	mqttclient "other/go-printer-termal/infrastructure/mqttclient"
	"strings"

	"github.com/alexbrainman/printer"
)

func main() {
	_ = cfg.LoadDotEnv(".env")
	settings, err := cfg.LoadSettings()
	if err != nil {
		log.Fatalf("load settings: %v", err)
	}

	sub, err := mqttclient.NewSubscriberFromEnv()
	if err != nil {
		log.Fatalf("mqtt connect: %v", err)
	}
	defer sub.Disconnect()

	topic := settings.TopicPrefix + settings.IDMerchant
	uc := application.PrintInvoiceUseCase{PaperWidth: settings.PaperWidth}

    err = sub.Subscribe(topic, func(msg []byte) {
        log.Printf("MQTT received | topic=%s | bytes=%d", topic, len(msg))
        log.Printf("MQTT payload: %s", string(msg))
        var payload domain.Payload
        if err := json.Unmarshal(msg, &payload); err != nil {
            log.Println("parse payload:", err)
            return
        }
        log.Printf("Parsed payload | order=%s | merchant=%s", payload.NoOrder, payload.MerchantName)
        p, err := printer.Open(settings.PrinterName)
        if err != nil {
            log.Println("open printer:", err)
            return
        }
        defer p.Close()
        log.Printf("Printer opened | name=%s", settings.PrinterName)
        if err = p.StartRawDocument("Receipt"); err != nil {
            log.Println("StartRawDocument:", err)
            return
        }
        log.Println("StartRawDocument OK")
        ep := escpos.NewEscposPrinter(p)
        if err = uc.ExecuteTestImage(ep, payload); err != nil {
            log.Println("print:", err)
        }
        log.Println("Image printed, performing cut")
        ep.Cut()
        if err = ep.PrintAndFlush(); err != nil {
            log.Println("flush:", err)
        }
        log.Println("Flush OK")
        if err = p.EndDocument(); err != nil {
            log.Println("EndDocument:", err)
        }
        log.Println(strings.TrimSpace("Pencetakan Selesai"))
    })
	if err != nil {
		log.Fatalf("mqtt subscribe: %v", err)
	}
	broker := os.Getenv("MQTT_BROKER_URL")
	if broker == "" {
		host := os.Getenv("MQTT_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("MQTT_PORT")
		if port == "" {
			port = "1883"
		}
		broker = "tcp://" + host + ":" + port
	}
	log.Printf("Client ready | Broker: %s | Topic: %s | Printer: %s", broker, topic, settings.PrinterName)
	select {}
}

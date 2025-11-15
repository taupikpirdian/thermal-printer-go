package mqttclient

import (
    "os"
    "strconv"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber struct {
    Client mqtt.Client
}

func NewSubscriberFromEnv() (*Subscriber, error) {
    broker := os.Getenv("MQTT_BROKER_URL")
    if broker == "" {
        host := os.Getenv("MQTT_HOST")
        if host == "" { host = "localhost" }
        port := os.Getenv("MQTT_PORT")
        if port == "" { port = "1883" }
        broker = "tcp://" + host + ":" + port
    }
    user := os.Getenv("MQTT_USERNAME")
    pass := os.Getenv("MQTT_PASSWORD")
    if user == "" { user = os.Getenv("MQTT_AUTH_USERNAME") }
    if pass == "" { pass = os.Getenv("MQTT_AUTH_PASSWORD") }
    clientID := os.Getenv("MQTT_CLIENT_ID")
    if clientID == "" {
        clientID = "go-printer-" + strconv.FormatInt(time.Now().UnixNano(), 10)
    }
    opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID).SetUsername(user).SetPassword(pass).SetCleanSession(true).SetAutoReconnect(true)
    c := mqtt.NewClient(opts)
    t := c.Connect()
    t.Wait()
    if t.Error() != nil {
        return nil, t.Error()
    }
    return &Subscriber{Client: c}, nil
}

func (s *Subscriber) Subscribe(topic string, handler func([]byte)) error {
    t := s.Client.Subscribe(topic, 0, func(_ mqtt.Client, m mqtt.Message) { handler(m.Payload()) })
    t.Wait()
    if t.Error() != nil {
        return t.Error()
    }
    return nil
}

func (s *Subscriber) Disconnect() {
    if s.Client != nil {
        s.Client.Disconnect(250)
    }
}
package clients

import (
	"encoding/json"
	"log"
	"time"

	"github.com/kube-tarian/container-bridge/client/pkg/clickhouse"
	"github.com/kube-tarian/container-bridge/client/pkg/config"

	"github.com/nats-io/nats.go"
)

// constant variables to use with nats stream and
// nats publishing
const (
	streamSubjects string = "CONTAINERMETRICS.*"
	eventSubject   string = "CONTAINERMETRICS.git"
	eventConsumer  string = "container-event-consumer"
)

type NATSContext struct {
	conf     *config.Config
	conn     *nats.Conn
	stream   nats.JetStreamContext
	dbClient clickhouse.DBInterface
}

func NewNATSContext(conf *config.Config, dbClient clickhouse.DBInterface) (*NATSContext, error) {
	log.Println("Waiting before connecting to NATS at:", conf.NatsAddress)
	time.Sleep(1 * time.Second)

	conn, err := nats.Connect(conf.NatsAddress, nats.Name("Github metrics"), nats.Token(conf.NatsToken))
	if err != nil {
		return nil, err
	}

	ctx := &NATSContext{
		conf:     conf,
		conn:     conn,
		dbClient: dbClient,
	}

	stream, err := ctx.CreateStream()
	if err != nil {
		ctx.conn.Close()
		return nil, err
	}

	ctx.stream = stream
	ctx.Subscribe(eventSubject, eventConsumer, dbClient)

	_, err = stream.StreamInfo("CONTAINERMETRICS")
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (n *NATSContext) CreateStream() (nats.JetStreamContext, error) {
	// Creates JetStreamContext
	stream, err := n.conn.JetStream()
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (n *NATSContext) Close() {
	n.conn.Close()
	n.dbClient.Close()
}

func (n *NATSContext) Subscribe(subject string, consumer string, conn clickhouse.DBInterface) {
	n.stream.Subscribe(subject, func(msg *nats.Msg) {
		type events struct {
			Events []json.RawMessage `json:"events"`
		}

		eventDocker := &events{}
		err := json.Unmarshal(msg.Data, &eventDocker)
		if err == nil {
			log.Println(eventDocker)
			msg.Ack()
			repoName := msg.Header.Get("REPO_NAME")
			type newEvent struct {
				RepoName string          `json:"repoName"`
				Event    json.RawMessage `json:"event"`
			}

			for _, event := range eventDocker.Events {
				event := &newEvent{
					RepoName: repoName,
					Event:    event,
				}

				eventsJSON, err := json.Marshal(event)
				if err != nil {
					log.Printf("Failed to marshall with repo name going ahead with only event, %v", err)
					eventsJSON = msg.Data
				}
				conn.InsertEvent(string(eventsJSON))

				if n.conf.TemporalEnabled {
					// Currently temporal workflow is not yet supported
					// Once it is supported the idea here is to trigger job in the workflow
					log.Println(`Currently temporal workflow is not yet supported!!!\n 
					Once it is supported the idea here is to trigger job in the workflow for this event,
					 for example to retrieve more information for the event!!!`)
				}
			}
		} else {
			log.Printf("Failed to unmarshal event, %v", err)
			conn.InsertEvent(string(msg.Data))
		}

		log.Println("Inserted metrics:", string(msg.Data))
	}, nats.Durable(consumer), nats.ManualAck())
}

package definition

import "time"

//AMQPPublishing used for sending to AMQP server
type AMQPPublishing struct {
	URL        string `json:"url"`        // url to the amqp server e.g. amqp://guest:guest@localhost:5672/vhost
	RoutingKey string `json:"routingKey"` // the routing key for posting the message
	Delay      int    `json:"delay"`      // message send delay in seconds
	Exchange   string `json:"exchange"`   // the name of the exchange to post to
	Body       string `json:"body"`       // payload of the message

	// Properties
	ContentType     string    `json:"contentType"`     // MIME content type
	ContentEncoding string    `json:"contentEncoding"` // MIME content encoding
	Priority        uint8     `json:"priority"`        // 0 to 9
	CorrelationID   string    `json:"correlationId"`   // correlation identifier
	ReplyTo         string    `json:"replyTo"`         // address to to reply to (ex: RPC)
	Expiration      string    `json:"expiration"`      // message expiration spec
	MessageID       string    `json:"messageId"`       // message identifier
	Timestamp       time.Time `json:"timestamp"`       // message timestamp
	Type            string    `json:"type"`            // message type name
	UserID          string    `json:"userId"`          // creating user id - ex: "guest"
	AppID           string    `json:"appId"`           // creating application id
}

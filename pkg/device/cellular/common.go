package cellular

type CellModem interface {
}

type SMSMessage struct {
	outgoing           bool
	destinationAddress string
	sendingAddress     string
	messageBody        []byte //
	read               bool   // Modem/SIM tracks read/unread
	messageId          int    // Modem/SIM's message identifier (from e.g. AT+CMGL/AT+CMGR)
}

func CreateOutgoingSMS(destination string, body []byte) SMSMessage {
	return SMSMessage{
		outgoing:           true,
		destinationAddress: destination,
		messageBody:        body,

		// TODO: Should we keep the same structure for unwritten messages?
		// Some of these fields might make sense once applied to a given modem
		sendingAddress: "",
		read:           false,
		messageId:      -1,
	}
}

type SmsService interface {
	SendMessage(message SMSMessage) error
	GetAllMessages() ([]SMSMessage, error)
}

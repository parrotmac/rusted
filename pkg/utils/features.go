package utils

// TODO: This might get converted to injectable modules that share a common interface
// For now, handlers will just switch on the value
type Features struct {
	MockSerial bool
	MockGPS bool
	MockMQTT bool
}

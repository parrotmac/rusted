package transport

import "time"

type HttpConfig struct {
	SeverBaseURL   string
	DefaultTimeout time.Duration
}

type HttpWrapper struct {
	Config *HttpConfig
}

func (w *HttpWrapper) GetBaseURL() string {
	return w.Config.SeverBaseURL
}

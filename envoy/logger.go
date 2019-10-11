package envoy

import log "github.com/sirupsen/logrus"

type Logger struct{}

func (logger Logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
func (logger Logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

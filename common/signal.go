package common

import (
	log "github.com/thinkboy/log4go"
	"os"
	"os/signal"
	"syscall"
)

const (
	Ver = "0.2"
)

// InitSignal register signals handler.
func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Info("comet[%s] get a signal %s", Ver, s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			return
		default:
			return
		}
	}
}

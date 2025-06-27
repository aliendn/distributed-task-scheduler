package cluster

import (
	"log"
	"time"
)

// Heartbeater sends periodic signals to simulate node health
type Heartbeater struct {
	interval time.Duration
	stopChan chan struct{}
	NodeID   string
}

// NewHeartbeater creates a heartbeat sender
func NewHeartbeater(nodeID string, interval time.Duration) *Heartbeater {
	return &Heartbeater{
		interval: interval,
		stopChan: make(chan struct{}),
		NodeID:   nodeID,
	}
}

func (hb *Heartbeater) Start() {
	ticker := time.NewTicker(hb.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Printf("[Heartbeat] Node %s is alive", hb.NodeID)
			case <-hb.stopChan:
				log.Printf("[Heartbeat] Node %s stopped", hb.NodeID)
				return
			}
		}
	}()
}

func (hb *Heartbeater) Stop() {
	select {
	case <-hb.stopChan:
		// already closed
	default:
		close(hb.stopChan)
	}
}

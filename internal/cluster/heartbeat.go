package cluster

import (
	"log"
	"sync"
	"time"
)

// Heartbeater Heartbeat sends periodic liveness signals.
type Heartbeater struct {
	interval time.Duration
	stopChan chan struct{}
	NodeID   string
	once     sync.Once
}

// NewHeartbeater creates a heartbeat sender.
func NewHeartbeater(nodeID string, interval time.Duration) *Heartbeater {
	return &Heartbeater{
		interval: interval,
		stopChan: make(chan struct{}),
		NodeID:   nodeID,
	}
}

// Start begins sending heartbeats.
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

// Stop signals the heartbeat loop to exit.
func (hb *Heartbeater) Stop() {
	hb.once.Do(func() {
		close(hb.stopChan)
	})
}

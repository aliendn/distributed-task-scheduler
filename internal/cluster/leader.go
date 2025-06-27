package cluster

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// LeaderElector handles leader election and role switching
type LeaderElector struct {
	NodeID      string
	IsLeader    bool
	leaderMutex sync.RWMutex
	stopChan    chan struct{}
	callback    func() // Called when this node becomes leader
}

// NewLeaderElector creates a new leader instance
func NewLeaderElector(onLeadershipGained func()) *LeaderElector {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = generateRandomNodeID()
	}
	return &LeaderElector{
		NodeID:   nodeID,
		stopChan: make(chan struct{}),
		callback: onLeadershipGained,
	}
}

// Start begins the election loop
func (le *LeaderElector) Start() {
	go le.runElectionLoop()
	log.Printf("[Cluster] Node %s starting leader election loop", le.NodeID)
}

// Stop shuts down election
func (le *LeaderElector) Stop() {
	close(le.stopChan)
}

func (le *LeaderElector) runElectionLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-le.stopChan:
			return
		case <-ticker.C:
			le.electLeader()
		}
	}
}

// Simulated leader election
func (le *LeaderElector) electLeader() {
	// TODO: Replace with etcd/raft/redis lock
	isLeader := rand.Intn(3) == 1 // 1 in 3 chance
	le.setLeadership(isLeader)
}

func (le *LeaderElector) setLeadership(isLeader bool) {
	le.leaderMutex.Lock()
	defer le.leaderMutex.Unlock()

	if le.IsLeader != isLeader {
		le.IsLeader = isLeader
		if isLeader {
			log.Printf("[Leader] Node %s became leader", le.NodeID)
			le.callback()
		} else {
			log.Printf("[Leader] Node %s is now a follower", le.NodeID)
		}
	}
}

func (le *LeaderElector) IsCurrentLeader() bool {
	le.leaderMutex.RLock()
	defer le.leaderMutex.RUnlock()
	return le.IsLeader
}

func generateRandomNodeID() string {
	return "node-" + time.Now().Format("150405") // fallback ID
}

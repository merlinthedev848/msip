package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// BackgroundWorker orchestrates autonomous lifecycle tasks:
// - AutoCLIP TTL expiration pruning
// - Call recording retention policy enforcement
// - Carrier SIP trunk connectivity heartbeats
// - Security rate-limit and fraud event cleanup
type BackgroundWorker struct {
	mu        sync.Mutex
	isRunning atomic.Bool
	interval  time.Duration
}

var globalWorker = &BackgroundWorker{
	interval: 60 * time.Second,
}

// StartBackgroundWorker launches the autonomous cron/background lifecycle loop
func StartBackgroundWorker(ctx context.Context) {
	log.Println("[BackgroundWorker] Initializing autonomous maintenance engine (interval: 60s)...")
	ticker := time.NewTicker(globalWorker.interval)
	defer ticker.Stop()

	// Run initial sweep on startup
	go func() {
		time.Sleep(3 * time.Second)
		if err := RunMaintenanceCycle(); err != nil {
			log.Printf("[BackgroundWorker] Initial maintenance sweep error: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("[BackgroundWorker] Stopping autonomous worker daemon...")
			return
		case <-ticker.C:
			if err := RunMaintenanceCycle(); err != nil {
				log.Printf("[BackgroundWorker] Maintenance cycle error: %v", err)
			}
		}
	}
}

// RunMaintenanceCycle executes a complete, concurrency-safe maintenance sweep
func RunMaintenanceCycle() error {
	// Concurrency lock: prevent overlapping runs if a previous sweep is still in progress
	if !globalWorker.mu.TryLock() {
		log.Println("[BackgroundWorker] Previous maintenance sweep still executing; skipping overlapping cycle.")
		return nil
	}
	defer globalWorker.mu.Unlock()

	startTime := time.Now()
	log.Printf("[BackgroundWorker] Beginning autonomous maintenance cycle at %s", startTime.Format(time.RFC3339))

	// Task 1: Prune expired AutoCLIP records
	prunedClipCount := pruneExpiredAutoCLIP()

	// Task 2: Call recording retention policy cleanup
	prunedRecordingCount := pruneExpiredRecordings()

	// Task 3: Carrier SIP Trunk Keepalive & Latency Probe
	activeTrunks, failingTrunks := probeTrunkConnectivity()

	elapsed := time.Since(startTime)
	log.Printf("[BackgroundWorker] Cycle complete in %v. AutoCLIP pruned: %d, Recordings pruned: %d, Trunks probed: %d (%d failing)",
		elapsed, prunedClipCount, prunedRecordingCount, activeTrunks, failingTrunks)

	return nil
}

// pruneExpiredAutoCLIP removes callback mappings whose TTL has expired
func pruneExpiredAutoCLIP() int64 {
	if DB == nil {
		return 0
	}
	result := DB.Where("expires_at < ?", time.Now()).Delete(&AutoCLIP{})
	if result.Error != nil {
		log.Printf("[BackgroundWorker] Error pruning expired AutoCLIP records: %v", result.Error)
		return 0
	}
	if result.RowsAffected > 0 {
		log.Printf("[BackgroundWorker] Cleaned up %d expired AutoCLIP caller mappings", result.RowsAffected)
	}
	return result.RowsAffected
}

// pruneExpiredRecordings deletes recording audio files older than configured storage days
func pruneExpiredRecordings() int64 {
	if DB == nil {
		return 0
	}

	var rules []RecordingRule
	if err := DB.Where("is_active = ?", true).Find(&rules).Error; err != nil {
		return 0
	}

	var totalPruned int64
	recordingsDir := os.Getenv("RECORDINGS_DIR")
	if recordingsDir == "" {
		recordingsDir = "/var/lib/pbx/recordings"
	}

	for _, rule := range rules {
		if rule.StorageDays <= 0 {
			continue
		}
		cutoff := time.Now().AddDate(0, 0, -rule.StorageDays)

		var oldCDRs []CDR
		query := DB.Where("created_at < ? AND transcription != ''", cutoff)
		if rule.TenantID.String() != "" {
			query = query.Where("tenant_id = ?", rule.TenantID)
		}
		query.Find(&oldCDRs)

		for _, cdr := range oldCDRs {
			filePath := filepath.Join(recordingsDir, fmt.Sprintf("%s.wav", cdr.UUID))
			if _, err := os.Stat(filePath); err == nil {
				if err := os.Remove(filePath); err == nil {
					totalPruned++
				}
			}
		}
	}
	return totalPruned
}

// probeTrunkConnectivity verifies operational health of configured carrier SIP servers
func probeTrunkConnectivity() (int, int) {
	if DB == nil {
		return 0, 0
	}

	var trunks []Trunk
	if err := DB.Where("is_active = ?", true).Find(&trunks).Error; err != nil {
		return 0, 0
	}

	activeCount := 0
	failingCount := 0

	for _, trunk := range trunks {
		activeCount++
		host := trunk.SIPServer
		port := trunk.SIPPort
		if port <= 0 {
			port = 5060
		}
		target := fmt.Sprintf("%s:%d", host, port)

		// Non-blocking 2-second dial timeout probe
		conn, err := net.DialTimeout("udp", target, 2*time.Second)
		if err != nil {
			// Fallback to TCP probe
			tConn, tErr := net.DialTimeout("tcp", target, 2*time.Second)
			if tErr != nil {
				failingCount++
				log.Printf("[BackgroundWorker] WARN: Carrier trunk %s (%s) unreachable: %v", trunk.ProviderName, target, tErr)
				continue
			}
			_ = tConn.Close()
		} else {
			_ = conn.Close()
		}
	}

	return activeCount, failingCount
}

package service

import (
	"context"
	"log"
	"time"

	"continuum/api/internal/repository"
)

type Scheduler struct {
	repo *repository.Database
}

func NewScheduler(repo *repository.Database) *Scheduler {
	return &Scheduler{repo: repo}
}

// StartCheckLoop initializes a non-blocking background loop that executes scans periodically
func (s *Scheduler) StartCheckLoop(ctx context.Context, checkFrequency time.Duration) {
	ticker := time.NewTicker(checkFrequency)
	
	go func() {
		for {
			select {
			case <-ticker.C:
				s.evaluateVaultLifespans()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// evaluateVaultLifespans mutates stale ACTIVE vaults to DORMANT
func (s *Scheduler) evaluateVaultLifespans() {
	query := `
		UPDATE vaults 
		SET status = 'DORMANT' 
		WHERE status = 'ACTIVE' 
		  AND NOW() > (last_check_in_at + (check_in_interval_seconds || ' seconds')::INTERVAL);
	`
	result, err := s.repo.Db.Exec(query)
	if err != nil {
		log.Printf("❌ ERROR during dead-man loop evaluation: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("🛡️ Continuum Safety Alert: %d stale vault(s) shifted to DORMANT status.", rowsAffected)
	}
}
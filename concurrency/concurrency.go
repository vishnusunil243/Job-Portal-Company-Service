package concurrency

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Concurrency struct {
	DB *gorm.DB
	mu sync.Mutex
}

func NewConcurrency(DB *gorm.DB) *Concurrency {
	return &Concurrency{
		DB: DB,
	}

}
func (c *Concurrency) Concurrency() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			c.mu.Lock()
			if err := c.DB.Exec(`
			UPDATE jobs SET status_id=3
			WHERE id IN (
				SELECT j.id FROM jobs j WHERE 
				NOW() > j.valid_until 
				AND j.status_id NOT IN (3)
			)	
			`).Error; err != nil {
				log.Print("error while performing crone jobs ", err)
				break
			}
		}
		c.mu.Unlock()
	}()
	fmt.Println("worked")
}

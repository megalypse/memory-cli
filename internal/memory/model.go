package memory

import "time"

type Memory struct {
	ID        int
	GroupID   int
	Name      string
	Content   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

package memory_group

import "time"

type MemoryGroup struct {
	ID          int
	Name        string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

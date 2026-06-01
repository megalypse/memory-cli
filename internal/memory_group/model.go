package memory_group

import "time"

type MemoryGroup struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
	DeletedAt   time.Time
}

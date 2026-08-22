package routine

import "time"

type Routine struct {
	ID          int64
	Name        string
	Description string
	Category    string
	Weight      int
	Coefficient int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

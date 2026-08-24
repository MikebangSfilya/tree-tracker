package routine

import "time"

type Routine struct {
	ID          int64
	Category    string
	Type        string
	Weight      int
	Coefficient int
	Temporary   bool
	Disposable  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

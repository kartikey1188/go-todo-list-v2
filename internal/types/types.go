package types

import (
	"fmt"
	"time"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Deadline    Date   `json:"deadline"`
}

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)

	s = s[1 : len(s)-1]

	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))
	return []byte(formatted), nil
}

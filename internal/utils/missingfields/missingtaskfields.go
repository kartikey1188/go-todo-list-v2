package missingfields

import (
	"fmt"
	"strings"

	"github.com/kartikey1188/go-todo-list-v2/internal/types"
)

func MissingTaskFields(task types.Task) string {
	missingFields := []string{}

	if task.Title == "" {
		missingFields = append(missingFields, "title")
	}
	if task.Description == "" {
		missingFields = append(missingFields, "description")
	}

	if len(missingFields) == 0 {
		return "" // No missing fields
	}

	return fmt.Sprintf("Missing Fields: %s", strings.Join(missingFields, ", "))
}

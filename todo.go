package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aquasecurity/table"
)

type TodoStatus string

const (
	StatusNotStarted TodoStatus = "Not Started"
	StatusInProgress TodoStatus = "In Progress"
	StatusCompleted  TodoStatus = "Completed"
	StatusRejected   TodoStatus = "Rejected"
)

var statusEmojiMap = map[TodoStatus]string{
	StatusInProgress: "🚧",
	StatusNotStarted: "⏳",
	StatusCompleted:  "✅",
	StatusRejected:   "❌",
}

type TodoPriority string

const (
	PriorityHigh   TodoPriority = "High"
	PriorityMedium TodoPriority = "Medium"
	PriorityLow    TodoPriority = "Low"
)

type Todo struct {
	Id            int
	Task          string
	Status        TodoStatus
	Priority      TodoPriority
	AssignedTo    string
	DueDate       *time.Time
	EstimatedTime *time.Duration
	TimeSpent     *time.Duration
	CompletedAt   *time.Time
	DependsOn     []int
}

type Todos []Todo

// IsComplete returns true if the task is marked as complete
func (t *Todo) IsComplete() bool {
	return t.Status == StatusCompleted
}

// CanStart checks if all dependencies are completed
func (t *Todo) CanStart(todos Todos) bool {
	// If no dependencies, can start immediately
	if len(t.DependsOn) == 0 {
		return true
	}

	// Check if all dependencies are completed
	for _, depID := range t.DependsOn {
		depTask, err := todos.FindById(depID)
		if err != nil || !depTask.IsComplete() {
			return false
		}
	}
	return true
}

// MarkComplete marks the task as completed and sets the completion timestamp
func (t *Todo) MarkComplete() {
	t.Status = StatusCompleted
	now := time.Now()
	t.CompletedAt = &now
}

// FindById finds a todo by its ID
func (todos Todos) FindById(id int) (*Todo, error) {
	for i := range todos {
		if todos[i].Id == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("task not found")
}

func (todos Todos) FindIndexById(id int) (int, error) {
	for i := range todos {
		if todos[i].Id == id {
			return i, nil
		}
	}
	return -1, errors.New("task not found")
}

func ParseStatus(status string) (TodoStatus, error) {
	cleanStatus := cleanStatusEmoji(status)

	// Normalize the status string
	switch strings.ToLower(strings.TrimSpace(cleanStatus)) {
	case "in progress":
		return StatusInProgress, nil
	case "complete", "completed", "done":
		return StatusCompleted, nil
	case "not started", "pending", "todo":
		return StatusNotStarted, nil
	case "reject", "rejected":
		return StatusRejected, nil
	default:
		return "", fmt.Errorf("invalid status: %s", status)
	}
}

// cleanStatusEmoji removes emoji symbols from status text
func cleanStatusEmoji(status string) string {
	// Common emojis used in status
	emojiPatterns := []string{"✅", "✓", "🟢", "🔄", "🚧", "🟡", "⏳", "❌", "🚫", "🔴", "⛔"}

	result := status
	for _, emoji := range emojiPatterns {
		result = strings.ReplaceAll(result, emoji, "")
	}

	return strings.TrimSpace(result)
}

func (todos *Todos) Add(title string, assignedTo string) {
	// Find the highest ID to ensure unique IDs
	highestID := 0
	for _, t := range *todos {
		if t.Id > highestID {
			highestID = t.Id
		}
	}

	todo := Todo{
		Id:            highestID + 1, // Assign next available ID
		Task:          title,
		Status:        StatusNotStarted,
		Priority:      PriorityLow,
		AssignedTo:    assignedTo,
		DueDate:       nil,
		EstimatedTime: nil,
		TimeSpent:     nil,
		CompletedAt:   nil,
		DependsOn:     nil,
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) Delete(id int) error {
	index, err := todos.FindIndexById(id)
	if err != nil {
		return err
	}

	*todos = slices.Delete(*todos, index, index+1)

	return nil
}

func (todos *Todos) EditTitle(id int, title string) error {
	index, err := todos.FindIndexById(id)
	if err != nil {
		return err
	}

	(*todos)[index].Task = title

	return nil
}

func (todos *Todos) UpdateStatus(id int, status TodoStatus) error {
	index, err := todos.FindIndexById(id)
	if err != nil {
		return err
	}

	if status == StatusCompleted {
		now := time.Now()
		(*todos)[index].CompletedAt = &now
	} else if (*todos)[index].Status == StatusCompleted {
		(*todos)[index].CompletedAt = nil
	}

	(*todos)[index].Status = status

	return nil
}

func (todos *Todos) Print() {
	table := table.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("Id", "Task", "Status", "Priority", "Assigned To", "Due Date", "Estimated Time", "Time Spent", "Completed At", "Depends On")
	for _, t := range *todos {
		dependsOn := ""
		if len(t.DependsOn) > 0 {
			for i, dep := range t.DependsOn {
				if i > 0 {
					dependsOn += ", "
				}
				dependsOn += strconv.Itoa(dep)
			}
		}

		dueDate := ""
		if t.DueDate != nil {
			dueDate = t.DueDate.Format(time.DateOnly)
		}

		estimatedTime := ""
		if t.EstimatedTime != nil {
			estimatedTime = t.EstimatedTime.String()
		}

		timeSpent := ""
		if t.TimeSpent != nil {
			timeSpent = t.TimeSpent.String()
		}

		completedAt := ""
		if t.CompletedAt != nil {
			completedAt = t.CompletedAt.Format(time.DateOnly)
		}

		statusEmoji := statusEmojiMap[t.Status]

		table.AddRow(
			strconv.Itoa(t.Id),
			t.Task,
			fmt.Sprintf("%s %s", statusEmoji, string(t.Status)),
			string(t.Priority),
			"@"+t.AssignedTo,
			dueDate,
			estimatedTime,
			timeSpent,
			completedAt,
			dependsOn,
		)
	}

	table.Render()
}

package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Entity[ID comparable] interface {
	GetID() ID
}

type Status string

const (
	StatusNew             Status = "Новое"
	StatusReopened        Status = "Переоткрыто"
	StatusInWork          Status = "В работе"
	StatusOnReview        Status = "На ревью"
	StatusReadyForMerge   Status = "Готово ко влитию"
	StatusMerged          Status = "Влито"
	StatusStopped         Status = "Приостановлено"
	StatusInClarification Status = "В уточнении"
	StatusTesting         Status = "Тестирование"
	StatusTestInProd      Status = "Протестировать в проде"
	StatusFinished        Status = "Завершено"
)

var (
	allStatuses = []Status{
		StatusNew,
		StatusReopened,
		StatusInWork,
		StatusOnReview,
		StatusReadyForMerge,
		StatusMerged,
		StatusStopped,
		StatusInClarification,
		StatusTesting,
		StatusTestInProd,
		StatusFinished,
	}
	statusIndex       = make(map[Status]int)
	ErrStatusNotValid = errors.New("Невалидный статус")
	ErrNoNextStatus   = errors.New("Следующего статуса нет!")
	ErrNoPrevStatus   = errors.New("Предыдущего статуса нет!")
)

func init() {
	for i, st := range allStatuses {
		statusIndex[st] = i + 1
	}
}

func (s Status) IsValid() bool {
	return statusIndex[s] != 0
}

func (s Status) String() string {
	return string(s)
}

func (s Status) Next() (Status, error) {
	idx := statusIndex[s]

	if idx == 0 {
		return "", ErrStatusNotValid
	}

	if len(allStatuses) == idx {
		return "", ErrNoNextStatus
	}

	return allStatuses[idx], nil
}

func (s Status) Prev() (Status, error) {
	idx := statusIndex[s]
	if idx == 0 {
		return "", ErrStatusNotValid
	}
	if idx == 1 {
		return "", ErrNoPrevStatus
	}
	return allStatuses[idx-2], nil
}

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	Status      Status
}

func (t Task) GetID() uuid.UUID {
	return t.ID
}

func (t Task) Validate() error {
	if t.ID == uuid.Nil {
		return fmt.Errorf("task id is required")
	}

	if t.Title == "" {
		return fmt.Errorf("task title is required")
	}

	if !t.Status.IsValid() {
		return fmt.Errorf("invalid task status: %q", t.Status)
	}

	return nil
}

package entities

import "github.com/satori/go.uuid"

// EventID уникальный идентификатор события
// По сути это value object
type EventID struct {
	value uuid.UUID
}

// NewEventID конструктор идентификатора события
func NewEventID(id string) (EventID, error) {
	var err error
	var IDObj uuid.UUID

	if id == "" {
		IDObj = uuid.NewV4()
	} else {
		IDObj, err = uuid.FromString(id)
	}

	return EventID{value: IDObj}, err
}

func (u EventID) String() string {
	return u.value.String()
}

// Equal сравнение ID с другим. Возвращает true если идентификаторы равны
func (u EventID) Equal(other EventID) bool {
	return u.value.String() == other.value.String()
}

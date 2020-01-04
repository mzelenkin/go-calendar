package inmemory

import (
	"context"
	"github.com/mzelenkin/go-calendar/internal/domain/entities"
	"github.com/mzelenkin/go-calendar/internal/storage"
	"time"
)

// EventInMemoryStorage хранилище пользователей в памяти
type EventInMemoryStorage struct {
	data map[string]entities.Event
}

func NewEventInMemoryStorage() (*EventInMemoryStorage, error) {
	return &EventInMemoryStorage{
		data: map[string]entities.Event{},
	}, nil
}

func (i *EventInMemoryStorage) Create(ctx context.Context, event *entities.Event) error {
	id := event.ID.String()

	if _, ok := i.data[id]; ok {
		return storage.EntityAlreadyExists
	}

	i.data[id] = *event

	return nil
}

func (i *EventInMemoryStorage) Update(ctx context.Context, event *entities.Event) error {
	eventIDString := event.ID.String()
	if _, ok := i.data[eventIDString]; !ok {
		return storage.EntityNotFound
	}

	i.data[eventIDString] = *event

	return nil
}

func (i *EventInMemoryStorage) ListAll(ctx context.Context, page int, itemsPerPage int) ([]entities.Event, error) {
	var ret []entities.Event
	var counter int
	startRec := page * itemsPerPage
	endRec := page*itemsPerPage + itemsPerPage

	for _, v := range i.data {

		// Пагинация
		if counter < startRec {
			continue
		}

		if counter >= endRec {
			break
		}

		counter++

		ret = append(ret, v)
	}

	return ret, nil
}

// FindBySpan возвращает все события, лежащие в указанном временном диапазоне start..end
// Параметры page и itemsPerPage
func (i *EventInMemoryStorage) FindBySpan(ctx context.Context, start time.Time, end time.Time, page int, itemsPerPage int) ([]entities.Event, error) {
	var ret []entities.Event
	var counter int
	startRec := page * itemsPerPage
	endRec := page*itemsPerPage + itemsPerPage

	for _, v := range i.data {

		// Пагинация
		if counter < startRec {
			continue
		}

		if counter >= endRec {
			break
		}

		counter++

		if overlaps(start, end, v.Start, v.End) {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (i *EventInMemoryStorage) FindByID(ctx context.Context, id entities.EventID) (*entities.Event, error) {
	eventIDString := id.String()
	item, ok := i.data[eventIDString]
	if !ok {
		return nil, storage.EntityNotFound
	}

	return &item, nil
}

func (i *EventInMemoryStorage) DeleteByID(ctx context.Context, id *entities.EventID) error {
	eventIDString := id.String()
	if _, ok := i.data[eventIDString]; !ok {
		return storage.EntityNotFound
	}

	delete(i.data, eventIDString)

	return nil
}

// overlaps проверяет пересечение двух диапазонов дат
func overlaps(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}

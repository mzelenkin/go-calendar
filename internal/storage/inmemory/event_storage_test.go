package inmemory

import (
	"context"
	"github.com/mzelenkin/go-calendar/internal/domain/entities"
	"github.com/mzelenkin/go-calendar/internal/storage"
	"reflect"
	"testing"
	"time"
)

func TestEventInMemoryStorage_Create(t *testing.T) {
	s, err := NewEventInMemoryStorage()
	if err != nil {
		t.Error(err)
	}

	// Создаем новое ID события
	id, err := entities.NewEventID("")
	if err != nil {
		t.Error(err)
	}

	event := entities.Event{
		ID:          id,
		Title:       "Событие #1",
		Start:       time.Now(),
		End:         time.Now().Add(2 * time.Hour),
		Description: "Тестовое событие",
	}

	err = s.Create(context.Background(), &event)
	if err != nil {
		t.Error(err)
	}

	event1, err := s.FindByID(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	eq := reflect.DeepEqual(&event, event1)
	if !eq {
		t.Fail()
	}
}

func TestEventInMemoryStorage_Update(t *testing.T) {
	s, err := NewEventInMemoryStorage()
	if err != nil {
		t.Error(err)
	}

	// Создаем пользователя
	id, err := entities.NewEventID("")
	if err != nil {
		t.Error(err)
	}

	event := entities.Event{
		ID:          id,
		Title:       "Событие #2",
		Start:       time.Now(),
		End:         time.Now().Add(1 * time.Hour),
		Description: "Тестовое событие",
	}

	err = s.Create(context.Background(), &event)
	if err != nil {
		t.Error(err)
	}

	// Модифицируем и обновляем в хранилище

	err = s.Update(context.Background(), &event)
	if err != nil {
		t.Error(err)
	}

	// Извлекаем и проверяем как обновилось
	event1, err := s.FindByID(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	eq := reflect.DeepEqual(&event, event1)
	if !eq {
		t.Fail()
	}

	err = s.DeleteByID(context.Background(), &event.ID)
	if err != nil {
		t.Error(err)
	}
	_, err = s.FindByID(context.Background(), id)
	if err != storage.EntityNotFound {
		t.Fail()
	}
}

// TestEventInMemoryStorage_GetByRange проверяет получение событий по диапазону
func TestEventInMemoryStorage_GetByRange(t *testing.T) {
	ctx := context.Background()
	s, err := NewEventInMemoryStorage()
	if err != nil {
		t.Error(err)
	}

	id, err := entities.NewEventID("")
	if err != nil {
		t.Error(err)
	}

	if s.Create(ctx, &entities.Event{
		ID:          id,
		Title:       "Событие #1",
		Start:       time.Now(),
		End:         time.Now().Add(1 * time.Hour),
		Description: "Тестовое событие на 1 час",
	}) != nil {
		t.Fatal(err)
	}

	id, err = entities.NewEventID("")
	if err != nil {
		t.Error(err)
	}

	if s.Create(ctx, &entities.Event{
		ID:          id,
		Title:       "Событие #2",
		Start:       time.Now().Add(1*time.Hour + 5*time.Minute),
		End:         time.Now().Add(1*time.Hour + 30*time.Minute),
		Description: "Тестовое событие на 25 минут",
	}) != nil {
		t.Fatal(err)
	}

	id, err = entities.NewEventID("")
	if err != nil {
		t.Error(err)
	}
	if s.Create(ctx, &entities.Event{
		ID:          id,
		Title:       "Событие #3",
		Start:       time.Now().Add(2 * time.Hour),
		End:         time.Now().Add(2*time.Hour + 10*time.Minute),
		Description: "Тестовое событие на 10 минут",
	}) != nil {
		t.Fatal(err)
	}

	// При выборке за час, в результат
	events, err := s.FindBySpan(ctx, time.Now(), time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fail()
	}

	// При выборе за 1:15 в выборку попадет 2 события
	events, err = s.FindBySpan(ctx, time.Now(), time.Now().Add(1*time.Hour+15*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fail()
	}
}

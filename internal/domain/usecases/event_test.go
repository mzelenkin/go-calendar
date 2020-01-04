package usecases

import (
	"context"
	"github.com/mzelenkin/go-calendar/internal/domain/entities"
	storage2 "github.com/mzelenkin/go-calendar/internal/storage"
	"github.com/mzelenkin/go-calendar/internal/storage/inmemory"
	"testing"
	"time"
)

// TestEventUsecases_CreateAndDelete проверяет создание и удаление
func TestEventUsecases_CreateAndDelete(t *testing.T) {
	ctx := context.Background()
	// inmemory хранилище не возвращает ошибок инициализации
	// По хорошему, тут должен быть mock, но пока лень 8(
	storage, _ := inmemory.NewEventInMemoryStorage()
	usecase := NewEventUsecases(storage)

	event := CreateEventRequest{
		Title:       "Тестовое событие №1",
		Start:       time.Now(),
		End:         time.Now(),
		Description: "",
	}

	id, err := usecase.Create(ctx, &event)
	if err != nil {
		t.Fatal(err)
	}

	eventId, err := entities.NewEventID(id)
	if err != nil {
		t.Fatal(err)
	}

	savedEvent, err := storage.FindByID(ctx, eventId)
	if err != nil {
		t.Fatal(err)
	}

	if event.Title != savedEvent.Title ||
		event.Start != savedEvent.Start ||
		event.End != savedEvent.End ||
		event.Description != savedEvent.Description {
		t.Fail()
	}

	err = usecase.Delete(ctx, eventId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = storage.FindByID(ctx, eventId)
	if err != storage2.EntityNotFound {
		t.Fail()
	}
}

// TestEventUsecases_Update проверяет обновление
func TestEventUsecases_Update(t *testing.T) {
	ctx := context.Background()

	storage, _ := inmemory.NewEventInMemoryStorage()
	usecase := NewEventUsecases(storage)

	event := CreateEventRequest{
		Title:       "Тестовое событие №1",
		Start:       time.Now(),
		End:         time.Now(),
		Description: "Это описание события №1",
	}

	id, err := usecase.Create(ctx, &event)
	if err != nil {
		t.Fatal(err)
	}

	eventId, err := entities.NewEventID(id)
	if err != nil {
		t.Fatal(err)
	}

	updatedEvent := UpdateEventRequest{
		ID:          eventId.String(),
		Title:       "Обновленное событие №1",
		Start:       time.Now(),
		End:         time.Now(),
		Description: "Это описание обновленного события",
	}

	err = usecase.Update(ctx, &updatedEvent)
	if err != nil {
		t.Fatal(err)
	}

	savedEvent, err := storage.FindByID(ctx, eventId)
	if err != nil {
		t.Fatal(err)
	}

	if updatedEvent.Title != savedEvent.Title ||
		updatedEvent.Start != savedEvent.Start ||
		updatedEvent.End != savedEvent.End ||
		updatedEvent.Description != savedEvent.Description {
		t.Fail()
	}
}

// TestEventUsecases_Crosses проверяет пересекающиеся события
func TestEventUsecases_Crosses(t *testing.T) {
	var err error
	ctx := context.Background()

	storage, _ := inmemory.NewEventInMemoryStorage()
	usecase := NewEventUsecases(storage)

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №1",
		Start:       time.Now(),
		End:         time.Now().Add(1 * time.Hour),
		Description: "Это описание обновленного события",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Пересекающееся событие",
		Start:       time.Now().Add(30 * time.Minute),
		End:         time.Now().Add(1 * time.Hour),
		Description: "Это описание обновленного события",
	})

	if err != ErrorDateBusy {
		t.Error(err)
		t.Fail()
	}
}

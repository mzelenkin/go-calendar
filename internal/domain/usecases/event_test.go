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
	// По хорошему, тут должен быть mock, но некогда
	// и у нас как раз есть покрытое тестами inmemory хранилище
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

// TestEventUsecases_Lists длинный длинный тест, проверяющий методы выборки за день, неделю и месяц
// В будущем надо будет его распилить на несколько и вообще тестить ByRange
func TestEventUsecases_Lists(t *testing.T) {
	var err error
	ctx := context.Background()

	storage, _ := inmemory.NewEventInMemoryStorage()
	usecase := NewEventUsecases(storage)

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №1",
		Start:       bod(time.Now()),
		End:         bod(time.Now()).Add(1 * time.Hour),
		Description: "Это описание события номер 1",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №1",
		Start:       bod(time.Now()).Add(19 * time.Hour),
		End:         bod(time.Now()).Add(22 * time.Hour),
		Description: "Это описание события номер 2",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №3",
		Start:       bod(time.Now()).Add(1*time.Hour).AddDate(0, 0, 1),
		End:         bod(time.Now()).Add(5*time.Hour).AddDate(0, 0, 1),
		Description: "Это описание события за пределами дня",
	})
	if err != nil {
		t.Fatal(err)
	}

	events, err := usecase.ListDay(ctx, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 2 {
		t.Fail()
	}

	// ============================
	// Чтобы два раза не создавать
	// Проверка выборки за неделю
	weekStart, weekEnd := weekRange(time.Now())
	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №4",
		Start:       weekStart.Add(1*time.Hour).AddDate(0, 0, 1),
		End:         weekStart.Add(5*time.Hour).AddDate(0, 0, 1),
		Description: "Это описание события в начале недели",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №5",
		Start:       weekEnd.Add(1*time.Hour).AddDate(0, 0, 1),
		End:         weekEnd.Add(2*time.Hour).AddDate(0, 0, 1),
		Description: "Это описание события за пределами недели",
	})
	if err != nil {
		t.Fatal(err)
	}

	weekEvents, err := usecase.ListWeek(ctx, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	// Вспоминаем, что выше тоже есть события на этой неделе
	// Итого должно получиться 4 события из 5
	if len(weekEvents) != 4 {
		t.Fail()
	}

	// ============================
	// Чтобы два раза не создавать
	// Проверка выборки за месяц
	now := time.Now()
	year := now.Year()
	month := now.Month()
	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №6",
		Start:       time.Date(year, month, 1, 01, 25, 0, 0, time.Local),
		End:         time.Date(year, month, 1, 02, 00, 0, 0, time.Local),
		Description: "Это описание события в начале месяца",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = usecase.Create(ctx, &CreateEventRequest{
		Title:       "Событие №7",
		Start:       time.Date(year, month, 1, 01, 25, 0, 0, time.Local).AddDate(0, 1, 0),
		End:         time.Date(year, month, 1, 02, 00, 0, 0, time.Local).AddDate(0, 1, 0),
		Description: "Это описание события в начале следующего месяца",
	})
	if err != nil {
		t.Fatal(err)
	}

	monthEvents, err := usecase.ListMonth(ctx, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	// Вспоминаем, что выше тоже есть события на этой неделе
	// Итого должно получиться 6 события из 7
	if len(monthEvents) != 6 {
		t.Fail()
	}
}

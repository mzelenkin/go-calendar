package usecases

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/mzelenkin/go-calendar/internal/domain/entities"
	"time"
)

// CreateEventRequest это DTO с входными данными для создания объекта Событие
type CreateEventRequest struct {
	Title       string    `validate:"required,min=3,max=50"`
	Start       time.Time `validate:"required,ltfield=End"`
	End         time.Time `validate:"required,gtfield=Start"`
	Description string
}

// UpdateEventRequest это DTO с входными данными для изменения объекта Событие целиком
type UpdateEventRequest struct {
	ID          string    `validate:"required,uuid"`
	Title       string    `validate:"required,min=3,max=50"`
	Start       time.Time `validate:"required,ltfield=End"`
	End         time.Time `validate:"required,gtfield=Start"`
	Description string
}

// EventUsecases сценарии использования для события
// При расширении эта структура может превратиться в фасад к use case'ам
type EventUsecases struct {
	storage EventStorage
}

func NewEventUsecases(storage EventStorage) *EventUsecases {
	return &EventUsecases{storage: storage}
}

// Create создает событие и возвращает его ID
func (u EventUsecases) Create(ctx context.Context, data *CreateEventRequest) (string, error) {
	validate := validator.New()

	err := validate.StructCtx(ctx, data)
	if err != nil {
		return "", err
	}

	id, err := entities.NewEventID("")
	if err != nil {
		return "", err
	}

	event := entities.Event{
		ID:          id,
		Title:       data.Title,
		Start:       data.Start,
		End:         data.End,
		Description: data.Description,
	}

	err = u.storage.Create(ctx, &event)

	return id.String(), err
}

// Update обновляет все событие целиком
func (u EventUsecases) Update(ctx context.Context, data *UpdateEventRequest) error {
	validate := validator.New()

	err := validate.StructCtx(ctx, data)
	if err != nil {
		return err
	}

	// Создаем объект EventID
	id, err := entities.NewEventID(data.ID)
	if err != nil {
		return err
	}

	event := entities.Event{
		ID:          id,
		Title:       data.Title,
		Start:       data.Start,
		End:         data.End,
		Description: data.Description,
	}

	err = u.storage.Update(ctx, &event)

	return err
}

// Delete удаляет сущность Событие по ее идентификатору
//
// Казалось бы зачем тут этот usecase, когда можно запросить напрямую хранилище?
// Дело в том, что во-первых это не удобно, т.е. нужно как-то запрашивать еще и
// хранилище. Но основная причина в том, что при расширении нам может понадобиться
// производить какие-то дополнительные действия и тогда мы можем добавить их сюда,
// не затрагивая остальной код.
func (u EventUsecases) Delete(ctx context.Context, id entities.EventID) error {
	return u.storage.DeleteByID(ctx, &id)
}

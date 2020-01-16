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

type ListResponseItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
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

	// В общем-то не важно, сколько тут повится записей, т.к. по идее либо что-то есть и тогда ошибка, либо их нет.
	item, err := u.storage.FindBySpan(ctx, data.Start, data.End)
	if err != nil {
		return "", err
	}

	if len(item) > 0 {
		return "", ErrorDateBusy
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

	item, err := u.storage.FindBySpan(ctx, data.Start, data.End)
	if err != nil {
		return err
	}

	// Если есть 1 событие и это не наше или событий больше 1, генерируем ошибку
	if len(item) == 1 && !id.Equal(item[0].ID) || len(item) > 1 {
		return ErrorDateBusy
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

// ListDay возвращает список событий за указанный день
func (u EventUsecases) ListDay(ctx context.Context, day time.Time) ([]ListResponseItem, error) {
	start := bod(day)
	end := eod(day)

	ret, err := u.findBySpan(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// ListWeek возвращает список событий за указанную неделю
func (u EventUsecases) ListWeek(ctx context.Context, day time.Time) ([]ListResponseItem, error) {
	start, end := weekRange(day)

	ret, err := u.findBySpan(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// ListMonth возвращает список событий за указанный месяц
func (u EventUsecases) ListMonth(ctx context.Context, day time.Time) ([]ListResponseItem, error) {
	start, end := monthRange(day)

	ret, err := u.findBySpan(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// findBySpan ищет события из указанного промежутка и маппит их на DTO
func (u EventUsecases) findBySpan(ctx context.Context, start time.Time, end time.Time) ([]ListResponseItem, error) {
	items, err := u.storage.FindBySpan(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var ret []ListResponseItem

	// Заполняем DTO
	for _, v := range items {
		ret = append(ret, ListResponseItem{
			ID:          v.ID.String(),
			Title:       v.Title,
			Start:       v.Start,
			End:         v.End,
			Description: v.Description,
		})
	}
	return ret, nil
}

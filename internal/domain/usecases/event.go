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

type ListAllResponseItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
}

// EventUsecases сценарии использования для события
// При расширении эта структура может превратиться в фасад к use case'ам
type EventUsecases struct {
	storage      EventStorage
	ItemsPerPage int
}

func NewEventUsecases(storage EventStorage) *EventUsecases {
	return &EventUsecases{storage: storage, ItemsPerPage: 25}
}

// Create создает событие и возвращает его ID
func (u EventUsecases) Create(ctx context.Context, data *CreateEventRequest) (string, error) {
	validate := validator.New()

	err := validate.StructCtx(ctx, data)
	if err != nil {
		return "", err
	}

	// В общем-то не важно, сколько тут повится записей, т.к. по идее либо что-то есть и тогда ошибка, либо их нет.
	item, err := u.storage.FindBySpan(ctx, data.Start, data.End, 0, 2)
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

	item, err := u.storage.FindBySpan(ctx, data.Start, data.End, 0, 2)
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

// ListAll возвращает список
func (u EventUsecases) ListAll(ctx context.Context, page int) ([]ListAllResponseItem, error) {
	items, err := u.storage.ListAll(ctx, page, u.ItemsPerPage)
	if err != nil {
		return nil, err
	}

	var ret []ListAllResponseItem

	// Заполняем DTO
	for _, v := range items {
		ret = append(ret, ListAllResponseItem{
			ID:          v.ID.String(),
			Title:       v.Title,
			Start:       v.Start,
			End:         v.End,
			Description: v.Description,
		})
	}

	return ret, nil
}

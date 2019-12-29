package usecases

import (
	"context"
	"github.com/mzelenkin/go-calendar/internal/domain/entities"
)

// EventStorage интерфейс хранилища событий, требующийся usecas'ам
// Мы не разделяем его на более мелкие (см. interface segregation)
// т.к. почти всегда используются все CRUD операции, однако по мере роста usecase'ов может понадобится разбиение
type EventStorage interface {
	Create(ctx context.Context, event *entities.Event) error
	Update(ctx context.Context, event *entities.Event) error

	FindByID(ctx context.Context, id entities.EventID) (*entities.Event, error)
	DeleteByID(ctx context.Context, id *entities.EventID) error
}

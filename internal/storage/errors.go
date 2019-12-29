package storage

// Ошибки хранилища
// Т.к. у реальных БД и т.п. ошибки разные, нам необходимо привести их к какой-то одной системе.

// NotFound - ошибка сущность не найдена
const EntityNotFound = StorageError("Entity not found")
const EntityAlreadyExists = StorageError("Entity already exists")

// StorageError тип для ошибок репозитория
type StorageError string

// Error реализует интерфейс error
func (r StorageError) Error() string {
	return string(r)
}

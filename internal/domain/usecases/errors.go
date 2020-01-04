package usecases

const ErrorDateBusy = UsecaseError("date busy")

// UsecaseError тип для ошибок сценария использования
type UsecaseError string

// Error реализует интерфейс error
func (r UsecaseError) Error() string {
	return string(r)
}

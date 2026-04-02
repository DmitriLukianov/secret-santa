package helpers

// ToPtr возвращает указатель на значение
func ToPtr[T any](v T) *T {
	return &v
}

// FromPtr возвращает значение по указателю или значение по умолчанию
func FromPtr[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

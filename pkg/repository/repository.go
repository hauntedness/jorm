package repository

type Repository[T any, K any] interface {
	FindById(K) (T, error)
	// insert([]T) (int, error)
	// update([]T) (int, error)
}

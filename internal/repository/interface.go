package repository

type Repository interface {
	Set(url string) string
	Get(id string) (string, error)
	Close() error
}

func NewRepository(filepath string) Repository {
	if len(filepath) == 0 {
		return newInMemoryRepository()
	}
	return newFileRepository(filepath)
}

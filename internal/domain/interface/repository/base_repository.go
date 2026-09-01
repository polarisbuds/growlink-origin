package repository_interface

//go:generate mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type BaseRepository[T any] interface {
	// Is exists record
	Exists(query string, args ...any) (bool, error)
	// get records
	Where(order, query string, args ...any) ([]T, error)
	// get record limit 1
	First(order string, query string, args ...any) (*T, error)
	// get all records
	All(order string) ([]T, error)
	// count records
	Count(query string, args ...any) (int64, error)

	// Insert record
	Insert(model *T) error
	// Insert records in batch
	InsertInBatch(models []T, batchSize int) error
	// soft delete records
	Delete(limit int, query string, args ...any) (int, error)
	// Update
	Update(name string, value any, query string, args ...any) (int, error)
}

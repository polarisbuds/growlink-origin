package interfaces

type BaseUnScopedRepository[T any] interface {
	// [Unscoped] Is exists record
	UExists(query string, args ...any) (bool, error)
	// [Unscoped] get records
	UWhere(order, query string, args ...any) ([]T, error)
	// [Unscoped] get record limit 1
	UFirst(order, query string, args ...any) (*T, error)
	// [Unscoped] get all records
	UAll(order string) ([]T, error)
	// [Unscoped] count records
	UCount(query string, args ...any) (int64, error)
}

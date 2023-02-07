package repository

type Repository interface {
	SaveValue(hashCode string, value interface{})
	GetValue(hash string) (string, error)
}

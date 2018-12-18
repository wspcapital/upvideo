package session

type Repository interface {
	FindById(id string) (*Entity, error)
	Update(item *Entity) error
	Cleanup(ttlMinutes int32)
	DeleteById(id string) error
}

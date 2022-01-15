package interfaces

type Repository interface {
	Insert(e Entity) (*Entity, error)
	Find(id string) (*Entity, error)
}

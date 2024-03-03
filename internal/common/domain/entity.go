package domain

type Entity struct {
	id uint64
}

func (e *Entity) GetId() uint64 {
	return e.id
}

func (e *Entity) Equal(other *Entity) bool {
	if e == other {
		return true
	}

	if other == nil {
		return false
	}

	return e.id == other.id
}

func NewEntity() *Entity {
	return &Entity{}
}

func NewEntityWithId(id uint64) *Entity {
	return &Entity{
		id: id,
	}
}

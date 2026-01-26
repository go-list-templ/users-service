package vo

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func NewID() ID {
	return ID{value: uuid.New()}
}

func UnsafeID(id uuid.UUID) ID {
	return ID{value: id}
}

func (i *ID) Value() uuid.UUID {
	return i.value
}

package vo

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func NewID() (ID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return ID{}, err
	}

	return ID{value: id}, nil
}

func UnsafeID(id uuid.UUID) ID {
	return ID{value: id}
}

func (i *ID) Value() uuid.UUID {
	return i.value
}

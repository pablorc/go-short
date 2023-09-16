package keygen

import "github.com/google/uuid"

func NewKey() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	return id.String()[:6]
}

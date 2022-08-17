package domain

import (
	"github.com/google/uuid"
)

type UserReq struct {
	Name 	string
}
type UserRes struct {
	Id		uuid.UUID
	Name	string
}
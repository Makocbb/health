package repositories

import "github.com/uptrace/bun"

type DBRepository interface {
	GetDB() *bun.DB
}

package advertisement

import (
	"marketplace/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	advertisementsTable = "advertisements"

	adIDColumn          = "id"
	adTitleColumn       = "title"
	adDescriptionColumn = "description"
	asImageUrlColumn    = "image_url"
	adPriceColumn       = "price"
	adAuthorIDColumn    = "author_id"
	adCreatedAtColumn   = "created_at"

	usersTable = "users"

	userIDColumn    = "id"
	userLoginColumn = "user_login"
)

type adRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.AdvertisementRepository {
	return &adRepository{db: db}
}

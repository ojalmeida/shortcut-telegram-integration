package crud

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrReviewNotFound = errors.New("user not found")
)

type Review struct {
	ID                    string `db:"id"`
	URL                   string `db:"url"`
	Number                int64  `db:"number"`
	RequesterID           string `db:"requester_id"`
	RequestedID           string `db:"requested_id"`
	SuppressedByRequester bool   `db:"suppressed_by_requester"`
	SuppressedByRequested bool   `db:"suppressed_by_requested"`
	NotifyAt              int64  `db:"notify_at"`
}

func CreateReview(r Review) (created Review, err error) {

	if r.ID == "" {

		var uniqueID uuid.UUID
		uniqueID, err = uuid.NewUUID()
		if err != nil {
			return
		}

		r.ID = uniqueID.String()

	}

	query := `INSERT INTO "reviews" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	err = db.Get(&created, query, r.ID, r.URL, r.Number, r.RequesterID, r.RequestedID, r.SuppressedByRequester, r.SuppressedByRequested, r.NotifyAt)

	return
}

func GetReview(id string) (retrieved Review, err error) {

	query := `SELECT * FROM "reviews" WHERE id = $1 LIMIT 1`
	err = db.Get(&retrieved, query, id)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			err = ErrReviewNotFound

		}

	}

	return
}

func GetReviews() (retrieved []Review, err error) {

	query := `SELECT * FROM "reviews"`
	err = db.Select(&retrieved, query)

	return
}

func UpdateReview(id string, r Review) (updated Review, err error) {

	query := `UPDATE "reviews" SET 
                	url = $2,
                	number = $3,
                   	requester_id = $4, 
                   	requested_id = $5,
                   	suppressed_by_requester = $6,
                   	suppressed_by_requested = $7,
                   	notify_at = $8
               WHERE id = $1 RETURNING *`

	err = db.Get(&updated, query, id, r.URL, r.Number, r.RequesterID, r.RequestedID, r.SuppressedByRequester, r.SuppressedByRequested, r.NotifyAt)

	return
}

func DeleteReview(id string) (err error) {

	query := `DELETE FROM "reviews" WHERE id = $1`
	_, err = db.Exec(query, id)

	return
}

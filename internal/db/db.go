// ./internal/db.go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/febriandani/backend-user-service/internal/infra"
	"github.com/febriandani/backend-user-service/protogen/golang/users"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DB struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

// NewDB creates a new array to mimic the behaviour of a in-memory database
func NewDB(db *infra.DatabaseList, logger *logrus.Logger) *DB {
	return &DB{
		db:  db,
		log: logger,
	}
}

// SaveUser adds a new order to the DB collection. Returns an error on duplicate ids
func (d *DB) SaveUser(ctx context.Context, tx *sql.Tx, user *users.User) (int64, error) {
	InsertUser := `INSERT INTO public.users
	(username, email, password, is_active, created_at, updated_at, created_by, updated_by)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?)	
	 returning user_id;
	`

	param := make([]interface{}, 0)

	param = append(param, user.Username)
	param = append(param, user.Email)
	param = append(param, user.Password)
	param = append(param, user.IsActive)
	param = append(param, time.Now().UTC())
	param = append(param, time.Now().UTC())
	param = append(param, user.CreatedBy)
	param = append(param, user.UpdatedBy)

	query, args, err := d.db.Backend.Write.In(InsertUser, param...)
	if err != nil {
		return 0, err
	}

	query = d.db.Backend.Write.Rebind(query)

	d.log.WithField("QueryDebug : ", query).Infof("Query SaveUser")

	var res *sql.Row
	if tx == nil {
		res = d.db.Backend.Write.QueryRow(ctx, query, args...)
	} else {
		res = tx.QueryRowContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	err = res.Err()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DB) CheckIsExistUser(ctx context.Context, user *users.User) (bool, error) {
	var res bool

	GetUser := `select exists(select 1 from public.users u where (u.username = '` + user.Username + `' or u.email = '` + user.Email + `') or (u.username = '` + user.Email + `' or u.email = '` + user.Email + `'));`

	query, args, err := d.db.Backend.Read.In(GetUser)
	if err != nil {
		return res, err
	}

	query = d.db.Backend.Read.Rebind(query)
	d.log.WithField("QueryDebug : ", query).Infof("Query isExists user")
	err = d.db.Backend.Write.Get(&res, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return res, err
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

// GetUserByID returns an order by the order_id
func (d *DB) GetUserByID(ctx context.Context, userID uint64) (*users.User, error) {
	var result users.User

	query := fmt.Sprintf(`SELECT user_id, username, email, is_active, created_at, created_by, updated_at, updated_by FROM public.users WHERE user_id = %d`, userID)

	rows, err := d.db.Backend.Write.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var createdAt, updatedAt time.Time
	err = rows.Scan(&result.UserId, &result.Username, &result.Email, &result.IsActive, &createdAt, &result.CreatedBy, &updatedAt, &result.UpdatedBy)
	if err != nil {
		return nil, err
	}

	// Convert time.Time to *timestamppb.Timestamp
	createdAtProto := timestamppb.New(createdAt)
	updatedAtProto := timestamppb.New(updatedAt)
	result.CreatedAt = createdAtProto
	result.UpdatedAt = updatedAtProto

	return &result, nil
}

func (d *DB) GetUserByEmailOrUsername(data string) (*users.User, error) {
	var result users.User

	query := fmt.Sprintf(`SELECT user_id, username, email, is_active, password FROM public.users WHERE username = '%s' or email = '%s';`, data, data)

	rows, err := d.db.Backend.Write.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(&result.UserId, &result.Username, &result.Email, &result.IsActive, &result.Password)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUserByIDs returns all users pertaining to the given order ids
func (d *DB) GetUserByIDs(userIDs []uint64) []*users.User {
	filtered := make([]*users.User, 0)

	// for _, idx := range userIDs {
	// 	for _, order := range d.collection {
	// 		if order.UserId == idx {
	// 			filtered = append(filtered, order)
	// 			break
	// 		}
	// 	}
	// }

	return filtered
}

// UpdateUser updates an order in place
func (d *DB) UpdateUser(order *users.User) {
	// for i, o := range d.collection {
	// 	if o.UserId == order.UserId {
	// 		d.collection[i] = order
	// 		return
	// 	}
	// }
}

// RemoveUser removes an order from the users collection
func (d *DB) RemoveUser(userID uint64) {
	// filtered := make([]*users.User, 0, len(d.collection)-1)
	// for i := range d.collection {
	// 	if d.collection[i].UserId != userID {
	// 		filtered = append(filtered, d.collection[i])
	// 	}
	// }
	// d.collection = filtered
}

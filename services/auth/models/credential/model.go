package credential

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/hash"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models"

	"github.com/lib/pq"
)

type Credential struct {
	Id                          int
	Uuid                        string
	Password                    sql.NullString
	CreatedAt                   pq.NullTime
	UpdatedAt                   pq.NullTime
	FailedAttempts              int
	LockedUntil                 pq.NullTime
	PasswordResetToken          sql.NullString
	PasswordResetTokenExpiresAt pq.NullTime
	Email                       sql.NullString
	Phone                       sql.NullString
}

func All(db *sql.DB) ([]*Credential, error) {
	rows, err := db.Query("select * from credentials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := []*Credential{}
	for rows.Next() {
		c, err := buildFromRow(rows)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, c)
	}

	return credentials, nil
}

func FindBy(db *sql.DB, fields map[string]interface{}) (*Credential, error) {
	var findStatement []string
	var findValues []interface{}

	index := 0
	for k, v := range fields {
		index++
		findStatement = append(findStatement, fmt.Sprintf("%v = $%v", k, index))
		findValues = append(findValues, v)
	}

	sql := "select * from credentials where " + strings.Join(findStatement, " and ")

	return buildFromRow(db.QueryRow(sql, findValues...))
}

func (credential *Credential) Create(db *sql.DB) error {
	hashValue, err := hash.HashPassword(credential.Password.String)
	if err != nil {
		return err
	}

	err = db.QueryRow(
		`insert into credentials
		 (email, phone, password, created_at, updated_at) values
		 ($1, $2, $3, $4, $5) returning id, uuid`,
		credential.Email, credential.Phone, hashValue, time.Now(), time.Now(),
	).Scan(&credential.Id, &credential.Uuid)

	return err
}

func (credential *Credential) Update(db *sql.DB, fields map[string]interface{}) error {
	var findStatement []string
	var findValues []interface{}

	index := 0
	for k, v := range fields {
		index++
		findStatement = append(findStatement, fmt.Sprintf("%v = $%v", k, index))
		findValues = append(findValues, v)
	}

	sql := fmt.Sprintf("update credentials set %s where id = %d", strings.Join(findStatement, ", "), credential.Id)

	_, err := db.Exec(sql, findValues...)

	return err
}

func (credential *Credential) IncrementFailedAttempt(db *sql.DB) error {
	_, err := db.Exec(`update credentials set failed_attempts = failed_attempts + 1
                       where id = $1 and failed_attempts = $2`, credential.Id, credential.FailedAttempts)

	return err
}

func (credential *Credential) Delete(db *sql.DB) error {
	_, err := db.Exec("delete from credentials where id = $1", credential.Id)

	return err
}

func (credential *Credential) UpdatePassword(db *sql.DB) error {
	hashValue, err := hash.HashPassword(credential.Password.String)
	if err != nil {
		return nil
	}

	_, err = db.Exec("update credentials set password = $1, password_reset_token = null where id = $2", hashValue, credential.Id)

	return err
}

func (credential *Credential) SetPasswordResetToken(db *sql.DB, token string) error {
	hashValue, err := hash.HashPassword(token)
	if err != nil {
		return nil
	}

	_, err = db.Exec("update credentials set password_reset_token = $1, where id = $2", hashValue, credential.Id)

	return err
}

func buildFromRow(row models.ScannableObject) (*Credential, error) {
	var credential Credential

	err := row.Scan(
		&credential.Id,
		&credential.Uuid,
		&credential.Password,
		&credential.CreatedAt,
		&credential.UpdatedAt,
		&credential.FailedAttempts,
		&credential.LockedUntil,
		&credential.PasswordResetToken,
		&credential.PasswordResetTokenExpiresAt,
		&credential.Email,
		&credential.Phone,
	)

	if err != nil {
		fmt.Println(err)
		return &credential, err
	}

	return &credential, nil
}

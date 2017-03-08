package models

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
)

const (
	PASSWORD_SALT_BYTES = 32
	PASSWORD_HASH_BYTES = 64
)

type User struct {
	Id              int64  `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password,omitempty"`
	IsAdministrator bool   `json:"isAdministrator"`
	salt            []byte
	hashedPassword  []byte
}

func NewUser(username string, password string, isAdministrator bool) *User {
	user := &User{}
	user.Username = username

	user.salt = mustGenerateSalt()
	user.hashedPassword = hashPassword(password, user.salt)
	user.IsAdministrator = isAdministrator

	return user
}

func (user *User) Type() string {
	return "user"
}

func (user *User) DatabaseTable() string {
	return "users"
}

func (user *User) RequireTransaction() bool {
	return false
}

func (user *User) IsCorrectPassword(password string) bool {
	hashedPassword := hashPassword(password, user.salt)
	return string(user.hashedPassword) == string(hashedPassword)
}

func (user *User) Save(execer SQLExecer) error {
	result, err := execer.Exec("INSERT INTO users(username, salt, hashed_password, is_administrator) VALUES(?,?,?,?)",
		user.Username,
		user.salt,
		user.hashedPassword,
		user.IsAdministrator,
	)

	if err != nil {
		return NewAPIError(500, "Failed to create user", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to retrieve last inserted id: %s", err)
	}

	user.Id = lastInsertedId
	return nil
}

func (user *User) Load(queryer SQLQueryer) error {
	panic("Not Implemented")
	return nil
}

func (user *User) Update(execer SQLExecer) error {
	panic("Not Implemented")
	return nil
}

func (user *User) Delete(execer SQLExecer) error {
	return nil
}

func (user *User) LoadFromUsername(username string, queryer SQLQueryer) error {
	err := queryer.QueryRow("SELECT id, username, salt, hashed_password, is_administrator FROM users WHERE username=?", username).
		Scan(&user.Id, &user.Username, &user.salt, &user.hashedPassword, &user.IsAdministrator)

	if err == sql.ErrNoRows {
		return NewAPIError(404, fmt.Sprintf("No user with id %d exist", user.Id), nil)
	} else if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to load user with id %d", user.Id), err)
	}

	return nil
}

func hashPassword(password string, salt []byte) []byte {
	// FIXME Use a much more secure hashing algoritm then md5.
	md5Hash := md5.New()
	io.WriteString(md5Hash, password)
	io.WriteString(md5Hash, string(salt))
	return md5Hash.Sum(nil)
}

func mustGenerateSalt() []byte {
	salt := make([]byte, PASSWORD_SALT_BYTES)

	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	return salt
}

// IMPORTANT Super ugly helper funtion to force user 1 to be default administrator when debugging.
func MustCreateDefaultAdministratorIfMissing(execer SQLExecer) {
	user := NewUser("admin", "admin", true)
	_, err := execer.Exec("INSERT OR REPLACE INTO users(id, username, salt, hashed_password, is_administrator) values (?,?,?,?,?)", 1, user.Username, user.salt, user.hashedPassword, user.IsAdministrator)
	if err != nil {
		panic(err)
	}
}

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	// Email   string `form:"email"`
	// Message string `form:"message" binding:"required"`
}

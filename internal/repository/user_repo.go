package repository

import (
	"database/sql"
	"log"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)
type UserRepository interface {
	Create(newUser *domain.User) error
	Update(updetedUser *domain.User) error
	Delete(id int) error
	FindOne(id int) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindAll() ([]domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}
func (u userRepository) Create(newUser *domain.User) error {
	query := `INSERT INTO users (user_id, name, email, password, profile_picture) VALUES ($1, $2, $3, $4, $5)`

	_, err := u.db.Exec(query, newUser.ID, newUser.Name, newUser.Email, newUser.Password, newUser.ProfilePicture)
	if err != nil {
		log.Println("Failed to add user:", err)
		return err
	} else {
		log.Println("Succsessfully added user")
	}
	return err
}

func (u userRepository) Update(updetedUser *domain.User) error {
	query := `UPDATE users SET name=$2, email=$3, password=$4, profile_picture=$5 WHERE user_id=$1`
	_, err := u.db.Exec(query, updetedUser.ID, updetedUser.Name, updetedUser.Email, updetedUser.Password, updetedUser.ProfilePicture)
	if err != nil {
		log.Println("Failed to update user:", err)
		return err
	} else {
		log.Println("Succsessfully updated")
	}

	return err
}

func (u *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := u.db.Exec(query, id)
	if err != nil {
		log.Println("Failed to delete user:", err)
		return err
	}
	log.Println("Successfully deleted user")
	return nil
}

func (u userRepository) FindOne(id int) (*domain.User, error) {
	query := `SELECT name, email, password, profile_picture, is_deleted FROM users WHERE user_id=$1;`
	row := u.db.QueryRow(query, id)
	var user domain.User
	err := row.Scan(&user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
	if err != nil {
		panic(err)
	}
	user.ID = id
	return &user, nil
}

func (u userRepository) FindAll() ([]domain.User, error) {
	query := `SELECT user_id, name, email, password, profile_picture, is_deleted FROM users;`
	rows, err := u.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}

func (u userRepository) FindByUsername(username string) (*domain.User, error) {
	query := `SELECT name, email, password, profile_picture, is_deleted FROM users WHERE name=$1;`
	row := u.db.QueryRow(query, username)
	var user domain.User
	err := row.Scan(&user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
	if err != nil {
		panic(err)
	}
	return &user, nil
}

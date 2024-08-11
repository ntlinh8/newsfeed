package repo

import (
	"crypto/sha256"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	db_model "newsfeed/pkg/repo/model"
	"newsfeed/pkg/service/model"
)

type PublicUserRepo struct {
	db *gorm.DB
}

type MySQLConfig struct {
	Username     string
	Password     string
	Addr         string
	DatabaseName string
}

func NewUserRepo(conf MySQLConfig) (*PublicUserRepo, error) {
	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Addr,
		conf.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dbConnStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %s", err)
	}
	return &PublicUserRepo{
		db: db,
	}, nil
}

func (urp *PublicUserRepo) CreateUser(user *model.User) (*model.User, error) {
	log.Println("[DEBUG] input user", user)

	h := sha256.New()
	h.Write([]byte(user.Password))
	bytes := h.Sum(nil)
	log.Println("[DEBUG] raw byte hashed", bytes)
	hashPassword := fmt.Sprintf("%x", bytes)

	log.Println("[DEBUG] hashed", hashPassword)

	dbUser := &db_model.DbUser{
		Username:     user.Username,
		HashPassword: hashPassword,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DOB:          user.DOB,
		Email:        user.Email,
	}

	res := urp.db.Create(dbUser)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to create user in db: %s", res.Error)
	}
	log.Println("[DEBUG] create sql res: +%v", res)
	log.Println("[DEBUG] created user: +%v", dbUser)

	user.UserId = dbUser.Id
	return user, nil
}

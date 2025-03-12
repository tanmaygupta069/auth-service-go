package mysql

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tanmaygupta069/auth-service/config"
)

var db *gorm.DB

var once sync.Once

type SqlInterface interface {
	Insert(user *User) error
	Delete(email string) error
	GetOne(key string,value string) (*User, error)
	GetAll() ([]User, error)
}

type SqlServiceImplementation struct {
	db *gorm.DB
}

func NewSqlClient() *SqlServiceImplementation {
	return &SqlServiceImplementation{
		GetSqlClient(),
	}
}

func InitializeSqlClient() {
	once.Do(func() {
		cfg, er := config.GetConfig()
		if er != nil {
			fmt.Println("error occured in sql client init")
		}
		dbConfig := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.MySqlConfig.User, cfg.MySqlConfig.Password, cfg.MySqlConfig.Host, cfg.MySqlConfig.Port, cfg.MySqlConfig.Database)
		d, err := gorm.Open(mysql.Open(dbConfig), &gorm.Config{})
		if err != nil {
			fmt.Println("error occured while connecting to mysql")
		}
		sqlDB, err := d.DB()
		if err != nil {
			fmt.Println("Failed to get DB instance:", err)
			return
		}
		if err := sqlDB.Ping(); err != nil {
			fmt.Println("Failed to ping DB:", err)
			return
		}

		if err := d.AutoMigrate(&User{}); err != nil {
			fmt.Println("Failed to auto-migrate:", err)
			return
		}

		fmt.Println("DB connection successful")
		db = d
	})
}

func GetSqlClient() *gorm.DB {
	if db == nil {
		InitializeSqlClient()
	}
	return db
}

func (s *SqlServiceImplementation) GetOne(key string,value string) (*User, error) {
	var user User
	err := db.Where(fmt.Sprintf("%s = ?", value),key).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SqlServiceImplementation) GetAll() ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *SqlServiceImplementation) Insert(user *User) error {
	return db.Create(user).Error
}

func (s *SqlServiceImplementation) Delete(email string) error {
	return db.Delete(&User{}, email).Error
}

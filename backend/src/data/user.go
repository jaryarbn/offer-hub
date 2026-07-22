package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const userInfoTable = "t_user_info"

var userInfoColumns = []string{
	"user_id",
	"username",
	"nick_name",
	"avatar",
	"introduction",
	"sex",
	"vip",
	"phone",
	"email",
	"user_status",
	"user_type",
	"create_time",
	"update_time",
}

var ErrUsernameAlreadyExists = errors.New("username already exists")

var ErrUserNotFound = errors.New("user not found")

// CreateUserRecord contains the fields needed when creating a user.
// Remaining t_user_info columns use their database defaults.
type CreateUserRecord struct {
	UserID   string `gorm:"column:user_id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	NickName string `gorm:"column:nick_name"`
}

func (CreateUserRecord) TableName() string {
	return userInfoTable
}

// UserRecord is the authentication view of a t_user_info row.
type UserRecord struct {
	UserID     string `gorm:"column:user_id"`
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	NickName   string `gorm:"column:nick_name"`
	Avatar     string `gorm:"column:avatar"`
	Sex        int    `gorm:"column:sex"`
	VIP        bool   `gorm:"column:vip"`
	Phone      string `gorm:"column:phone"`
	Email      string `gorm:"column:email"`
	UserStatus int    `gorm:"column:user_status"`
	UserType   int    `gorm:"column:user_type"`
}

func (UserRecord) TableName() string {
	return userInfoTable
}

// UserInfoRecord excludes the password hash from profile queries.
type UserInfoRecord struct {
	UserID       string    `gorm:"column:user_id"`
	Username     string    `gorm:"column:username"`
	NickName     string    `gorm:"column:nick_name"`
	Avatar       string    `gorm:"column:avatar"`
	Introduction string    `gorm:"column:introduction"`
	Sex          int       `gorm:"column:sex"`
	VIP          bool      `gorm:"column:vip"`
	Phone        string    `gorm:"column:phone"`
	Email        string    `gorm:"column:email"`
	UserStatus   int       `gorm:"column:user_status"`
	UserType     int       `gorm:"column:user_type"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func (UserInfoRecord) TableName() string {
	return userInfoTable
}

func (data *Data) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := data.MySQL.WithContext(ctx).
		Model(&UserRecord{}).
		Where("username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return false, fmt.Errorf("query username existence: %w", err)
	}
	return count > 0, nil
}

func (data *Data) CreateUser(ctx context.Context, record CreateUserRecord) error {
	err := data.MySQL.WithContext(ctx).Create(&record).Error
	if isDuplicateKeyError(err) {
		return ErrUsernameAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (data *Data) GetUserByUsername(ctx context.Context, username string) (UserRecord, error) {
	var record UserRecord
	err := data.MySQL.WithContext(ctx).
		Where("username = ?", username).
		Take(&record).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserRecord{}, ErrUserNotFound
	}
	if err != nil {
		return UserRecord{}, fmt.Errorf("query user by username: %w", err)
	}
	return record, nil
}

func (data *Data) GetUserByID(ctx context.Context, userID string) (UserInfoRecord, error) {
	var record UserInfoRecord
	err := data.MySQL.WithContext(ctx).
		Select(userInfoColumns).
		Where("user_id = ?", userID).
		Take(&record).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserInfoRecord{}, ErrUserNotFound
	}
	if err != nil {
		return UserInfoRecord{}, fmt.Errorf("query user by user_id: %w", err)
	}
	return record, nil
}

func isDuplicateKeyError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

package models

// User model
type User struct {
	Model
	Email          string          `json:"email"           gorm:"type:varchar(512); not null;"`
	EmailVerified  bool            `json:"email_verified"  gorm:"type:tinyint(1); not null;"`
	Name           string          `json:"name"            gorm:"type:varchar(512); not null;"`
	LastName       string          `json:"last_name"       gorm:"type:varchar(512); not null"`
	Picture        string          `json:"picture"         gorm:"type:varchar(2048); not null;"`
	AccessControls []AccessControl `json:"access_controls" gorm:"foreignkey:UserId;"`
}

// Indices implementation of the Modelable interface
func (u *User) Indices() map[string][]string {
	return map[string][]string{
		"unique_account_email_deleted": []string{"account_id", "email", "deleted_at"},
	}
}

// Table implementation of the Modelable interface
func (u *User) Table() string {
	return "users"
}

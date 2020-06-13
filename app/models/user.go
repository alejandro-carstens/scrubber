package models

// User model
type User struct {
	Model
	Email         string `json:"email"           gorm:"type:varchar(512); not null;"`
	EmailVerified bool   `json:"email_verified"  gorm:"type:tinyint(1); not null;"`
	FullName      string `json:"full_name"       gorm:"type:varchar(512); not null;"`
	Name          string `json:"name"            gorm:"type:varchar(512); not null;"`
	LastName      string `json:"last_name"       gorm:"type:varchar(512); not null"`
	Picture       string `json:"picture"         gorm:"type:varchar(2048); not null;"`
	Roles         []Role `json:"roles,omitempty" gorm:"many2many:users_roles;association_foreignkey:user_id;association_foreignkey:role_id;"`
}

// Indices implementation of the Modelable interface
func (u *User) Indices() map[string][]string {
	return map[string][]string{
		"unique_email":      []string{"email"},
		"full_name_deleted": []string{"full_name", "deleted_at"},
	}
}

// Table implementation of the Modelable interface
func (u *User) Table() string {
	return "users"
}

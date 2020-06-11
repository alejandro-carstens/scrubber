package models

type UserRole struct {
	Model
	UserID uint64 `json:"user_id" gorm:"type:bigint unsigned;not null;"`
	RoleID uint64 `json:"role_id" gorm:"type:bigint unsigned;not null;"`
}

// Indices implementation of the Modelable interface
func (ur *UserRole) Indices() map[string][]string {
	return map[string][]string{
		"user_deleted": []string{"user_id", "deleted_at"},
		"role_deleted": []string{"role_id", "deleted_at"},
	}
}

// Table implementation of the Modelable interface
func (ur *UserRole) Table() string {
	return "users_roles"
}

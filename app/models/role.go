package models

type Role struct {
	Model
	Name        string       `json:"name"  gorm:"type:varchar(512); not null"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:users_roles;association_foreignkey:user_id;association_foreignkey:role_id;"`
	Permissions []Permission `json:"permissions"`
}

// Indices implementation of the Modelable interface
func (u *Role) Indices() map[string][]string {
	return map[string][]string{
		"unique_name": []string{"name"},
	}
}

// Table implementation of the Modelable interface
func (u *Role) Table() string {
	return "roles"
}

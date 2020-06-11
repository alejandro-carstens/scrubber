package models

type Permission struct {
	Model
	RoleID uint64 `json:"user_id" gorm:"type:bigint unsigned;not null;"`
	Action string `json:"action"  gorm:"varchar(255);not null;"`
	Scope  string `json:"scope"   gorm:"varchar(255);not null;"`
}

// Indices implementation of the Modelable interface
func (p *Permission) Indices() map[string][]string {
	return map[string][]string{
		"role_deleted": []string{"role_id", "deleted_at"},
	}
}

// Table implementation of the Modelable interface
func (p *Permission) Table() string {
	return "permissions"
}

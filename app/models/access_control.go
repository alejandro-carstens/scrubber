package models

type AccessControl struct {
	Model
	UserID uint64 `json:"user_id" gorm:"type:bigint unsigned;not null;"`
	Action string `json:"action"  gorm:"varchar(255);not null;"`
	Scope  string `json:"scope"   gorm:"varchar(255);not null;"`
}

// Indices implementation of the Modelable interface
func (ac *AccessControl) Indices() map[string][]string {
	return map[string][]string{
		"user_action": []string{"user_id", "action", "deleted_at"},
	}
}

// Table implementation of the Modelable interface
func (ac *AccessControl) Table() string {
	return "access_controls"
}

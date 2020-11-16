package domain

type RolesUser struct {
	ID  uint `gorm:"primarykey"`
	Rid int
	Uid int
}

func (RolesUser) TableName() string {
	return "roles_user"
}

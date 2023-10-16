package chat

import (
	"fmt"
)

var (
	masterMap = map[uint64]map[uint64]bool{}

	defaultPermission PermissionData = 0
)

func loadMasters(config map[uint64][]uint64) {
	for k, v := range config {
		masterMap[k] = map[uint64]bool{}
		for _, vv := range v {
			masterMap[k][vv] = true
		}
	}
}

type Permission struct {
	UserID      string `gorm:"column:user_id;primaryKey;type:varchar(20)"`
	UserName    string `gorm:"column:user_name;type:varchar(255)"`
	Permissions uint64 `gorm:"column:permissions;index"`
}

func (p Permission) TableName() string {
	return "permissions"
}

type PermissionData uint64

const (
	permissionNone        PermissionData = 0
	permissionChat        PermissionData = 1 << 0
	permissionChatPlus    PermissionData = 1 << 1
	permissionImageSmall  PermissionData = 1 << 2
	permissionImageMedium PermissionData = 1 << 3
	permissionImageLarge  PermissionData = 1 << 4
	permissionImageMulti  PermissionData = 1 << 5

	PermissionNone        = "none"
	PermissionChat        = "chat"
	PermissionChatPlus    = "chat+"
	PermissionImageSmall  = "image-small"
	PermissionImageMedium = "image-medium"
	PermissionImageLarge  = "image-large"
	PermissionImageMulti  = "image-multi"
)

var (
	permissionMap = map[string]PermissionData{
		PermissionNone:        permissionNone,
		PermissionChat:        permissionChat,
		PermissionChatPlus:    permissionChatPlus,
		PermissionImageSmall:  permissionImageSmall,
		PermissionImageMedium: permissionImageMedium,
		PermissionImageLarge:  permissionImageLarge,
		PermissionImageMulti:  permissionImageMulti,
	}
)

func (p PermissionData) HasPermission(permission PermissionData) bool {
	return p&permission != 0
}

func (p PermissionData) AddPermission(permission PermissionData) PermissionData {
	return p | permission
}

func (p PermissionData) RemovePermission(permission PermissionData) PermissionData {
	return p &^ permission
}

func (p PermissionData) AllPermissions() []string {
	var result []string
	for k, v := range permissionMap {
		if p.HasPermission(v) {
			result = append(result, k)
		}
	}

	return result
}

func parsePermission(permission PermissionData) uint64 {
	return uint64(permission)
}

func format(permission string) (p PermissionData, exist bool) {
	if p, exist = permissionMap[permission]; exist {
		return p, true
	}

	return permissionNone, false
}

func addPermission(userID string, nickname string, permissions ...string) error {
	var fp = defaultPermission
	for _, ps := range permissions {
		if rp, exist := format(ps); exist {
			fp = fp.AddPermission(rp)
		} else {
			return fmt.Errorf("permission %s not exist", ps)
		}
	}

	var p Permission
	has, tryGetErr := database.Has(Permission{}.TableName(), "user_id = ?", userID)

	if tryGetErr != nil {
		return fmt.Errorf("cannot get permission: %w", tryGetErr)
	} else if !has {
		p = Permission{
			UserID:      userID,
			UserName:    nickname,
			Permissions: parsePermission(defaultPermission),
		}
	} else if getErr := database.GetOne(&p, "user_id = ?", userID); getErr != nil {
		return fmt.Errorf("cannot get permission: %w", getErr)
	}

	op := PermissionData(p.Permissions).AddPermission(fp)

	if has {
		return database.UpdateOne(&Permission{
			Permissions: parsePermission(op),
		}, "user_id = ?", userID)
	} else {
		return database.InsertOne(&Permission{
			UserID:      userID,
			UserName:    nickname,
			Permissions: parsePermission(op),
		})
	}
}

func removePermission(userID string, permissions ...string) error {
	var p Permission
	if has, tryGetErr := database.Has(Permission{}.TableName(), "user_id = ?", userID); tryGetErr != nil {
		return fmt.Errorf("cannot get permission: %w", tryGetErr)
	} else if !has {
		return fmt.Errorf("user %s not exist", userID)
	} else if getErr := database.GetOne(&p, "user_id = ?", userID); getErr != nil {
		return fmt.Errorf("cannot get permission: %w", getErr)
	}

	op := PermissionData(p.Permissions)
	for _, ps := range permissions {
		if rp, exist := format(ps); exist {
			op = op.RemovePermission(rp)
		}
	}

	return database.UpdateOne(&Permission{
		Permissions: parsePermission(op),
	}, "user_id = ?", userID)
}

func checkPermission(userID string, permission string) (legal bool, err error) {
	var p Permission
	needCheck, exist := format(permission)
	if !exist {
		return false, fmt.Errorf("permission %s not exist", permission)
	}

	if defaultPermission.HasPermission(needCheck) {
		return true, nil
	}

	if has, tryGetErr := database.Has(Permission{}.TableName(), "user_id = ?", userID); tryGetErr != nil {
		return false, fmt.Errorf("cannot get permission: %w", tryGetErr)
	} else if !has {
		return false, nil
	} else if getErr := database.GetOne(&p, "user_id = ?", userID); getErr != nil {
		return false, fmt.Errorf("cannot get permission: %w", getErr)
	} else {
		return PermissionData(p.Permissions).HasPermission(needCheck), nil
	}
}

func getPermission(userID string) (permission Permission, err error) {
	if has, tryGetErr := database.Has(Permission{}.TableName(), "user_id = ?", userID); tryGetErr != nil {
		return Permission{}, fmt.Errorf("cannot get permission: %w", tryGetErr)
	} else if !has {
		return Permission{}, nil
	} else if getErr := database.GetOne(&permission, "user_id = ?", userID); getErr != nil {
		return Permission{}, fmt.Errorf("cannot get permission: %w", getErr)
	} else {
		return permission, nil
	}
}

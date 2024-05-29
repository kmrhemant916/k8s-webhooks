package global

const (
	RoleFindQueryByName = "SELECT * FROM `roles` WHERE name = ?"
	PermissionFindQueryByName = "SELECT * FROM `permissions` WHERE name = ?"
	RolePermissionFindQueryByID = "SELECT * FROM `role_permissions` WHERE role_id = ? AND permission_id = ?"
	RolePermissionFindQueryByWildcardRoleID = "SELECT * FROM `role_permissions` WHERE role_id IN (?) AND permission_id = ?"
	GroupFindQueryByName = "SELECT * FROM `groups` WHERE name = ?"
	GroupFindQueryByID = "SELECT * FROM `groups` WHERE group_id = ?"
	GroupRoleFindQueryByGroupRoleID = "SELECT * FROM `group_roles` WHERE role_id = ? AND group_id = ?"
	GroupRoleFindQueryByRoleID = "SELECT * FROM `group_roles` WHERE role_id = ?"
	UserFindQueryByID = "SELECT * FROM `users` WHERE user_id = ?"
	UserGroupFindQueryByID = "SELECT * FROM `user_groups` WHERE user_id = ? and group_id = ?"
	UserFindQueryByEmail = "SELECT * FROM `users` WHERE email = ?"
	UserGroupFindQueryByUserID = "SELECT * FROM `user_groups` WHERE user_id = ?"
	GroupRoleFindQueryByID = "SELECT * FROM `group_roles` WHERE group_id = ?"
	RoleFindQueryByID = "SELECT * FROM `roles` WHERE id = ?"
	RolePermissionFindQueryByRoleID = "SELECT * FROM `role_permissions` WHERE role_id = ?"
	PermissionFindQueryByID = "SELECT * FROM `permissions` WHERE id = ?"
)
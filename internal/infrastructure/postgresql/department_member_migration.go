package infrastructure

import "gorm.io/gorm"

// migrateDepartmentMemberRoles backfills department_members.role_id from users / user_departments before columns are dropped.
func migrateDepartmentMemberRoles(db *gorm.DB) error {
	var deptMembersTable int64
	if err := db.Raw(`
		SELECT count(*) FROM information_schema.tables
		WHERE table_schema = current_schema() AND table_name = 'department_members'
	`).Scan(&deptMembersTable).Error; err != nil {
		return err
	}
	if deptMembersTable == 0 {
		return nil
	}
	if err := db.Exec(`
		ALTER TABLE department_members
		ADD COLUMN IF NOT EXISTS role_id bigint
	`).Error; err != nil {
		return err
	}
	var usersHaveRoleID int64
	if err := db.Raw(`
		SELECT count(*) FROM information_schema.columns
		WHERE table_schema = current_schema()
		  AND table_name = 'users'
		  AND column_name = 'role_id'
	`).Scan(&usersHaveRoleID).Error; err != nil {
		return err
	}
	if usersHaveRoleID > 0 {
		if err := db.Exec(`
			UPDATE department_members dm
			SET role_id = u.role_id
			FROM users u
			WHERE dm.deleted_at IS NULL
			  AND u.deleted_at IS NULL
			  AND dm.user_id = u.id
			  AND dm.role_id IS NULL
			  AND u.role_id IS NOT NULL
		`).Error; err != nil {
			return err
		}
	}
	var usersHaveDeptRole int64
	if err := db.Raw(`
		SELECT count(*) FROM information_schema.columns
		WHERE table_schema = current_schema()
		  AND table_name = 'users'
		  AND column_name IN ('department_id', 'role_id')
	`).Scan(&usersHaveDeptRole).Error; err != nil {
		return err
	}
	if usersHaveDeptRole == 2 {
		if err := db.Exec(`
			INSERT INTO department_members (created_at, updated_at, is_head, is_manager, department_id, user_id, role_id)
			SELECT NOW(), NOW(), false, false, u.department_id, u.id, u.role_id
			FROM users u
			WHERE u.deleted_at IS NULL
			  AND u.department_id IS NOT NULL
			  AND u.department_id > 0
			  AND u.role_id IS NOT NULL
			  AND u.role_id > 0
			  AND NOT EXISTS (
			    SELECT 1 FROM department_members dm
			    WHERE dm.deleted_at IS NULL
			      AND dm.user_id = u.id
			      AND dm.department_id = u.department_id
			  )
		`).Error; err != nil {
			return err
		}
	}
	var userDeptTable int64
	if err := db.Raw(`
		SELECT count(*) FROM information_schema.tables
		WHERE table_schema = current_schema() AND table_name = 'user_departments'
	`).Scan(&userDeptTable).Error; err != nil {
		return err
	}
	if userDeptTable > 0 && usersHaveDeptRole == 2 {
		if err := db.Exec(`
			INSERT INTO department_members (created_at, updated_at, is_head, is_manager, department_id, user_id, role_id)
			SELECT NOW(), NOW(), false, false, ud.department_id, ud.user_id, u.role_id
			FROM user_departments ud
			INNER JOIN users u ON u.id = ud.user_id AND u.deleted_at IS NULL
			WHERE ud.deleted_at IS NULL
			  AND u.role_id IS NOT NULL
			  AND u.role_id > 0
			  AND NOT EXISTS (
			    SELECT 1 FROM department_members dm
			    WHERE dm.deleted_at IS NULL
			      AND dm.user_id = ud.user_id
			      AND dm.department_id = ud.department_id
			  )
		`).Error; err != nil {
			return err
		}
	}
	return nil
}

func dropLegacyUserDepartmentColumns(db *gorm.DB) error {
	steps := []string{
		`ALTER TABLE users DROP COLUMN IF EXISTS role_id`,
		`ALTER TABLE users DROP COLUMN IF EXISTS department_id`,
		`DROP TABLE IF EXISTS user_departments CASCADE`,
	}
	for _, sql := range steps {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

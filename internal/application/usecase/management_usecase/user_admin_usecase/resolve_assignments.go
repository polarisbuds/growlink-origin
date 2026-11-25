package user_admin_usecase

import (
	"context"
	"errors"
	"fmt"

	outputmodel "growlink/internal/application/model/output_model"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type deptAssignmentResolved struct {
	DepartmentID uint
	RoleID       uint
}

type teamAssignmentResolved struct {
	TeamID uint
}

func validUUID(uuids ...string) error {
	var errs []error
	for _, id := range uuids {
		if _, err := uuid.Parse(id); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (uc *UserAdmin) resolveDepartmentAssignments(rows []DepartmentAssignmentInput) ([]deptAssignmentResolved, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("所属部署を選択してください")
	}

	departmentUUIDs := make([]string, len(rows))
	roleUUIDs := make([]string, len(rows))
	deptSeen := make(map[string]any, len(rows))
	roleSeen := make(map[string]any, len(rows))
	for idx, row := range rows {
		if err := validUUID(row.DepartmentUuid, row.RoleUuid); err != nil {
			return nil, fmt.Errorf("無効なUUID文字列が指定されています: %w", err)
		}
		departmentUUIDs[idx] = row.DepartmentUuid
		roleUUIDs[idx] = row.RoleUuid

		deptSeen[row.DepartmentUuid] = struct{}{}
		roleSeen[row.RoleUuid] = struct{}{}
	}

	group, _ := errgroup.WithContext(context.Background())
	deptIDMapByUUID := make(map[string]uint, len(rows))
	group.Go(func() error {
		departments, err := uc.deptRepo.GetByUUIDs(departmentUUIDs)
		if err != nil {
			return fmt.Errorf("部署情報の取得に失敗しました: %w", err)
		}
		for _, dept := range departments {
			if _, ok := deptSeen[dept.UUID.String()]; ok {
				deptIDMapByUUID[dept.UUID.String()] = dept.ID
			}
		}
		if len(deptIDMapByUUID) != len(rows) {
			return fmt.Errorf("存在しない部署が指定されました")
		}
		return nil
	})

	roleIDMapByUUID := make(map[string]uint, len(rows))
	group.Go(func() error {
		roles, err := uc.roleRepo.GetByUUIDs(roleUUIDs)
		if err != nil {
			return fmt.Errorf("役職情報の取得に失敗しました: %w", err)
		}
		for _, role := range roles {
			if _, ok := roleSeen[role.UUID.String()]; ok {
				roleIDMapByUUID[role.UUID.String()] = role.ID
			}
		}
		if len(roleIDMapByUUID) != len(rows) {
			return fmt.Errorf("存在しない役職が指定されました")
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	out := make([]deptAssignmentResolved, len(rows))
	for idx, row := range rows {
		out[idx] = deptAssignmentResolved{
			DepartmentID: deptIDMapByUUID[row.DepartmentUuid],
			RoleID:       roleIDMapByUUID[row.RoleUuid],
		}
	}
	return out, nil
}

func (uc *UserAdmin) resolveTeamAssignments(rows []TeamAssignmentInput) ([]teamAssignmentResolved, error) {
	teamUUIDs := make([]string, len(rows))
	for idx, row := range rows {
		teamUUIDs[idx] = row.TeamUuid
	}
	teams, err := uc.teamRepo.GetByUUIDs(teamUUIDs)
	if err != nil {
		return nil, fmt.Errorf("チームを選択してください")
	}

	out := make([]teamAssignmentResolved, len(teams))
	for idx, team := range teams {
		out[idx] = teamAssignmentResolved{
			TeamID: team.ID,
		}
	}
	return out, nil
}

func (uc *UserAdmin) departmentAssignmentsOutput(rows []deptAssignmentResolved) ([]outputmodel.UserDepartmentAssignmentOutput, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("所属部署を選択してください")
	}

	deptIDToUUID := map[uint]string{}
	group, _ := errgroup.WithContext(context.Background())
	group.Go(func() error {
		depts, err := uc.deptRepo.GetAll()
		if err != nil {
			return err
		}
		for _, dept := range depts {
			deptIDToUUID[dept.ID] = dept.UUID.String()
		}
		return nil
	})

	roleIDToUUID := map[uint]string{}
	group.Go(func() error {
		roles, err := uc.roleRepo.GetAll()
		if err != nil {
			return err
		}
		for _, role := range roles {
			roleIDToUUID[role.ID] = role.UUID.String()
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	out := make([]outputmodel.UserDepartmentAssignmentOutput, len(rows))
	for idx, row := range rows {
		out[idx] = outputmodel.UserDepartmentAssignmentOutput{
			DepartmentUUID: deptIDToUUID[row.DepartmentID],
			RoleUUID:       roleIDToUUID[row.RoleID],
		}
	}
	return out, nil
}

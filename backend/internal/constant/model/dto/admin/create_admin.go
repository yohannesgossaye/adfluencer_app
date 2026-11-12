package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
)

const (
	AdminRoleSuperAdmin = "super_admin"
	AdminApikey         = "cdghsbjncnjdnscjncs712637617637713289328981bsdcnjnjsncn"
)

type CreateAdminRequest struct {
	AdminName  string `json:"admin_name" validate:"required"`
	AdminEmail string `json:"admin_email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	AdminRole  string `json:"admin_role" validate:"required"`
}

func (r *CreateAdminRequest) Validate() error {
	if r.AdminName == "" {
		return customErr.ErrAdminNameRequired
	}
	if r.AdminEmail == "" {
		return customErr.ErrAdminEmailRequired
	}
	if r.Password == "" {
		return customErr.ErrPasswordRequired
	}
	if r.AdminRole == "" {
		return customErr.ErrAdminRoleRequired
	}
	if r.AdminRole != AdminRoleSuperAdmin {
		return customErr.ErrInvalidAdminRole
	}
	return nil
}

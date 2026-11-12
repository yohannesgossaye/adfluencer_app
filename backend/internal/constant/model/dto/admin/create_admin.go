package admin

import (
	"strings"

	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
)

// AdminRoleSuperAdmin is the only admin level for now.
const AdminRoleSuperAdmin = "super_admin"

// CreateAdminRequest is the body of POST /admin/create_admin.
type CreateAdminRequest struct {
	AdminName  string `json:"admin_name"`
	AdminEmail string `json:"admin_email"`
	Password   string `json:"password"`
	AdminRole  string `json:"admin_role"`
}

// Normalize trims the name and lower-cases the email exactly as login does,
// so the new admin can always log in with the address they typed.
func (r *CreateAdminRequest) Normalize() {
	r.AdminName = strings.TrimSpace(r.AdminName)
	r.AdminEmail = strings.ToLower(strings.TrimSpace(r.AdminEmail))
	r.AdminRole = strings.TrimSpace(r.AdminRole)
}

// Validate applies the same rules as sign-up. Call Normalize first.
func (r *CreateAdminRequest) Validate() error {
	switch {
	case r.AdminName == "":
		return customErr.ErrAdminNameRequired
	case len(r.AdminName) < 2 || len(r.AdminName) > 50:
		return customErr.ErrInvalidNameLength
	case r.AdminEmail == "":
		return customErr.ErrAdminEmailRequired
	case !authDto.ValidEmail(r.AdminEmail):
		return customErr.ErrInvalidEmail
	case r.Password == "":
		return customErr.ErrPasswordIsRequired
	case len(r.Password) < 8 || len(r.Password) > 20:
		return customErr.ErrInvalidPasswordLength
	case r.AdminRole == "":
		return customErr.ErrAdminRoleRequired
	case r.AdminRole != AdminRoleSuperAdmin:
		return customErr.ErrInvalidAdminRole
	}
	return nil
}

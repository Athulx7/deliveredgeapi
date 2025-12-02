package controllers

import (
	"DeliverEdgeapi/app/jwtconnections"
	"DeliverEdgeapi/app/services"
	"database/sql"

	// "golang.org/x/crypto/bcrypt"
	"github.com/revel/revel"
)

type AuthController struct {
	*revel.Controller
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCompany struct {
	UserID       int
	Email        string
	PasswordHash string
	CompanyID    int
	UserCode     string
	UserName     string
	RoleID       int
	RoleName     string
	CompanyName  string
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
}

func (c AuthController) Login() revel.Result {
	var req LoginRequest

	if err := c.Params.BindJSON(&req); err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid request body"})
	}

	query := `
	SELECT u.user_id, u.email, u.password_hash, u.company_id, u.user_code, u.username, u.role_id, r.role_name,
    c.company_name, c.db_host, c.db_user, c.db_password, c.db_name FROM tbl_global_users u 
	JOIN tbl_company_master c ON u.company_id = c.company_id 
	LEFT JOIN global_roles r ON u.role_id = r.role_id 
	WHERE u.email = ? AND u.status = 1 AND c.active_status = 1
	`

	var u UserCompany
	err := services.AdminDB.QueryRow(query, req.Email).Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.CompanyID,
		&u.UserCode, &u.UserName, &u.RoleID,
		&u.RoleName,
		&u.CompanyName, &u.DBHost, &u.DBUser,
		&u.DBPassword, &u.DBName,
	)
	if err == sql.ErrNoRows {
		return c.RenderJSON(map[string]string{"error": "Invalid email or password"})
	} else if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	// if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
	// 	return c.RenderJSON(map[string]string{"error": "Invalid password"})
	// }

	if services.TenantDBs == nil {
		services.TenantDBs = make(map[string]*sql.DB)
	}

	tenantDB, exists := services.TenantDBs[u.DBName]
	if !exists {
		tenantDB, err = services.ConnectTenantDB(u.DBHost, u.DBUser, u.DBPassword, u.DBName)
		if err != nil {
			return c.RenderJSON(map[string]string{"error": "Cannot connect tenant DB: " + err.Error()})
		}
		services.TenantDBs[u.DBName] = tenantDB
	}

	token, _ := jwtconnections.GenerateJWT(
		u.UserID, u.CompanyID, u.CompanyName, u.DBName,
		u.UserCode, u.UserName, u.RoleID, u.RoleName,
	)

	return c.RenderJSON(map[string]interface{}{
		"message":       "Login successful",
		"token":         token,
		"company_name":  u.CompanyName,
		"database_name": u.DBName,
	})
}

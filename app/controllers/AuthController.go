package controllers

import (
	"database/sql"
	"DeliverEdgeapi/app/services"
	"DeliverEdgeapi/app/jwtconnections"
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
	SELECT u.user_id, u.email, u.password_hash, u.company_id, u.user_Code,
	       c.company_name, c.db_host, c.db_user, c.db_password, c.db_name
	FROM tbl_global_users u
	JOIN tbl_company_master c ON u.company_id = c.company_id
	WHERE u.email = ? AND u.status = 1 AND c.active_status = 1;
`

	var u UserCompany
	err := services.AdminDB.QueryRow(query, req.Email).Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.CompanyID, &u.UserCode,
		&u.CompanyName, &u.DBHost, &u.DBUser, &u.DBPassword, &u.DBName,
	)
	if err == sql.ErrNoRows {
		return c.RenderJSON(map[string]string{"error": "Invalid email or password"})
	} else if err != nil {
		return c.RenderJSON(map[string]string{"error in the else": err.Error()})
	}

	// if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
	// 	return c.RenderJSON(map[string]string{"error": "Invalid password"})
	// }

	tenantDB, err := services.ConnectTenantDB(u.DBHost, u.DBUser, u.DBPassword, u.DBName)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Cannot connect tenant DB: " + err.Error()})
	}
	defer tenantDB.Close()

	token, _ := jwtconnections.GenerateJWT(u.UserID, u.CompanyID, u.CompanyName, u.DBName,u.UserCode)

	return c.RenderJSON(map[string]interface{}{
		"message":       "Login successful",
		"token":         token,
		"company_name":  u.CompanyName,
		"database_name": u.DBName,
	})
}

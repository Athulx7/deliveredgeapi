package controllers

import (
	"DeliverEdgeapi/app/jwtconnections"
	"DeliverEdgeapi/app/services"
	"fmt"

	"github.com/revel/revel"
)

type MenuController struct {
	*revel.Controller
}

func (c MenuController) ListMenus() revel.Result {

	body := make(map[string]interface{})
	if err := c.Params.BindJSON(&body); err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid request body"})
	}

	token, err := jwtconnections.ExtractToken(body)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	claims, err := jwtconnections.ValidateJWT(token)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid Token"})
	}

	dbName := claims.DBName
	roleID := claims.RoleID

	db := services.TenantDBs[dbName]
	if db == nil {
		return c.RenderJSON(map[string]string{
			"error": "Tenant DB not found. Please login again.",
		})
	}

	query := `
	SELECT id, label, route,icon, is_active, isSideMenuActive 
	FROM tbl_menus 
	WHERE role_id = ? AND is_active = 1
	`
	rows, err := db.Query(query, roleID)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": fmt.Sprintf("DB query error: %v", err)})
	}
	defer rows.Close()

	menus := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, url, icon string
		var isActive, isSideMenuActive bool

		if err := rows.Scan(&id, &name, &url, &icon, &isActive, &isSideMenuActive); err != nil {
			return c.RenderJSON(map[string]string{"error": fmt.Sprintf("Row scan error: %v", err)})
		}

		menus = append(menus, map[string]interface{}{
			"id":               id,
			"name":             name,
			"url":              url,
			"icon":             icon,
			"is_active":        isActive,
			"isSideMenuActive": isSideMenuActive,
		})
	}

	return c.RenderJSON(map[string]interface{}{
		"success": true,
		"role_id": roleID,
		"menus":   menus,
	})
}

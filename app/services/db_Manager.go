package services

import (
    "database/sql"
    "fmt"

    _ "github.com/microsoft/go-mssqldb"
    "github.com/revel/revel"
)

var AdminDB *sql.DB

func InitAdminDB() {
    connStr := revel.Config.StringDefault("db.admin.conn", "")
    if connStr == "" {
        revel.AppLog.Fatal("❌ Missing db.admin.conn in app.conf")
    }

    db, err := sql.Open("mssql", connStr)
    if err != nil {
        revel.AppLog.Fatal("❌ SQL Open Error:", err)
    }

    // Test connection
    err = db.Ping()
    if err != nil {
        revel.AppLog.Fatal("❌ Cannot connect to SQL Server:", err)
    }

    revel.AppLog.Info("✅ Connected to SQL Server Admin DB successfully")
	fmt.Println("✅ Connected to SQL Server Admin DB successfully")

    AdminDB = db
}

func ConnectTenantDB(dbHost, dbUser, dbPass, dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", dbHost, dbUser, dbPass, dbName)
	db, err := sql.Open("mssql", connStr)
	if err != nil {
        revel.AppLog.Errorf("❌ Error opening tenant DB connection for %s: %v", dbName, err)
		fmt.Println("Error in the connectiTenantDB 😭😳🤪😵‍💫", dbName, err)
        return nil, err
    }
	if err = db.Ping(); err != nil {
        revel.AppLog.Errorf("❌ Cannot connect to tenant DB (%s): %v", dbName, err)
		fmt.Println("Error in the connectiTenantDB 😭😳🤪😵‍💫", dbName, err)
        return nil, err
    }
	revel.AppLog.Infof("✅ Connected to Tenant DB: %s", dbName)
	fmt.Println("Suucessfulll in the connectiTenantDB ✅ Connected 🤪😭💕😘😍❤️😂", dbName, err)
    return db, nil
}

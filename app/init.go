package app

import (
	"DeliverEdgeapi/app/services"
	"fmt"

	_ "github.com/revel/modules"
	"github.com/revel/revel"
)

var (
	AppVersion string
	BuildTime  string
)

func init() {
	revel.Filters = []revel.Filter{
		CORSFilter,
		revel.PanicFilter,
		revel.RouterFilter,
		revel.FilterConfiguringFilter,
		revel.ParamsFilter,
		revel.SessionFilter,
		revel.FlashFilter,
		revel.ValidationFilter,
		revel.I18nFilter,
		HeaderFilter,
		revel.InterceptorFilter,
		revel.CompressFilter,
		revel.BeforeAfterFilter,
		revel.ActionInvoker,
	}

	fmt.Println("Inside Init.go — connecting to DeliverEdge Admin DB ...")
	revel.OnAppStart(services.InitAdminDB)
	fmt.Println("✅ Connected to DeliverEdge Admin DB Successfully!")

	revel.AppLog.Infof("Running in %s mode", revel.RunMode)
}

func CORSFilter(c *revel.Controller, fc []revel.Filter) {
	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Response.Out.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Response.Out.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Response.Out.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.Response.SetStatus(200)
		return
	}

	revel.AppLog.Infof("✅ CORS Filter hit: %s [%s]", c.Request.URL.Path, c.Request.Method)
	fc[0](c, fc[1:])
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:])
}

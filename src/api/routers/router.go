// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//api/color_schemes/get
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/beautify",
			beego.NSInclude(
				&controllers.BaseController{},
			),
		),
	)
	beego.AddNamespace(ns)

	// if beego.BConfig.RunMode == "dev" || beego.BConfig.RunMode == "testing" {
	// 	var filterLog = func(ctx *context.Context) {
	// 		dataformat := []logger.FormatInfo{
	// 			{"================ this is a request with response start ================", ""},
	// 			{"request URI:   ", ctx.Request.RequestURI},
	// 			{"request header:", ctx.Request.Header},
	// 			{"request body:  ", ctx.Request.Body},
	// 			{"response body: ", ctx.Input.Data()},
	// 			{"================ this is a request with response end ==================", ""},
	// 		}

	// 	}
	// 	beego.InsertFilter("/*", beego.AfterExec, filterLog, false)
	// }
}

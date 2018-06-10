package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["api/controllers:BaseController"] = append(beego.GlobalControllerRouter["api/controllers:BaseController"],
		beego.ControllerComments{
			Method: "SchemeDetial",
			Router: `/schemeDetial`,
			AllowHTTPMethods: []string{"POST"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["api/controllers:BaseController"] = append(beego.GlobalControllerRouter["api/controllers:BaseController"],
		beego.ControllerComments{
			Method: "SaveSchemes",
			Router: `/schemes`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["api/controllers:BaseController"] = append(beego.GlobalControllerRouter["api/controllers:BaseController"],
		beego.ControllerComments{
			Method: "SaveSchemes_Old",
			Router: `/schemes`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}

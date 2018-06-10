package comm

import (
	"reflect"

	"github.com/astaxie/beego"
)

var statusText = map[int]string{
	200: "successful",
	201: "查不到数据",
	//系统级错误
	400: "参数错误",
	500: "内部服务错误",
	501: "保存配色方案失败",
	502: "更新配色方案失败",
	503: "保存配色方案详情失败",
}

//ResultJSON 通用数据返回格式
type ResultJSON struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//GetCodeText 返回对应的错误码说明
func GetCodeText(code int) string {
	str, ok := statusText[code]
	if ok {
		return str
	}
	return statusText[1000]
}

//ResultCodeObj 返回错误码
func ResultCodeObj(code int) ResultJSON {
	var _initial = ResultJSON{
		Code: code,
		Msg:  "",
	}

	if code != 200 {
		_initial.Msg = GetCodeText(code)
	}

	return _initial
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

//FuncExit 返回错误码
func CtrlFuncExit(code int, err error, c beego.Controller) {
	R := ResultCodeObj(code)
	c.Data["json"] = R
	c.ServeJSON()
	if err != nil {
		c.Ctx.Abort(code, err.Error())
	} else {
		c.Ctx.Abort(code, R.Msg)
	}

}

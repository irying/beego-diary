package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"gin-blog/pkg/exception"
	"net/http"
	"gin-blog/models"
	"gin-blog/pkg/util"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/logging"
)


// @Summary 获取多个文章标签
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{"lists":[{"id":3,"created_on":1516849721,"modified_on":0,"name":"3333","created_by":"4555","modified_by":"","state":0}],"total":29},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context)  {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	code := exception.SUCCESS

	c.JSON(http.StatusOK, gin.H{
		"code":code,
		"msg":exception.GetMsg(code),
		"data":data,
	})
}

// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context)  {
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}

	valid.Required(name, "name").Message("名称不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100,"created_by").Message("创建人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Range(state, 0 ,1, "state").Message("状态只允许0或1")

	
	code := exception.INVALID_PARAMS

	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			data := make(map[string]interface{})
			data["name"] = name
			data["created_by"] = createdBy
			data["state"] = state
			models.AddTag(data)
			code = exception.SUCCESS
		} else {
			code = exception.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": exception.GetMsg(code),
		"data": make(map[string]string),
	})
}
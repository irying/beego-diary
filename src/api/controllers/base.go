package controllers

import (
	"api/comm"
	"api/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

// Operations about base
type BaseController struct {
	beego.Controller
}
type saveSchemesStruct struct {
	ID       string `json:"id"`
	Property struct {
		DecorateType   string `json:"decorateType"`
		DecorateOverly string `json:"decorateOverly"`
		StylePrefer    string `json:"stylePrefer"`
		Name           string `json:"name"`
		ApplyCategory  string `json:"applyCategory"` //"title,text,contents,sectionTitle,endPage"
	}
	Color map[string][]string `json:"color"`
}

//ParamStructSchemeDetial 获取配色方案原始JSON接口参数
type ParamStructSchemeDetial struct {
	ColorSchemeID string
	Category      string
}

type colorCount struct {
	//SchemesID     int    //配色方案的ID
	BackGround int //各层背景色数量
	//Level         int //层级
	Decoration    int //各层装饰色数量 i
	DiagramSeries int //各层图示系列色数量 s
	Substrate     int //各层衬底色数量
	// SectionNumber int    //节编号色e 数量
	// Title         int    //各层标题色数量
	// Text          int    //各层正文色数量
	DirectItem int //直接对象颜色组
}

func init() {
	orm.RegisterModel(new(models.Schemes), new(models.SchemesItems), new(models.SchemesJSON), new(models.SchemesColor))
}

// @Title 上报配色方案
// @Description  上报配色方案,配色工具调用该接口将抽取的配色方案保存到数据库
// @Param	json    	        body 	string	true		"上报的JSON数据"
// @Success 200 {string} 保存成功
// @Failure 失败原因，对应错误码
// @router /schemes [post]
func (u *BaseController) SaveSchemes() {
	//存方案原始表

}

// @Title 上报配色方案
// @Description  上报配色方案,配色工具调用该接口将抽取的配色方案保存到数据库
// @Param	json    	        body 	string	true		"上报的JSON数据"
// @Success 200 {string} 保存成功
// @Failure 失败原因，对应错误码
// @router /schemes [post]
func (u *BaseController) SaveSchemes_Old() {

	requestBody := saveSchemesStruct{}
	body := strings.Replace(string(u.Ctx.Input.RequestBody), "\r\n", "", -1)
	if err := json.Unmarshal([]byte(body), &requestBody); err != nil {
		comm.CtrlFuncExit(500, err, u.Controller)
	}

	strJSON, err := json.Marshal(requestBody)
	if err != nil {
		comm.CtrlFuncExit(400, err, u.Controller)
	}

	//遍历MAP
	//keys := reflect.ValueOf(requestBody.Color).MapKeys()
	keys := make([]string, 0, len(requestBody.Color))
	//创建层级字典
	LevelMap := make(map[int]map[string]colorCount)
	LevelMap[0] = make(map[string]colorCount)
	FBGNum := 0
	FSubstrate := 0
	//idworker, _ := comm.NewIdWorker(1)
	//NID, _ := idworker.NextId()
	//fmt.Println(NID)
	NID := requestBody.ID
	templateID := strings.Split(NID, "_")[0]
	schemesColors := []models.SchemesColor{}
	existKeysMAP := []string{}
	for key, value := range requestBody.Color {
		keyIsAdded := true
		keys = append(keys, key)
		//截取key
		keysplit := strings.Split(key, "_")
		layer := len(keysplit) - 1

		//去掉a ,f, i, e后缀
		newKey := strings.Join(keysplit[0:layer], "_")
		if layer <= 1 {
			if layer == 0 {
				FBGNum++
			} else if strings.ToLower(keysplit[layer][0:1]) == "c" {
				FSubstrate++
			}
			newKey = key
		} else if layer > 1 {
			layer--
		}
		parentKey := strings.Join(keysplit[0:layer], "_") //b0_c0_s0 ==> b0_c0
		if parentKey == "" {
			parentKey = newKey
		}

		//更新父节点的父节点数量
		layerIndex := len(strings.Split(parentKey, "_")) - 1
		parentPLv, isExist := LevelMap[layerIndex]
		//fmt.Println(newKey, parentKey, layerIndex)
		if !isExist {
			parentPLv = make(map[string]colorCount)
			LevelMap[layerIndex] = parentPLv
		}

		p := parentPLv[parentKey]
		if !comm.Contain(newKey, existKeysMAP) {
			existKeysMAP = append(existKeysMAP, newKey)
			keyIsAdded = false
		}

		lastChar := strings.ToLower(keysplit[layer][0:1])
		if !keyIsAdded {
			switch lastChar {
			case "s": //图示项目系列颜色组
				p.DiagramSeries++
			case "c": //衬底
				p.Substrate++
			case "d": //直接对象颜色组
				p.DirectItem++
			case "i": //图示装饰颜色组
				p.Decoration++
			}
		}
		parentPLv[parentKey] = p
		if layerIndex != layer {
			lv, isExist := LevelMap[layer]
			if !isExist {
				lv = make(map[string]colorCount)
				LevelMap[layer] = lv
			}
			if _, hasKey := lv[newKey]; !hasKey {
				lv[newKey] = colorCount{}
			}
		}

		//存颜色
		for _, c := range value {
			colors := strings.Split(c, ",")
			base := colors[0]
			contrast := colors[1]
			schemesColors = append(schemesColors, models.SchemesColor{
				SchemeID:   NID,
				TemplateID: templateID,
				Base:       base,
				Contrast:   contrast,
				ItemKey:    key,
				Type:       lastChar,
			})
		}
	}

	timestamp := time.Now().Unix()
	o := orm.NewOrm()
	o.Begin()
	schemes := models.Schemes{Name: requestBody.Property.Name, ApplyPage: requestBody.Property.ApplyCategory, Levels: len(LevelMap), Style: requestBody.Property.StylePrefer,
		CreateTime: timestamp, TemplateID: templateID, Background: FBGNum, Substrate: FSubstrate, SchemeID: NID}
	//fmt.Printf("%+v\n", schemes)
	schemesItems := []models.SchemesItems{}
	for l, data := range LevelMap {
		for key, colorMap := range data {
			schemesItems = append(schemesItems, models.SchemesItems{
				Level:         l + 1,
				Background:    colorMap.BackGround,
				Decoration:    colorMap.Decoration,
				DiagramSeries: colorMap.DiagramSeries,
				Directitem:    colorMap.DirectItem,
				SchemeID:      NID,
				TemplateID:    templateID,
				ItemKey:       key,
			})
		}
	}
	//fmt.Printf("%+v\n", schemesItems)
	//查找配色方案是否存在
	findschemes := o.QueryTable("schemes").Filter("scheme_id", NID)
	if findschemes.Exist() {
		if _, err := findschemes.Update(orm.Params{
			"Name":       schemes.Name,
			"ApplyPage":  schemes.ApplyPage,
			"Levels":     schemes.Levels,
			"Style":      schemes.Style,
			"UpdateTime": schemes.CreateTime,
			"Background": schemes.Background,
			"Substrate":  schemes.Substrate,
		}); err != nil {
			o.Rollback()
			comm.CtrlFuncExit(502, err, u.Controller)
		}
	} else {
		_, err := o.Insert(&schemes)
		if err != nil {
			o.Rollback()
			comm.CtrlFuncExit(501, err, u.Controller)
		}
	}

	//存方案原始表
	schemesJSON := models.SchemesJSON{JSON: string(strJSON), SchemeID: NID, Category: schemes.ApplyPage, TemplateID: schemes.TemplateID}
	findschemes = o.QueryTable("schemes_json").Filter("scheme_id", NID)
	if findschemes.Exist() {
		if _, err := findschemes.Update(orm.Params{
			"JSON":     schemesJSON.JSON,
			"SchemeID": schemesJSON.SchemeID,
			"Category": schemesJSON.Category,
		}); err != nil {
			o.Rollback()
			comm.CtrlFuncExit(502, err, u.Controller)
		}
	} else {
		_, err := o.Insert(&schemesJSON)
		if err != nil {
			o.Rollback()
			comm.CtrlFuncExit(501, err, u.Controller)
		}
	}

	//先删除原来的数据
	if res, err := o.Raw("DELETE FROM schemes_items where scheme_id = ?", NID).Exec(); err != nil {
		o.Rollback()
		comm.CtrlFuncExit(503, err, u.Controller)
	} else {
		num, _ := res.RowsAffected()
		fmt.Println("DELETE FROM items nums: ", num)
	}

	if res, err := o.Raw("DELETE FROM schemes_color where scheme_id = ?", NID).Exec(); err != nil {
		o.Rollback()
		comm.CtrlFuncExit(503, err, u.Controller)
	} else {
		num, _ := res.RowsAffected()
		fmt.Println("DELETE FROM colors nums: ", num)
	}

	if num, err := o.InsertMulti(len(schemesItems), schemesItems); err != nil {
		o.Rollback()
		comm.CtrlFuncExit(503, err, u.Controller)
	} else {
		fmt.Printf("Insert %d schemesItems data!\r\n", num)
	}

	if num, err := o.InsertMulti(len(schemesColors), schemesColors); err != nil {
		o.Rollback()
		comm.CtrlFuncExit(503, err, u.Controller)
	} else {
		fmt.Printf("Insert %d schemesColors data!\r\n", num)
	}

	o.Commit()
	//fmt.Println(insertItems[0], strings.Join(insertItems[1:], ","))
	fmt.Printf("%+v\n", LevelMap[0])

	RETJson := comm.ResultCodeObj(200)
	RETJson.Data = LevelMap
	u.Data["json"] = RETJson
	u.ServeJSON()
}

// SchemeDetial 获取本色方案原始数据
// @Title 获取本色方案原始数据
// @Description  客户端获取某个配色方案原始JSON数据
// @router /schemeDetial [POST]
func (u *BaseController) SchemeDetial() {
	requestBody := ParamStructSchemeDetial{}
	body := strings.Replace(string(u.Ctx.Input.RequestBody), "\r\n", "", -1)
	if err := json.Unmarshal([]byte(body), &requestBody); err != nil {
		comm.CtrlFuncExit(500, err, u.Controller)
	}
	if requestBody.ColorSchemeID == "" {
		comm.CtrlFuncExit(400, nil, u.Controller)
	}
	o := orm.NewOrm()
	SID := strings.Split(requestBody.ColorSchemeID, "_")
	QueryID := SID[0]
	QueryField := "template_id"
	if len(SID) > 1 {
		QueryID = requestBody.ColorSchemeID
		QueryField = "scheme_id"
	}
	//schemesjson := models.SchemesJSON{SchemeID: requestBody.ColorSchemeID, Category: requestBody.Category}
	var schemesjson []models.SchemesJSON
	sql := "SELECT * FROM schemes_json WHERE " + QueryField + " = ?"
	fmt.Println(requestBody.Category)
	if requestBody.Category != "" {
		sql += " and category = '" + requestBody.Category + "'"
	}
	num, err := o.Raw(sql, QueryID).QueryRows(&schemesjson)
	if err != nil {
		comm.CtrlFuncExit(500, nil, u.Controller)
	} else if num == 0 {
		comm.CtrlFuncExit(201, nil, u.Controller)
	} else {
		// var data = make(map[string]string, num)
		// for field, item := range schemesjson {
		// 	fmt.Println(field, item)
		// }
		// fmt.Println(data)
		RETJson := comm.ResultCodeObj(200)
		RETJson.Data = schemesjson
		u.Data["json"] = RETJson
		u.ServeJSON()
	}

	// err := o.Read(&schemesjson)
	// if err == orm.ErrNoRows {
	// 	comm.CtrlFuncExit(201, nil, u.Controller)
	// } else {
	// 	RETJson := comm.ResultCodeObj(200)
	// 	RETJson.Data = schemesjson.JSON
	// 	u.Data["json"] = RETJson
	// 	u.ServeJSON()
	// }
}

// func getColorCount(cc *colorCount, keystr string) {
// 	//找到衬底数量
// 	substrateMath := findKey(`(`+cc.Key+`_c\d),\b`, keystr)
// 	cc.Substrate = len(substrateMath)
// 	//fmt.Println(substrateMath, len(substrateMath))
// 	//找直接对象颜色组
// }

func findKey(ex string, keystr string) [][]string {
	r := regexp.MustCompile(ex)
	match := r.FindAllStringSubmatch(keystr, -1)

	return match
}

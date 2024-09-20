package tool

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/url"
	"runtime"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func DetailErr(errStr string, c *gin.Context) {
	pc, _, _, _ := runtime.Caller(1)

	file, line := runtime.FuncForPC(pc).FileLine(pc)
	logger.Error(file, line, errStr)
	data := "服务器错误"
	var hit = false
	if strings.Contains(errStr, "Duplicate entry") {
		data = "数据重复,请检查名称或其它参数"
		hit = true
	} else if strings.Contains(errStr, "已经预约过了") {
		data = "已经预约过了，请在预约端查询入口权限"
		hit = true
	} else if strings.Contains(errStr, "strconv.Atoi: parsing :") {
		data = "缺少必要传参,请开发人员核查接口"
		hit = true
	} else if strings.Contains(errStr, "Unknown column") {
		strSplit := strings.Split(errStr, "Unknown column")
		data += ": 不清楚的字段名 " + strSplit[1]
		hit = true
	}
	if IsChinese(errStr) && !hit {
		data = errStr
	}

	c.JSON(401, gin.H{"data": data, "err": errStr})
}

func DefaultStuct(c *gin.Context) map[string]interface{} {
	var _temp, _ = io.ReadAll(c.Request.Body)
	var obj map[string]interface{}
	json.Unmarshal(_temp, &obj)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(_temp))

	return obj
}

func GetParams(c *gin.Context) {
	var queyrs = c.Request.URL.Query()
	for key, query := range queyrs {
		st := template.HTMLEscapeString(query[len(query)-1])
		params, _ := url.ParseQuery(c.Request.URL.RawQuery)
		params.Set(key, st)
		c.Request.URL.RawQuery = params.Encode()
	}
}

func GetGinTemp(c *gin.Context) []byte {
	var _temp, _ = io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(_temp))
	return _temp
}

func StructToMap(_struct interface{}) map[string]interface{} {
	var _map = make(map[string]interface{})
	_json, _ := json.Marshal(_struct)
	json.Unmarshal(_json, &_map)
	return _map
}

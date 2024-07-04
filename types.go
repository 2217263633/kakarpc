package myrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/url"
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
	logger.Error(errStr, c.Request.URL.Path, c.Request.Method)
	data := "服务器错误"

	if strings.Contains(errStr, "Duplicate entry") {
		data = "数据重复,请检查名称或其它参数"
	} else if strings.Contains(errStr, "已经预约过了") {
		data = "已经预约过了，请在预约端查询入口权限"
	} else if strings.Contains(errStr, "strconv.Atoi: parsing :") {
		data = "缺少必要传参,请开发人员核查接口"
	} else if strings.Contains(errStr, "Unknown column") {
		strSplit := strings.Split(errStr, "Unknown column")
		data += ": 不清楚的字段名 " + strSplit[1]
	}
	if IsChinese(errStr) {
		data = errStr
	}

	c.JSON(401, gin.H{"data": data, "err": errStr})
}

type Orange struct {
	Size int
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
		// logger.Info(key, query, st)
		// var rsu = make([]string, 0)
		// queyrs[key] = []string{st}
		params.Set(key, st)
		c.Request.URL.RawQuery = params.Encode()
		// logger.Info(queyrs[key])
		// c.Request.URL.Query().Set(key, st)
	}
	// logger.Info(c.Request.URL.Query())
}

type SqlStruct struct {
	Values     string // 不要写 select
	Tabel_name string // 不要写 from
	Where      string // 自己写 where 或者 on
	Order      string // 要写 order by
}

func (_sql SqlStruct) ToString() string {
	sql := fmt.Sprintf(`select %s from %s `,
		_sql.Values, _sql.Tabel_name)

	if _sql.Where != "" {
		sql += _sql.Where
	}

	if _sql.Order != "" {
		sql += " " + _sql.Order
	}

	return sql
}

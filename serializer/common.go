/**
 * @Author: LOFTER
 * @Description:
 * @File:  common
 * @Date: 2021/2/6 10:20 下午
 */
package serializer

import (
	"time"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// DataList 基础列表结构
type DataList struct {
	Items interface{} `json:"items"`
	Total uint        `json:"total"`
}

// BuildListResponse 列表构建器
func BuildListResponse(items interface{}, total uint, code int) Response {
	return Response{
		Data: DataList{
			Items: items,
			Total: total,
		},
		Code: code,
	}
}

func FormatTimeRFC(name time.Time) string {
	return name.Format("2006-01-02 15:04:05")
}

//时间字符串=>时间戳
func FormatTimeStamp(name string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", name, loc)
	return tt.Unix()

}

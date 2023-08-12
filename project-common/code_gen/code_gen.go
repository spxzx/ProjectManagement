package code_gen

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"html/template"
	"log"
	"os"
	"strings"
)

func connectMysql() *gorm.DB {
	//配置MySQL连接参数
	dsn := fmt.Sprintf("root:spx@yp@tcp(127.0.0.1:3309)/project_management?charset=utf8&parseTime=True&loc=Local")
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	return db
}

type Result struct {
	Field string
	Type  string
}
type StructResult struct {
	StructName string
	Result     []*Result
}
type MessageResult struct {
	MessageName string
	Result      []*Result
}

func Name(name string) string {
	names := name[:]
	isSkip := false
	var sb strings.Builder
	for idx, val := range names {
		if idx == 0 {
			s := names[:idx+1]
			s = strings.ToUpper(s)
			sb.WriteString(s)
			continue
		}
		if isSkip {
			isSkip = false
			continue
		}
		if val == 95 {
			s := names[idx+1 : idx+2]
			s = strings.ToUpper(s)
			sb.WriteString(s)
			isSkip = true
			continue
		} else {
			s := names[idx : idx+1]
			sb.WriteString(s)
		}
	}
	return sb.String()
}

func getType(t string) string {
	if strings.Contains(t, "bigint") {
		return "int64"
	}
	if strings.Contains(t, "varchar") {
		return "string"
	}
	if strings.Contains(t, "text") {
		return "string"
	}
	if strings.Contains(t, "tinyint") {
		return "int"
	}
	if strings.Contains(t, "int") &&
		!strings.Contains(t, "tinyint") &&
		!strings.Contains(t, "bigint") {
		return "int"
	}
	if strings.Contains(t, "double") {
		return "float64"
	}
	return ""
}

func getMessageType(t string) string {
	if strings.Contains(t, "bigint") {
		return "int64"
	}
	if strings.Contains(t, "varchar") {
		return "string"
	}
	if strings.Contains(t, "text") {
		return "string"
	}
	if strings.Contains(t, "tinyint") {
		return "int32"
	}
	if strings.Contains(t, "int") &&
		!strings.Contains(t, "tinyint") &&
		!strings.Contains(t, "bigint") {
		return "int32"
	}
	if strings.Contains(t, "double") {
		return "double"
	}
	return ""
}

func GenStruct(table string, structName string) {
	db := connectMysql()
	var re []*Result
	db.Raw(fmt.Sprintf("describe %s", table)).Scan(&re)
	for _, v := range re {
		v.Field = Name(v.Field)
		v.Type = getType(v.Type)
	}
	tmpl, err := template.ParseFiles("./struct.tpl")
	log.Println(err)
	sr := StructResult{structName, re}
	tmpl.Execute(os.Stdout, sr)
}

func GenProtoMessage(table string, msgName string) {
	db := connectMysql()
	var re []*Result
	db.Raw(fmt.Sprintf("describe %s", table)).Scan(&re)
	for _, v := range re {
		v.Field = Name(v.Field)
		v.Type = getMessageType(v.Type)
	}
	var fm template.FuncMap = make(map[string]any)
	fm["Add"] = func(v int, add int) int {
		return v + add
	}
	t := template.New("message.tpl")
	t.Funcs(fm)
	tmpl, err := t.ParseFiles("./message.tpl")
	log.Println(err)
	sr := MessageResult{MessageName: msgName, Result: re}
	err = tmpl.Execute(os.Stdout, sr)
	log.Println(err)
}

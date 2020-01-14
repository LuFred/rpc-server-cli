package dep

import (
	"strings"
	"os"
	"io"
	"fmt"
	"regexp"
	"log"
)

var (
	//匹配单张表
	patTableStr="(?i)(CREATE TABLE if not exists[\\s]*`[\\w]*`)[\\s]*[(]{1}[\\s`\\w,\\W]*?(ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;)"
	//匹配表名,取第3个子集
	patTableNameStr="(?i)(CREATE TABLE if not exists[\\s]*`([\\w]*)`)[\\s]*[(]{1}"
	//匹配字段
	patFieldStr="(?i)([a-zA-Z_]+)`[\\s]*?([\\w]+)([(][\\s\\w\\d]*[)]){0,1}[\\s\\w\\d']*?(comment[\\s]*?'([\\s\\w\\W\\d(),<>]*?)'){0,1},"

)
type Field struct {
	DBType string //数据库字段类型
	DBTypeNumber string //例如：若数据库类型为varchar(10),则DBTypeNumber=10
	Comment string //描述
	DBName string //数据库字段名
	GoName string //go中的字段名
	GoType string //go中的字段类型
}
type Entity struct {
	EntityName string //golang entity 名称
	DBName string//db 表名
	Field []Field
}

var entityTemplate=`package entity
type %s struct {
%s
}

func (d *%s) TableName() string {
	return "%s"
}
`

func ReadSqlString (filePath string) (string,error){
	sqlString:=""
	buffer:=make([]byte,2048)
	file, err := os.Open(filePath)
	if err != nil {
		return "",err
	}
	defer file.Close()
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		sqlString=sqlString+string(buffer[:bytesread])
	}
	return sqlString,nil
}
func GetTable(s string)[]string{
	var resu []string
	reg,_:=regexp.Compile(patTableStr)
	allS:=reg.FindAllString(s,-1)

	for _,v:=range allS{
		resu=append(resu,v)
	}
	return resu
}

func GetDBTableName(s string)(tableName string){
	reg,_:=regexp.Compile(patTableNameStr)
	allS:=reg.FindAllStringSubmatch(s,-1)
	for _,v:=range allS{
		if len(v)>2{
			tableName=v[2]
		}

	}
	return
}

func GetFeild(s string)(field []Field){
	reg,_:=regexp.Compile(patFieldStr)
	allS:=reg.FindAllStringSubmatch(s,-1)
	for _,v:=range allS{
		newField:=Field{}
		if len(v)>1{
			newField.DBName=v[1]
			newField.DBType=v[2]
			if v[3]!=""{
				newField.DBTypeNumber=strings.Trim(v[3],"(")
				newField.DBTypeNumber=strings.Trim(newField.DBTypeNumber,")")
			}

			if len(v)>4{
				newField.Comment=v[5]
			}
			if newField.DBName!=""{
				newField.GoName=MapGoFieldName(newField.DBName)
			}
			if newField.DBType!=""{
				newField.GoType=MapGoFieldType(newField.DBType)
			}
		}
		field=append(field,newField)
		//log.Println(newField)

	}
	return
}

func MapGoFieldType(dbType string)(goType string){
	switch dbType {
	case "int","INT":
		goType="int32"
	case "bigint","BIGINT":
		goType="int64"
	case "tinyint","TINYINT":
		goType="bool"
	case "varchar","VARCHAR":
		goType="string"
	case "double","DOUBLE","float","FLOAT":
		goType="float64"

	}
	return
}

func MapGoFieldName(dbName string)(goName string){
	wordList:=strings.Split(dbName,"_")
	for _,v:=range wordList{
		r:=[]rune(v)
		for i:=range r{
			if i==0{
				if r[i]>=97&&r[i]<=122{
					goName+=string(r[i]-32)
				}else{
					goName+=string(r[i])
				}
			}else{
				if r[i]>=65&&r[i]<=96{
					goName+=string(r[i]+32)
				}else{
					goName+=string(r[i])
				}
			}
		}
	}
	return
}

func CreateGOEntity(entityList []Entity){
	//创建实体文件
	dbDir:=fmt.Sprintf("%s%s%s_service%sdb",Root,PathRule,ServiceName,PathRule)
	entityDir:=fmt.Sprintf(dbDir+"%s%s",PathRule,"entity")
	os.RemoveAll(dbDir)
	//hasDir,_:=pathExists(dir)
	//if !hasDir{

	//}
	_= os.Mkdir(dbDir, os.ModePerm)
	_= os.Mkdir(entityDir, os.ModePerm)
	for _,v:=range entityList{
		entityFileName:=fmt.Sprintf(entityDir+"%s%s_entity.go",PathRule,v.DBName)
		//创建go文件
		eFile, err := os.Create(entityFileName)
		if err != nil {
			log.Println(err)
			return
		}
		//读取文件
		fieldStr:=""
		for _,j:=range v.Field{
			fieldStr+=fmt.Sprintf("	%s %s	`db:\"%s\" description:\"%s\"`\r\n",
				j.GoName,j.GoType,j.DBName,j.Comment,
			)
		}
		entityStr:=fmt.Sprintf(entityTemplate,v.EntityName+"Entity",fieldStr,v.EntityName+"Entity",v.DBName)
		eFile.WriteString(entityStr)
		eFile.Close()
	}

}
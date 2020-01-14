package dep

import (
	"fmt"
	"os"
	"log"
	"strings"
)
var protoTemplate=`
syntax = "proto3"; // 指定proto版本
package %s; // 指定包名

%s

`

var protoEnumTemplate=`syntax = "proto3"; // 指定proto版本
package {{serviceName}}; // 指定包名

enum QueryDeleteStatusEnum {
   ALL=0;
   DELETED=1;
   UNDELETED=2;
}

enum QueryOrderByEnum {
    ASC=0;
    DESC=1;
}


`
var protoServiceTemplate=`
syntax = "proto3"; // 指定proto版本
package %s; // 指定包名
%s

service %s {

%s

}
`
func CreateGOProto(entityList []Entity){

	//创建实体文件
	dir:=fmt.Sprintf("%s%s%s_service%s%s",Root,PathRule,ServiceName,PathRule,"proto")
	_= os.Mkdir(dir, os.ModePerm)
	pb:=fmt.Sprintf("%s%s%s_service%sproto%spb",Root,PathRule,ServiceName,PathRule,PathRule)
	_= os.Mkdir(pb, os.ModePerm)
	createEntityProto(entityList)
	createServiceProto(entityList)
	createEnumProto()

}
func createEnumProto(){
	protoFileName:=fmt.Sprintf("%s%s%s_service%sproto%senum.proto",Root,PathRule,ServiceName,PathRule,PathRule)
	//创建proto文件
	file, err := os.Create(protoFileName)
	if err != nil {
		log.Println(err)
		return
	}
	protoContent:=strings.Replace(protoEnumTemplate,"{{serviceName}}",ServiceName,-1)

	file.WriteString(protoContent)
	file.Close()
}
func createServiceProto(entityList []Entity){
	protoFileName:=fmt.Sprintf("%s%s%s_service%sproto%sservice.proto",Root,PathRule,ServiceName,PathRule,PathRule)
	//创建proto文件
	file, err := os.Create(protoFileName)
	if err != nil {
		log.Println(err)
		return
	}
	protoContent:=fmt.Sprintf(protoServiceTemplate,ServiceName,"%s",firstLetterToUpper(ServiceName),"%s")
	//初始化import
	importStr:=""
	for _,v:=range entityList{
		importStr+=fmt.Sprintf("import \"%s.proto\";\r\n",v.DBName)

	}
	protoContent=fmt.Sprintf(protoContent,importStr,"%s")

	//初始化函数
	handStr:=""
	for _,v:=range entityList{
		etyName:=v.EntityName
		handStr+=fmt.Sprintf("    /*\r\n    * %s\r\n    */\r\n",etyName)
		// create
		handStr+=fmt.Sprintf("    //Create%s\r\n"+
			"    rpc Create%s (Create%sRequest) returns (Create%sReply) {};\r\n",etyName,etyName,etyName,etyName)
		// get
		handStr+=fmt.Sprintf("    //Get%sById\r\n"+
			"    rpc Get%sById (Get%sByIdRequest) returns (Get%sByIdReply) {};\r\n",etyName,etyName,etyName,etyName)
		// update
		handStr+=fmt.Sprintf("    //Update%sById\r\n"+
			"    rpc Update%sById (Update%sByIdRequest) returns (Update%sByIdReply) {};\r\n",etyName,etyName,etyName,etyName)
		// delete
		handStr+=fmt.Sprintf("    //Delete%sById\r\n"+
			"    rpc Delete%sById (Delete%sByIdRequest) returns (Delete%sByIdReply) {};\r\n",etyName,etyName,etyName,etyName)
		//get list
		handStr+=fmt.Sprintf("    //Get%s\r\n"+
			"    rpc Get%s (Get%sRequest) returns (Get%sReply) {};\r\n",etyName,etyName,etyName,etyName)

	}
	protoContent=fmt.Sprintf(protoContent,handStr)
	file.WriteString(protoContent)
	file.Close()
}
func createEntityProto(entityList []Entity){
	//有多少表，创建多少个实体proto对象
	for _,v:=range entityList{
		protoFileName:=fmt.Sprintf("%s%s%s_service%sproto%s%s.proto",Root,PathRule,ServiceName,PathRule,PathRule,v.DBName)
		//创建proto文件
		file, err := os.Create(protoFileName)
		if err != nil {
			log.Println(err)
			return
		}
		protoContent:=fmt.Sprintf(protoTemplate,ServiceName,"%s")
		operationStr:=fmt.Sprintf(
			`%s
%s
%s
%s
%s
%s`,
			infoStringInit(v),
			createaStringInit(v),
			getStringInit(v),
			updateStringInit(v),
			deleteStringInit(v),
			getListStringInit(v))
		protoContent=fmt.Sprintf(protoContent,operationStr)
		file.WriteString(protoContent)
		file.Close()
	}
}

func infoStringInit(entity Entity)string{
	var str=""
	etyName:=entity.EntityName
	str+=fmt.Sprintf("message %sInfo { \r\n",etyName)
	for i,v:=range entity.Field{
		str+=fmt.Sprintf("    %s %s = %d;\r\n",v.GoType,v.DBName,i+1)
	}
	str+="}\r\n"
	return str
}
func createaStringInit(entity Entity)string{
	var str=""
	etyName:=entity.EntityName
	//初始化request
	str+=fmt.Sprintf("message Create%sRequest { \r\n",etyName)
	var index=1
	for _,v:=range entity.Field{
		if v.DBName!="id"&&v.DBName!="gmt_create"&&v.DBName!="gmt_modified"{
			str+=fmt.Sprintf("    %s %s = %d;\r\n",v.GoType,v.DBName,index)
			index++
		}
	}
	str+="}\r\n"
	//初始化reply
	str+=fmt.Sprintf("message Create%sReply { \r\n",etyName)
	str+=fmt.Sprintf("    %sInfo data=1;\r\n"+
							"}\r\n",etyName)


	return str
}
func deleteStringInit(entity Entity)string{
	var str=""
	idType:="int64"
	for _,v:=range entity.Field{
		if v.DBName=="id"{
			idType=v.GoType
		}
	}
	etyName:=entity.EntityName
	str+=fmt.Sprintf(
		"message Delete%sByIdRequest{\r\n"+
   				"    %s id = 1;\r\n"+
				"}\r\n"+
				"message Delete%sByIdReply{\r\n"+
				"    int32 err_code=1;\r\n"+
				"    string err_msg=2;\r\n"+
				"}\r\n",
				etyName,idType,etyName)

	return str
}
func getStringInit(entity Entity)string{
	var str=""
	idType:="int64"
	for _,v:=range entity.Field{
		if v.DBName=="id"{
			idType=v.GoType
		}
	}
	etyName:=entity.EntityName
	str+=fmt.Sprintf("message Get%sByIdRequest{\r\n    %s id = 1;\r\n}\r\nmessage Get%sByIdReply {\r\n"+
		"    %sInfo data=1;\r\n}\r\n",
		etyName,idType,etyName,etyName)

	return str
}
func getListStringInit(entity Entity)string{
	var str=""
	etyName:=entity.EntityName
	str+=fmt.Sprintf("message Get%sRequest{\r\n    int32 offset = 1;\r\n    int32 limit=2;\r\n}\r\nmessage Get%sReply {\r\n"+
		"    repeated %sInfo data=1;\r\n    int32 total=2;\r\n    }\r\n",
		etyName,etyName,etyName)

	return str
}
func updateStringInit(entity Entity)string{
	var str=""
	etyName:=entity.EntityName
	//初始化request
	str+=fmt.Sprintf("message Update%sByIdRequest { \r\n",etyName)
	var index=1
	for _,v:=range entity.Field{
		if v.DBName!="gmt_create"&&v.DBName!="gmt_modified"{
			str+=fmt.Sprintf("    %s %s = %d;\r\n",v.GoType,v.DBName,index)
			index++
		}
	}
	str+="}\r\n"
	//初始化reply
	str+=fmt.Sprintf("message Update%sByIdReply { \r\n",etyName)
	str+=fmt.Sprintf("    %sInfo data=1;\r\n"+
							 "    int32 err_code=2;\r\n"+
							 "    string err_msg=3;\r\n}\r\n",etyName)

	return str
}
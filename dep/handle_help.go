package dep

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var handleCoreBaselTemplate = `package core


`
var handleCoreErrorlTemplate = `package core

import (
	"fmt"

	"github.com/lufred/{{serviceName}}_service/util/log"
	"github.com/lufred/{{serviceName}}_service/config"
)


// LogicError .
type LogicError struct {
	Code         int
	Message      string
	Type         string //system:系统错误|logic ：业务逻辑错误
	SeriviceName string
	handle       string
	Err          error


}
const(
	SystemErrorType="[system]"
	LogicErrorType="[logic]"
)
func (e LogicError) Error() string {
	if e.Err != nil {
		log.Debugf(fmt.Sprintf("ServiceName:%s Handle:%s ErrorType:%s Code:%d Message:%s error:%s", e.SeriviceName, e.handle,e.Type, e.Code, e.Message, e.Err.Error()))
		return fmt.Sprintf("ServiceName:%s Handle:%s ErrorType:%s Code:%d Message:%s error:%s", e.SeriviceName, e.handle,e.Type, e.Code, e.Message, e.Err.Error())
	}
	log.Debugf(fmt.Sprintf("ServiceName:%s Handle:%s ErrorType:%s Code:%d Message:%s", e.SeriviceName, e.handle,e.Type, e.Code, e.Message))
	return fmt.Sprintf("ServiceName:%s Handle:%s ErrorType:%s Code:%d Message:%s", e.SeriviceName, e.handle,e.Type, e.Code, e.Message)

}

//ProgramException 程序错误
func ProgramException(hd string, e error) error {
	err := LogicError{
		SeriviceName: config.ProConfig.ServiceName,
		handle:       hd,
		Err:          e,
		Type:SystemErrorType,
	}
	log.Error(err)
	return err
}

//LogicException 逻辑错误
func LogicException(code int, msg, hd string, e error) error {
	err := LogicError{
		SeriviceName: config.ProConfig.ServiceName,
		handle:       hd,
		Err:          e,
		Type:        LogicErrorType,
		Code:         code,
		Message:      msg,
	}
	log.Error(err)
	return err
}


`

var handleErrorlMapTemplate = `package core

// error status code
const (
	ErrorStatusParameterError = 4000
	ErrorStatusNotPermission  = 4001
	ErrorStatusAlreadyExists  = 4002
	ErrorStatusInternalServer = 4003
	ErrorStatusNotFound       = 4004

)

var errorStatusText = map[int]string{
	ErrorStatusParameterError:   "Parameter Error",
	ErrorStatusAlreadyExists:  "Object Already Exists",
	ErrorStatusNotFound:       "Object Not Found",
	ErrorStatusInternalServer: "Internal Server Error",
	ErrorStatusNotPermission: "Unauthorized",
}

func ErrorStatusText(code int) string {
	return errorStatusText[code]
}
`
var handleLogicBaseTemplate = `package logic

//{{serviceHandleStruct}}Server
type {{serviceHandleStruct}}Server struct{}

`
var handleLogicTemplate = `package logic

import (
	"github.com/lufred/goorm"
	"golang.org/x/net/context"
	pb "github.com/lufred/{{serviceName}}_service/proto/pb"
	"github.com/lufred/{{serviceName}}_service/handle/core"
	"github.com/lufred/{{serviceName}}_service/db/entity"
	"github.com/lufred/{{serviceName}}_service/util"
    "math"
)

//Create{{goDbEntity}}
func (s *{{serviceHandleStruct}}Server) Create{{goDbEntity}}(context context.Context, in *pb.Create{{goDbEntity}}Request) (*pb.Create{{goDbEntity}}Reply, error) {

	r := new(pb.Create{{goDbEntity}}Reply)
	o := goorm.NewOrm()
	curTime := util.GetTimeMillisecond()
	ety := &entity.{{goDbEntity}}Entity{
		{{entityInsertAllField}}
	}
	id,err := o.Insert(ety)
	if err != nil {
		return r, core.ProgramException("Create{{goDbEntity}}", err)
	}
	{{idToid}}
	r.Data = convert{{goDbEntity}}EntityToPb(ety)

	return r, nil
}


//Update{{goDbEntity}}ById
func (s *{{serviceHandleStruct}}Server) Update{{goDbEntity}}ById(context context.Context, in *pb.Update{{goDbEntity}}ByIdRequest) (*pb.Update{{goDbEntity}}ByIdReply, error) {

	r := new(pb.Update{{goDbEntity}}ByIdReply)
	o := goorm.NewOrm()
	curTime := util.GetTimeMillisecond()
	ety := &entity.{{goDbEntity}}Entity{}
	selecter:=goorm.Cond{
		"id": in.Id,
	}
	err:=o.One(ety,selecter)
	if err!=nil{
		if err==goorm.ErrNoMoreRows{
			r.ErrCode=core.ErrorStatusNotFound
			r.ErrMsg=core.ErrorStatusText(core.ErrorStatusNotFound)
		}else{
			r.ErrCode=core.ErrorStatusInternalServer
			r.ErrMsg=core.ErrorStatusText(core.ErrorStatusInternalServer)
			return r, core.ProgramException("Update{{goDbEntity}}ById", err)
		}
	}
	ety.GmtModified=curTime
	{{updataEntityAllField}}
    o.Update(ety)
	if err != nil {
		return r, core.ProgramException("Update{{goDbEntity}}ById", err)
	}
	r.Data = convert{{goDbEntity}}EntityToPb(ety)
	return r, nil
}


//Get{{goDbEntity}}ById
func (s *{{serviceHandleStruct}}Server) Get{{goDbEntity}}ById(context context.Context, in *pb.Get{{goDbEntity}}ByIdRequest) (*pb.Get{{goDbEntity}}ByIdReply, error) {
	r := new(pb.Get{{goDbEntity}}ByIdReply)
	o := goorm.NewOrm()
	ety := &entity.{{goDbEntity}}Entity{}
	selecter:=goorm.Cond{
		"id": in.Id,
	}
	err:=o.One(ety,selecter)
	if err!=nil{
		if err==goorm.ErrNoMoreRows{
			return r,nil
		}else{
			return r, core.ProgramException("Get{{goDbEntity}}ById", err)
		}
	}
	r.Data = convert{{goDbEntity}}EntityToPb(ety)
	return r, nil
}

//Get{{goDbEntity}}
func (s *{{serviceHandleStruct}}Server) Get{{goDbEntity}}(context context.Context, in *pb.Get{{goDbEntity}}Request) (*pb.Get{{goDbEntity}}Reply, error) {
	r := new(pb.Get{{goDbEntity}}Reply)
	o := goorm.NewOrm()
	var list []entity.{{goDbEntity}}Entity
	var offset,limit int32
	selecter := goorm.Cond{}
	if in.Limit==-1{
		limit=math.MaxInt32
	}else{
		limit=in.Limit
	}
	if in.Offset<0{
		offset=0
	}else{
		offset=in.Offset
	}
	count, _ := o.Count(&entity.{{goDbEntity}}Entity{}, selecter)
	r.Total=int32(count)
	if count < int64(offset)|| limit <1 {
		return r, nil
	}
	err:=o.SelectLimit(&list,selecter,int64(offset),int64(limit))
	if err!=nil{
			return r, core.ProgramException("Get{{goDbEntity}}", err)
	}
	if len(list) > 0 {
		r.Data = make([]*pb.{{goDbEntity}}Info, len(list))
		for i := range list {
			r.Data[i] = convert{{goDbEntity}}EntityToPb(&list[i])

		}
	} else {
		r.Data = make([]*pb.{{goDbEntity}}Info, 0)
	}
	return r, nil
}

//Delete{{goDbEntity}}ById
func (s *{{serviceHandleStruct}}Server) Delete{{goDbEntity}}ById(context context.Context, in *pb.Delete{{goDbEntity}}ByIdRequest) (*pb.Delete{{goDbEntity}}ByIdReply, error) {
	r := new(pb.Delete{{goDbEntity}}ByIdReply)
	o := goorm.NewOrm()
	ety := &entity.{{goDbEntity}}Entity{}
	selecter:=goorm.Cond{
		"id": in.Id,
	}
	err:=o.One(ety,selecter)
	if err!=nil{
		if err==goorm.ErrNoMoreRows{
			r.ErrCode=core.ErrorStatusNotFound
			r.ErrMsg=core.ErrorStatusText(core.ErrorStatusNotFound)
		}else{
			r.ErrCode=core.ErrorStatusInternalServer
			r.ErrMsg=core.ErrorStatusText(core.ErrorStatusInternalServer)
			return r, core.ProgramException("Delete{{goDbEntity}}ById", err)
		}
	}
	_,err=o.Delete(ety)
	if err!=nil{
		r.ErrCode=core.ErrorStatusInternalServer
		r.ErrMsg=core.ErrorStatusText(core.ErrorStatusInternalServer)
		return r, core.ProgramException("Delete{{goDbEntity}}ById", err)

	}
	return r, nil
}

func convert{{goDbEntity}}EntityToPb(in *entity.{{goDbEntity}}Entity) *pb.{{goDbEntity}}Info{
	p:=&pb.{{goDbEntity}}Info{}
	{{convertPbEntity}}
	return p
}
`
var (
	handleBaseDir = ""
)

func init() {

}

func CreateHandle(entityList []Entity) {
	handleBaseDir = BaseDir + PathRule + "handle"
	createDir(handleBaseDir)
	createHandleRun(entityList)

}
func createHandleRun(entityList []Entity) {
	createCore()
	createHandle(entityList)

}
func createHandle(entityList []Entity) {

	coreDir := handleBaseDir + PathRule + "logic"
	createDir(coreDir)
	//create logic base
	file, err := os.Create(coreDir + PathRule + "base.go")
	if err != nil {
		log.Printf("createCore %v", err)
		return
	}
	//grpc 服务实体对象
	s := firstLetterToUpper(ServiceName)
	baseStr := strings.Replace(handleLogicBaseTemplate, "{{serviceHandleStruct}}", s, -1)
	file.WriteString(baseStr)
	file.Close()
	//create logic

	for _, v := range entityList {
		logicFile := fmt.Sprintf(coreDir+PathRule+"%s_logic.go", v.DBName)
		file, err = os.Create(logicFile)
		if err != nil {
			log.Printf("createHandle %v", err)
			return
		}

		logicEntityContentTemplate := handleLogicTemplate
		logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{serviceName}}", ServiceName, -1)
		logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{goDbEntity}}", v.EntityName, -1)
		logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{serviceHandleStruct}}", s, -1)
		entityInsertAllFieldStr := ""
		updataEntityAllFieldStr := ""
		convertPbEntityStr := ""
		for _, f := range v.Field {
			if f.GoName == "GmtCreate" {
				entityInsertAllFieldStr = entityInsertAllFieldStr + f.GoName + ":curTime" + ",\r\n        "
			} else if f.GoName == "GmtModified" {
				entityInsertAllFieldStr = entityInsertAllFieldStr + f.GoName + ":0" + ",\r\n        "
			} else if f.GoName != "Id" {
				entityInsertAllFieldStr = entityInsertAllFieldStr + f.GoName + ":in." + f.GoName + ",\r\n        "
			}

			if f.GoName != "Id" && f.GoName != "GmtCreate" && f.GoName != "GmtModified" {
				updataEntityAllFieldStr = updataEntityAllFieldStr + "ety." + f.GoName + "=in." + f.GoName + "\r\n    "
			}
			convertPbEntityStr = convertPbEntityStr + "p." + f.GoName + "=in." + f.GoName + "\r\n    "
		}
		for _, f := range v.Field {
			if f.GoName != "Id" {
				continue
			}
			if f.GoType == "int32" {
				logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{idToid}}", "ety.Id=int32(id)", -1)
			} else {
				logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{idToid}}", "ety.Id=id", -1)
			}
		}

		logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{entityInsertAllField}}", entityInsertAllFieldStr, -1)
		logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{updataEntityAllField}}", updataEntityAllFieldStr, -1)
		logicEntityContentTemplate = strings.Replace(logicEntityContentTemplate, "{{convertPbEntity}}", convertPbEntityStr, -1)
		file.WriteString(logicEntityContentTemplate)
		file.Close()
	}

}
func createCore() {
	coreDir := handleBaseDir + PathRule + "core"
	createDir(coreDir)
	//创建local配置文件
	file, err := os.Create(coreDir + PathRule + "base.go")
	if err != nil {
		log.Printf("createCore %v", err)
		return
	}
	file.WriteString(handleCoreBaselTemplate)
	file.Close()
	//create error
	file, err = os.Create(coreDir + PathRule + "error.go")
	if err != nil {
		log.Printf("createCore %v", err)
		return
	}
	logicHandleCoreErrorlTemplate := handleCoreErrorlTemplate
	logicHandleCoreErrorlTemplate = strings.Replace(handleCoreErrorlTemplate, "{{serviceName}}", ServiceName, -1)
	file.WriteString(logicHandleCoreErrorlTemplate)
	file.Close()
	//create error map
	file, err = os.Create(coreDir + PathRule + "err_map.go")
	if err != nil {
		log.Printf("createCore %v", err)
		return
	}
	file.WriteString(handleErrorlMapTemplate)
	file.Close()

}

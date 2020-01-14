package main

import (
	"flag"
	"fmt"
	"github.com/lufred/rpc-server-cli/dep"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var (
	serviceName   = flag.String("service", "demo", "服务名称")
	listenAddress = flag.String("listen", "0.0.0.0", "监听地址 默认:0.0.0.0")
	port          = flag.String("port", "4399", "监听端口号 默认:4399")
)

func init() {
	flag.Parse()
	sysType := runtime.GOOS
	switch sysType {
	case "windows":
		// windows系统
		dep.PathRule = "\\"
		dep.ExecSystem = dep.System_Windows
	case "linux":
		dep.ExecSystem = dep.System_Linux
		dep.PathRule = "/"
	default:
		dep.ExecSystem = dep.System_Darwin
		dep.PathRule = "/"
	}
}

func main() {
	root, _ := os.Getwd()
	dep.Root = root
	dep.ServiceName = *serviceName
	dep.Port = *port
	dep.ListenAddress = *listenAddress
	dep.BaseDir = fmt.Sprintf("%s%s%s_service", dep.Root, dep.PathRule, dep.ServiceName)
	dir := fmt.Sprintf("%s%s%s_service", dep.Root, dep.PathRule, dep.ServiceName)
	os.RemoveAll(dir)
	_ = os.Mkdir(dir, os.ModePerm)
	path := fmt.Sprintf("%s%s%s", root, dep.PathRule, "table.sql")
	sqlStr, err := dep.ReadSqlString(path)
	if err != nil {
		log.Println(err)
		return
	}
	tableList := dep.GetTable(sqlStr)

	entityList := make([]dep.Entity, 0)
	for _, v := range tableList {
		newEntity := dep.Entity{}
		dbTableName := dep.GetDBTableName(v)
		newEntity.DBName = dbTableName
		newEntity.EntityName = dep.StringUnderlineToHump(dbTableName, "")
		newEntity.Field = dep.GetFeild(v)
		entityList = append(entityList, newEntity)
	}
	//创建config
	dep.CreateGoConfig()
	//创建dep
	dep.CreateDep()
	//创建doc
	dep.CreateDoc(sqlStr)

	//创建middleware
	dep.CreateMiddleware()
	//创建handle
	dep.CreateHandle(entityList)
	//创建goalng entity
	dep.CreateGOEntity(entityList)
	//创建main
	dep.CreateMain(entityList)

	//创建proto
	dep.CreateGOProto(entityList)

	//生成pb文件
	cmdCreatePb := ""
	cmdCreatePb = fmt.Sprintf(root+"%s%s_service%sproto", dep.PathRule, dep.ServiceName, dep.PathRule)

	err = os.Chdir(cmdCreatePb)
	if err != nil {
		log.Println(err)
	}
	var cmd *exec.Cmd
	switch dep.ExecSystem {
	case dep.System_Windows:
		cmd = exec.Command("sh", "-c", "protoc -I ./pb ./*.proto --go_out=plugins=grpc:./pb/ --proto_path=./")
	case dep.System_Linux:
		cmd = exec.Command("sh", "-c", "protoc -I ./pb ./*.proto --go_out=plugins=grpc:./pb/ --proto_path=./")
	case dep.System_Darwin:
		cmd = exec.Command("sh", "-c", "protoc -I ./pb ./*.proto --go_out=plugins=grpc:./pb/ --proto_path=./")
	default:
		cmd = exec.Command("sh", "-c", "protoc -I ./pb ./*.proto --go_out=plugins=grpc:./pb/ --proto_path=./")
	}
	_, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}

}

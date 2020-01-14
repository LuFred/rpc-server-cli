package dep

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var mainTemplate = `package main

import (
	"flag"
	"github.com/lufred/goorm"
	"net"

	"fmt"
	"google.golang.org/grpc"
	"github.com/lufred/{{serviceName}}_service/config"
	"github.com/lufred/{{serviceName}}_service/dep"
	"github.com/lufred/{{serviceName}}_service/middleware"
	"github.com/lufred/{{serviceName}}_service/util/gopool"
	"github.com/lufred/{{serviceName}}_service/util/log"
	"github.com/lufred/{{serviceName}}_service/util/redis"
	"github.com/lufred/{{serviceName}}_service/handle/logic"
	pb "github.com/lufred/{{serviceName}}_service/proto/pb"
	"github.com/lufred/{{serviceName}}_service/db/entity"
)

var s *grpc.Server

func main() {
	host := fmt.Sprintf("%s:%s", config.ProConfig.Listen, config.ProConfig.Port)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	opInterceptor, ioCloser, err := middleware.InterceptorOpentracing(config.ProConfig.TracingTransportURL, config.ProConfig.ServiceName)
	if err != nil {

		log.Errorf("initGrpc: %s", err.Error())
	} else {
		defer ioCloser.Close()
		dep.UseInterceptor(opInterceptor)
	}
	s=dep.NewGrpcServer()
	pb.Register{{serviceHandleStruct}}Server(s,&logic.{{serviceHandleStruct}}Server{})
	log.Infof("listen %s", host)
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
func init() {
	initConfig()
	initGoPool()
	initLog()
	initORM()
	initRedis()

}

func initConfig() {
	flag.Parse()
	config.RegisterConfig()
}
func initLog() {
	log.Debugged = config.ProConfig.Debug
}
func initGoPool() {
	pool.InitGoPool(1024)
}
func initRedis() {
	redis.EnableCache = false
	redis.InitRedis(config.ProConfig.Redis.Host, config.ProConfig.Redis.Password, config.ProConfig.Redis.Database)
}
func initORM() {
	goorm.SetLogging(true)

	{{ormEntityMap}}
	err := goorm.RegisterDataBase("default",
		"mysql",
		config.ProConfig.Db.Host,
		config.ProConfig.Db.DB,
		config.ProConfig.Db.UserName,
		config.ProConfig.Db.PWD,
		10,
		10,
	)
	if err != nil {
		panic(err)
	}

}

`
var (
	mainBaseDir = ""
)

func CreateMain(entityList []Entity) {
	mainBaseDir = BaseDir
	createMainRun(entityList)
}
func createMainRun(entityList []Entity) {
	file, err := os.Create(mainBaseDir + PathRule + "main.go")
	if err != nil {
		log.Printf("createCore %v", err)
		return
	}
	s := firstLetterToUpper(ServiceName)
	mainStr := strings.Replace(mainTemplate, "{{serviceName}}", ServiceName, -1)
	mainStr = strings.Replace(mainStr, "{{serviceHandleStruct}}", s, -1)
	ormMapStr := ""
	for _, v := range entityList {
		ormMapStr += fmt.Sprintf("goorm.RegisterModel(&entity.%sEntity{})\r\n    ", v.EntityName)
	}
	mainStr = strings.Replace(mainStr, "{{ormEntityMap}}", ormMapStr, -1)
	file.WriteString(mainStr)
	file.Close()

}

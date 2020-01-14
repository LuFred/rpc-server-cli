package dep

import (
	"os"
	"log"
)

var depGrpcTamlTemplate=`
package dep

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"math"
)

var (
	interceptor []grpc.UnaryServerInterceptor
)

func NewGrpcServer() *grpc.Server {
	grpc.MaxRecvMsgSize(math.MaxInt32)
	grpc.MaxSendMsgSize(math.MaxInt32)

	return grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptor...)))

}
func UseInterceptor(unaryServer grpc.UnaryServerInterceptor) {
	interceptor = append(interceptor, unaryServer)
}

`

var (
	depBaseDir=""
)
func init(){

}

func CreateDep(){
	depBaseDir=BaseDir+PathRule+"dep"
	createDir(depBaseDir)
	createDepGo()



}
func createDepGo(){
	//创建local配置文件
	file,err:=os.Create(depBaseDir+PathRule+"grpc.go")
	if err != nil {

		log.Printf("createDepGo %v",err)
		return
	}
	file.WriteString(depGrpcTamlTemplate)
	file.Close()

}
package dep

import (
	"os"
	"log"
)

var midOpentracingTamlTemplate=`
package middleware

import (
	"errors"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"google.golang.org/grpc"
	"io"
)

var ErrNewHTTPTransport = errors.New("opentracing: cannot initialize opentracing http transport")

func InterceptorOpentracing(host, serviceName string) (grpc.UnaryServerInterceptor, io.Closer, error) {
	transport, err := zipkin.NewHTTPTransport(
		host,
		zipkin.HTTPBatchSize(5),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		return nil, nil, ErrNewHTTPTransport
	}
	tracer, closer := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
	)
	var opts []grpc_opentracing.Option
	opts = append(opts, grpc_opentracing.WithTracer(tracer))
	interceptor := grpc_middleware.ChainUnaryServer(
		grpc_opentracing.UnaryServerInterceptor(opts...),
	)
	return interceptor, closer, nil

}

`

var (
	middlewareBaseDir=""
)
func init(){

}

func CreateMiddleware(){
	middlewareBaseDir=BaseDir+PathRule+"middleware"
	createDir(middlewareBaseDir)
	createMiddlewareRun()



}
func createMiddlewareRun(){
	createOpentracingGo()

}
func createOpentracingGo(){
	//创建local配置文件
	file,err:=os.Create(middlewareBaseDir+PathRule+"open_tracing.go")
	if err != nil {

		log.Printf("createOpentracingGo %v",err)
		return
	}
	file.WriteString(midOpentracingTamlTemplate)
	file.Close()
}
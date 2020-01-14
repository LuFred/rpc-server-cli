package dep

import (
	"strings"
	"os"
	"log"
)

var configTamlTemplate=`
# service name
servicename: {{serviceName}}
# enable debug mode
debug: true
# opentracing address
tracingurl: 
# listen address
listen: {{serviceListenAddress}}
#listen port
port: {{servicePort}}
# mysql db
database:
  host:
  username:
  pwd:
  db:
# redis db
redis:
  host:
  password:
  database:

`
var configGoTemplate=`
package config

import (
	"flag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

// Config struct defines the config structure
type Config struct {
	ServiceName         string       `+"`yaml:\"servicename\"`"+`
	Debug               bool         `+"`yaml:\"debug\"`"+`
	TracingTransportURL string       `+"`yaml:\"tracingurl\"`"+`
	Db                  *Database    `+"`yaml:\"database\"`"+`
	Listen              string       `+"`yaml:\"listen\"`"+`
	Port                string        `+"`yaml:\"port\"`"+`
	Redis               *RedisConfig `+"`yaml:\"redis\"`"+`
}

type RedisConfig struct {
	Host     string `+"`yaml:\"host\"`"+`
	Password string `+"`yaml:\"password\"`"+`
	Database int    `+"`yaml:\"database\"`"+`
}

type Database struct {
	Host     string `+"`yaml:\"host\"`"+`
	UserName string `+"`yaml:\"username\"`"+`
	PWD      string `+"`yaml:\"pwd\"`"+`
	DB       string `+"`yaml:\"db\"`"+`
}

var (
	ProConfig = Config{}
	env       = flag.String("env", "local", "运行环境")
	conf      = flag.String("conf", "nil", "配置文件路径")
)

//RegisterConfig 初始化config
func RegisterConfig() {
	var file []byte
	var err error
	if !flag.Parsed() {
		os.Stderr.Write([]byte("ERROR: config before flag.Parse"))
		os.Exit(1)
		return
	}
	root, err := os.Getwd()
	if *env == "local" {
		root = root + "/config/local.yml"
	} else {
		if *conf == "nil" {
			if err != nil {
				log.Fatalln(err.Error())
			}
			root = root + "/config/local.json"
		} else {
			root = *conf
		}
	}
	if err != nil {
		log.Fatalln("ERROR: Read config file error")
		return
	}
	file, err = ioutil.ReadFile(root)
	err = yaml.Unmarshal(file, &ProConfig)
	if err != nil {
		panic(err)
	}
}

`
var (
	configBaseDir=""
)
func init(){

}

func CreateGoConfig(){
	configBaseDir=BaseDir+PathRule+"config"
	createDir(configBaseDir)
	createConfigGo()
	createYml()



}
func createConfigGo(){
	//创建local配置文件
	file,err:=os.Create(configBaseDir+PathRule+"config.go")
	if err != nil {
		log.Printf(configBaseDir+PathRule+"config.go")
		log.Printf("CreateConfigGo %v",err)
		return
	}
	file.WriteString(configGoTemplate)
	file.Close()

}
func createYml(){
	//配置服务名
	content:=strings.Replace(configTamlTemplate,`{{serviceName}}`,ServiceName,-1)
	//配置监听地址
	content=strings.Replace(content,`{{serviceListenAddress}}`,ListenAddress,-1)
	//配置端口号
	content=strings.Replace(content,`{{servicePort}}`,Port,-1)

	//创建local配置文件
	file,err:=os.Create(configBaseDir+PathRule+"local.yml")
	if err != nil {
		log.Println(err)
		return
	}
	file.WriteString(content)
	file.Close()
	//创建dev配置文件
	file,err=os.Create(configBaseDir+PathRule+"dev.yml")
	if err != nil {
		log.Println(err)
		return
	}
	file.WriteString(content)
	file.Close()
	//创建prod配置文件
	file,err=os.Create(configBaseDir+PathRule+"prod.yml")
	if err != nil {
		log.Println(err)
		return
	}
	file.WriteString(content)
	file.Close()



}
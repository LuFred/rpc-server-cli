package dep

import (
	"os"
	"log"
)

var (
	docBaseDir=""
)
func init(){

}

func CreateDoc(sqlString string){
	docBaseDir=BaseDir+PathRule+"doc"
	createDir(docBaseDir)
	createDocRun(sqlString)



}
func createDocRun(sqlString string){
	//创建local配置文件
	file,err:=os.Create(docBaseDir+PathRule+"table.sql")
	if err != nil {

		log.Printf("createDocRun %v",err)
		return
	}
	file.WriteString(sqlString)
	file.Close()

}
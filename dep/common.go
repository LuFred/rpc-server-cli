package dep

import (
	"os"
	"strings"
)
var(
	Root=""
	ServiceName=""
	BaseDir=""
	Port=""
	ListenAddress=""
	GoOS=""
	PathRule=""
	ExecSystem= ""
)
const (
System_Darwin="darwin"
System_Windows="windows"
System_Linux="linux"
)
func createDir(dir string)(bool,error){
	os.RemoveAll(dir)
	err:=os.Mkdir(dir, os.ModePerm)
	if err!=nil{
		return false,err
	}
	return true,nil
}
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
//firstLetterToUpper 首字母转大写
func firstLetterToUpper(s string)string{
	if len(s)>0 {
		r:=[]rune(s)
		if r[0]>=97&&r[0]<=122 {
			r[0] = r[0] - 32
		}
		return string(r)
	}
	return s
}

//字符串下划线转驼峰 例如 a_bc_d => ABcD
//StringUnderlineToHump
func StringUnderlineToHump (dbTableName,suffix string)(goName string){
	wordList:=strings.Split(dbTableName,"_")
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
	goName+=suffix
	return
}
package mvc

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/temp"
	"github.com/galaxy-book/polaris-backend/common/core/util/str"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const FacadeBuildGoFile = `package {{.Package}}

import (
	"errors"
	"fmt"
	{{.CtxPack}}
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	{{.GinUtilPack}}
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/{{.VoPackage}}"
)

{{.FuncDesc}}
`

const FacadeBuildFunc = `
func {{.MethodName}}({{.ArgsDesc}}) {{.OutArgsDesc}} {
	respVo := &{{.OutArgsDesc}}{}
	
	reqUrl := {{.Api}}
	queryParams := map[string]interface{}{}
{{.QueryParamsDesc}}
{{.RequestDesc}}
	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("{{.SvcName}} response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	{{.ReturnDesc}}
}
`

type FacadeBuilder struct {
	StorageDir string
	Package    string
	VoPackage  string
	Greeters   []interface{}
}

type ReqTypeInfo struct {
	reqType reflect.Type
	isPtr   bool
}

var skipMethods = []string{"Health"}

func (fb FacadeBuilder) Build() {
	for _, greeter := range fb.Greeters {
		greeterInfo := GetGreeterInfo(greeter)
		goFileContent := fb.greeterRender(greeter, greeterInfo)

		if strings.Trim(goFileContent, " ") == "" {
			continue
		}
		version := strings.ToLower(greeterInfo.Version)
		if version != "" {
			version += "_"
		}
		filePath := fmt.Sprintf("%s/%s_%s_%sfacade.go", fb.StorageDir, greeterInfo.ApplicationName, strings.ToLower(greeterInfo.HttpType), version)

		data := []byte(goFileContent)
		if ioutil.WriteFile(filePath, data, 0644) == nil {
			fmt.Printf("build success: %s\n", filePath)
		}
	}
}

func (fb FacadeBuilder) greeterRender(greeter interface{}, greeterInfo Greeter) string {
	greeterValue := reflect.ValueOf(greeter)
	greeterType := reflect.TypeOf(greeter)

	funcDesc := ""
	methodNum := greeterType.NumMethod()

	if methodNum == 0 {
		return ""
	}
	hasCtx := false
	hasApi := false
	for i := 0; i < methodNum; i++ {
		methodName := greeterType.Method(i).Name
		if exist, _ := slice.Contain(skipMethods, methodName); exist{
			fmt.Println("skip method", methodName)
			continue
		}
		hasApi = true
		method := greeterValue.MethodByName(methodName)
		subFuncDesc, useCtx := fb.greeterMethodRender(method, methodName, greeterInfo)
		if useCtx {
			hasCtx = true
		}
		funcDesc += subFuncDesc + "\n"
	}
	if ! hasApi{
		return ""
	}

	goFileTemplateData := map[string]string{}
	goFileTemplateData["Package"] = fb.Package
	goFileTemplateData["VoPackage"] = fb.VoPackage
	goFileTemplateData["FuncDesc"] = funcDesc
	if hasCtx {
		goFileTemplateData["CtxPack"] = "\"context\""
		goFileTemplateData["GinUtilPack"] = "\"github.com/galaxy-book/polaris-backend/common/extra/gin/util\""
	} else {
		goFileTemplateData["CtxPack"] = ""
		goFileTemplateData["GinUtilPack"] = ""
	}
	result, _ := temp.Render(FacadeBuildGoFile, goFileTemplateData)
	return result
}

func (fb FacadeBuilder) greeterMethodRender(method reflect.Value, methodName string, greeterInfo Greeter) (string, bool) {
	templateData := map[string]string{}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("err at %s, err %v", methodName, r)
		}
	}()

	//MethodName
	templateData["MethodName"] = methodName + strings.ToUpper(greeterInfo.Version)

	//ArgsDesc
	ptrTab := map[string]ReqTypeInfo{}
	templateData["ArgsDesc"] = fb.buildArgDesc(method, ptrTab)

	_, useCtx := ptrTab["ctx"]

	//ReturnDesc
	templateData["OutArgsDesc"] = fb.buildOutArgsDesc(method, ptrTab)

	//BuildApi
	templateData["Api"] = "fmt.Sprintf(\"%s/" + BuildApi(greeterInfo.ApplicationName, greeterInfo.Version, methodName) + "\", config.GetPreUrl(\"" + greeterInfo.ApplicationName + "\"))"

	//SvcName
	templateData["SvcName"] = greeterInfo.ApplicationName

	//QueryParamsDesc
	templateData["QueryParamsDesc"] = fb.buildQueryParamsDesc(greeterInfo.HttpType, ptrTab)

	//RequestBodyLogDesc
	//templateData["RequestBodyLogDesc"] = fb.buildRequestBodyLogDesc(greeterInfo.HttpType)

	//RequestDesc
	templateData["RequestDesc"] = fb.buildRequestDesc(greeterInfo.HttpType, ptrTab)

	//ReturnDesc
	templateData["ReturnDesc"] = fb.buildReturnDesc(ptrTab)

	result, _ := temp.Render(FacadeBuildFunc, templateData)
	return result, useCtx
}

func GetGreeterInfo(greeter interface{}) Greeter {
	greeterValue := reflect.ValueOf(greeter)
	if greeterValue.Kind() == reflect.Ptr {
		greeterValue = greeterValue.Elem()
	}
	return greeterValue.FieldByName("Greeter").Interface().(Greeter)
}

func (fb FacadeBuilder) buildArgDesc(method reflect.Value, ptrTab map[string]ReqTypeInfo) string {
	methodType := method.Type()
	numIn := methodType.NumIn()
	argsDesc := ""
	for i := 0; i < numIn; i++ {
		inType := methodType.In(i)
		inTypeString := inType.String()

		isPtr := false
		if inType.Kind() == reflect.Ptr {
			inType = inType.Elem()
			isPtr = true
		}

		argName := ""
		if inType.String() == "context.Context" {
			argName = "ctx"
		} else if inType.Kind() == reflect.Struct {
			argName = "req"
		} else {
			argName = "arg" + strconv.Itoa(i)
		}
		ptrTab[argName] = ReqTypeInfo{isPtr: isPtr, reqType: inType}

		split := ", "
		if i == numIn-1 {
			split = ""
		}
		argsDesc += fmt.Sprintf("%s %s%s", argName, inTypeString, split)
	}
	return argsDesc
}

func (fb FacadeBuilder) buildOutArgsDesc(method reflect.Value, ptrTab map[string]ReqTypeInfo) string {
	methodType := method.Type()
	numOut := methodType.NumOut()
	//限制出参只能有一个
	if numOut > 1 {
		panic("Only one value can be returned")
	}

	outArgsDesc := ""
	for i := 0; i < numOut; i++ {
		outType := methodType.Out(i)
		outTypeString := outType.String()

		isPtr := false
		if outType.Kind() == reflect.Ptr {
			outType = outType.Elem()
			isPtr = true
		}

		argName := ""
		if outType.String() == "context.Context" {
			argName = "ctx"
		} else if outType.Kind() == reflect.Struct {
			argName = "resp"
		} else {
			argName = "ret" + strconv.Itoa(i)
		}
		ptrTab[argName] = ReqTypeInfo{isPtr: isPtr, reqType: outType}

		split := ", "
		if i == numOut-1 {
			split = ""
		}
		outArgsDesc += fmt.Sprintf("%s%s", outTypeString, split)
	}
	return outArgsDesc
}

func (fb FacadeBuilder) buildQueryParamsDesc(httpType string, ptrTab map[string]ReqTypeInfo) string {
	reqObj, ok := ptrTab["req"]

	queryParamsDesc := ""
	queryBodyDesc := ""
	fullUrlDesc := "fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)"
	if ok {
		reqType := reqObj.reqType
		numField := reqType.NumField()
		assemblyQueryDesc(numField, reqType, reqObj, &queryParamsDesc, &queryBodyDesc)
	}

	if queryBodyDesc == "" && httpType == http.MethodPost {
		queryBodyDesc = "requestBody := \"\"\n"
	}
	if queryBodyDesc != "" {
		fullUrlDesc += "\n\tfullUrl += \"|\" + requestBody\n"
	}

	return queryParamsDesc + queryBodyDesc + fullUrlDesc
}


func (fb FacadeBuilder) buildRequestBodyLogDesc(httpType string) string {
	if httpType == http.MethodPost{
		return "requestBody"
	}
	return "\"\""
}

func assemblyQueryDesc(numField int, reqType reflect.Type, reqObj ReqTypeInfo, queryParamsDesc, queryBodyDesc *string) {
	for i := 0; i < numField; i++ {
		field := reqType.Field(i)
		fieldName := field.Name

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if _, ok := BasicTypeConverter[fieldType.String()]; ok {
			//基本类型
			*queryParamsDesc += fmt.Sprintf("\tqueryParams[\"%s\"] = %s.%s\n", str.LcFirst(fieldName), "req", fieldName)
		} else if reqType.Kind() == reflect.Struct {
			if *queryBodyDesc != "" {
				panic("req body more than one")
			}
			//body
			*queryBodyDesc = fmt.Sprintf("\trequestBody := json.ToJsonIgnoreError(%s.%s)\n", "req", fieldName)
		}
	}
	if reqObj.isPtr {
		*queryParamsDesc = "\tif req != nil {\n" + *queryParamsDesc + "\n\t}\n"
	}
}

func (fb FacadeBuilder) buildRequestDesc(httpType string, ptrTab map[string]ReqTypeInfo) string {
	headerOptionsDesc := ""
	requestDesc := ""

	respObj, ok := ptrTab["resp"]
	if !ok {
		panic("no resp respObj")
	}

	respObjPtrFlag := "*"
	if respObj.isPtr {
		respObjPtrFlag = ""
	}

	ctx, ok := ptrTab["ctx"]
	if ok {
		ctxPtrFlag := ""
		if ctx.isPtr {
			ctxPtrFlag = "*"
		}
		headerOptionsDesc += fmt.Sprintf("\theaderOptions, err := util.BuildHeaderOptions(%s%s)\n", ctxPtrFlag, "ctx")
		headerOptionsDesc += fmt.Sprint("\tif err != nil{\n")
		headerOptionsDesc += fmt.Sprint("\trespVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))\n")
		headerOptionsDesc += fmt.Sprintf("\treturn %srespVo\n", respObjPtrFlag)
		headerOptionsDesc += fmt.Sprint("\t}\n")
	}

	httpRequestMethodArgsDesc := ""
	if httpType == http.MethodPost {
		httpRequestMethodArgsDesc += ", requestBody"
	}
	if ok {
		httpRequestMethodArgsDesc += ", headerOptions..."
	}
	httpMethodType := httpType[0:1] + strings.ToLower(httpType[1:len(httpType)])
	requestDesc = fmt.Sprintf("\trespBody, respStatusCode, err := http.%s(reqUrl, queryParams%s)\n", httpMethodType, httpRequestMethodArgsDesc)
	return headerOptionsDesc + requestDesc
}

func (fb FacadeBuilder) buildReturnDesc(ptrTab map[string]ReqTypeInfo) string {
	respObj, ok := ptrTab["resp"]
	if !ok {
		panic("no resp respObj")
	}
	respObjPtrFlag := "*"
	if respObj.isPtr {
		respObjPtrFlag = ""
	}

	return fmt.Sprintf("return %srespVo", respObjPtrFlag)
}

package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type Node interface {
	NodeName() string
	Children() []Node
	AddChild(n Node)
}

type DataObject struct {
	Name       string
	Fields     []*Field
	ChildNodes []Node
}

type Constructor struct {
	FromInterface *Interface
	Method        *Method
}

type Interface struct {
	Name        string
	Methods     []*Method
	Constructor *Constructor
	ChildNodes  []Node
}

func (do *DataObject) NodeName() string {
	return do.Name
}

func (do *DataObject) Children() []Node {
	return do.ChildNodes
}

func (do *DataObject) AddChild(n Node) {
	do.ChildNodes = append(do.ChildNodes, n)
	return
}

func (do *DataObject) HasTimeType() (r bool) {
	for _, f := range do.Fields {
		if f.Type == "time.Time" {
			return true
		}
	}
	return false
}

func (do *DataObject) HasArrayType() (r bool) {
	for _, f := range do.Fields {
		if f.IsArray {
			return true
		}
	}
	return false
}

func (do *DataObject) HasMapType() (r bool) {
	for _, f := range do.Fields {
		if f.IsMap {
			return true
		}
	}
	return false
}

func (inf *Interface) NodeName() string {
	return inf.Name
}

func (inf *Interface) Children() []Node {
	return inf.ChildNodes
}

func (inf *Interface) AddChild(n Node) {
	inf.ChildNodes = append(inf.ChildNodes, n)
	return
}

type Method struct {
	Name                    string
	Params                  []*Field
	Results                 []*Field
	ConstructorForInterface *Interface
}

func (m *Method) ResultsForJavascriptFunction(prefix string) (r string) {
	rs := []string{}
	for _, r := range m.Results {
		rs = append(rs, prefix+"."+strings.Title(r.Name))
	}
	r = strings.Join(rs, ", ")
	return
}

func (m *Method) ParamsForJavascriptFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, p.Name)
	}
	r = strings.Join(ps, ", ")
	return
}

// In objective-c, method name will be the name of the first parameters
func (m *Method) ParamsForObjcFunctionWithCallback(apiprefix, interfaceName string) (r string) {
	ps := m.paramsForObj()

	successBlockParams := ""
	if m.ConstructorForInterface != nil {
		interName := m.ConstructorForInterface.Name
		successBlockParams = apiprefix + interName + "* " + interName
	} else {
		suffix := "results"
		if len(m.Results) == 1 {
			suffix = "error"
		}
		successBlockParams = m.ResultsForObjcFunction(interfaceName) + suffix
	}

	if len(m.Params) == 0 {
		ps[0] = ps[0] + ":(void (^)(" + successBlockParams + "))successBlock"
	} else {
		ps = append(ps, "success:(void (^)("+successBlockParams+"))successBlock")
	}
	ps = append(ps, "failure:(void (^)(NSError *error))failureBlock")

	r = strings.Join(ps, " ")
	return r
}

func (m *Method) paramsForObj() []string {
	mname := lowerFirstLetter(m.Name)
	if len(m.Params) == 0 {
		return []string{mname}
	}

	ps := []string{}
	for i, p := range m.Params {
		op := p.ToLanguageField("objc")
		name := op.Name

		if i == 0 {
			name = mname
		}

		ps = append(ps, name+":("+op.FullObjcTypeName()+")"+op.Name)
	}

	return ps
}

// see Method#ParamsForObjcFunctionWithCallback
func (m *Method) ParamsForObjcFunction() (r string) {
	ps := m.paramsForObj()

	r = strings.Join(ps, " ")
	return
}

func (m *Method) ParamsForJavaFunction() (r string) {
	if len(m.Params) == 0 {
		r = ""
		return
	}

	ps := []string{}
	for _, p := range m.Params {
		op := p.ToLanguageField("java")
		ps = append(ps, op.FullJavaTypeName()+" "+op.Name)
	}
	r = strings.Join(ps, ",")
	return
}

func (m *Method) ObjcReturnResultsOrOnlyOne() (r string) {
	if len(m.Results) == 1 {
		r = "results." + ObjcSnake(m.Results[0].Name)
		return
	}
	return "results"
}

func (m *Method) JavaReturnResultsOrOnlyOne() (r string) {
	if len(m.Results) == 1 {
		r = "results.get" + strings.Title(m.Results[0].Name) + "()"
		return
	}
	return "results"
}

func (m *Method) ResultsForObjcFunction(interfaceName string) (r string) {
	if len(m.Results) > 1 {
		r = m.Results[0].Prefix + interfaceName + m.Name + "Results *"
		return
	}
	if len(m.Results) == 0 {
		panic("method " + m.Name + "returned zero values")
	}
	r = m.Results[0].ToLanguageField("objc").Type
	return
}

func lowerFirstLetter(name string) string {
	if len(name) == 1 {
		return strings.ToLower(name)
	}
	return strings.ToLower(string(name[0])) + name[1:]
}

func (m *Method) ResultsForJavaFunction(interfaceName string) (r string) {
	if len(m.Results) > 1 {
		r = m.Name + "Results"
		return
	}
	if len(m.Results) == 0 {
		panic("method " + m.Name + "returned zero values")
	}
	r = m.Results[0].ToLanguageField("java").Type
	return
}

func (m *Method) ParamsForGoServerFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, "p.Params."+strings.Title(p.Name))
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ParamsForGoServerConstructorFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, "p.This."+strings.Title(p.Name))
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ResultsForGoServerFunction(prefix string) (r string) {
	rs := []string{}
	for _, r := range m.Results {
		rs = append(rs, prefix+"."+strings.Title(r.Name))
	}
	r = strings.Join(rs, ", ")
	return
}

func Snake(name string) (r string) {
	// first two letters are upcase like URL, MD5, HTTPRequest etc, keep it as it is.
	if len(name) >= 2 {
		if unicode.IsUpper([]rune(name)[0]) && (unicode.IsUpper([]rune(name)[1]) || unicode.IsNumber([]rune(name)[1])) {
			return name
		}
	}

	r = strings.ToLower(name[:1]) + name[1:]
	return
}

func ObjcSnake(name string) (r string) {

	// if name start with "new", "alloc", "copy", "mutableCopy", make the first letter upper case.
	r = Snake(name)
	if strings.Index(r, "new") == 0 ||
		strings.Index(r, "alloc") == 0 ||
		strings.Index(r, "copy") == 0 ||
		strings.Index(r, "mutableCopy") == 0 {
		r = strings.ToUpper(name[:1]) + name[1:]
	}
	return
}

func (m *Method) ParamsForGoClientFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, Snake(p.Name)+" "+p.FullGoTypeName())
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ResultsForGoClientFunction() (r string) {
	rs := []string{}
	for _, r := range m.Results {
		rs = append(rs, Snake(r.Name)+" "+r.FullGoTypeName())
	}
	r = strings.Join(rs, ", ")
	return
}

func (m *Method) ParamsForJson() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, `"`+strings.Title(p.Name)+`": `+p.Name)
	}
	r = strings.Join(ps, ", ")
	r = "{ " + r + " }"
	return
}

type Field struct {
	IsArray                     bool
	IsMap                       bool
	Type                        string
	Name                        string
	Star                        bool
	ImportName                  string
	PropertyAnnotation          string
	SetPropertyConvertFormatter string
	GetPropertyConvertFormatter string
	Primitive                   bool
	ConstructorType             string
	PkgName                     string
	Prefix                      string
	MapSpec                     [2]string
}

func (f Field) IsError() bool {
	return f.Type == "error"
}

func (f Field) FullGoTypeName() (r string) {
	if f.IsArray {
		r = r + "[]"
	}
	if f.IsMap {
		r = r + "map"
	}
	if f.Star {
		r = r + "*"
	}
	if f.ImportName != "" {
		r = r + f.ImportName + "."
	}
	r = r + f.Type
	return
}

func (f Field) FullObjcTypeName() (r string) {
	if f.IsArray {
		return "NSArray *"
	}
	if f.IsMap {
		return "NSDictionary *"
	}
	if f.Primitive {
		r = f.Type
		return
	}
	r = f.Prefix + f.Type
	return
}

func (f Field) FullJavaTypeName() (r string) {
	if f.IsArray {
		return "ArrayList<" + f.Type + ">"
	}
	if f.IsMap {
		return "Map"
	}
	if f.Primitive {
		r = f.Type
		return
	}
	r = f.Prefix + f.Type
	return
}

func (f Field) SetPropertyFromObjcDict(key string) (r string) {
	val := "[dict valueForKey:@\"" + key + "\"]"
	if len(strings.Split(f.SetPropertyConvertFormatter, "%s")) == 3 {
		r = fmt.Sprintf(f.SetPropertyConvertFormatter, f.Prefix+strings.Title(f.PkgName), val)
		return
	}
	r = fmt.Sprintf(f.SetPropertyConvertFormatter, val)
	return
}

func (f Field) SetPropertyObjc() (r string) {
	r = "[dict valueForKey:@\"" + strings.Title(f.Name) + "\"]"
	return
}

func (f Field) GetPropertyToObjcDict(key string) (r string) {
	if len(strings.Split(f.GetPropertyConvertFormatter, "%s")) == 3 {
		r = fmt.Sprintf(f.GetPropertyConvertFormatter, f.Prefix+strings.Title(f.PkgName), key)
		return
	}

	r = fmt.Sprintf(f.GetPropertyConvertFormatter, key)
	return
}

func (f Field) GetPropertyObjc() (r string) {
	r = "self." + ObjcSnake(f.Name)
	return
}

func findDefiniationNode(t string, apiset *APISet) (r Node) {
	for _, do := range apiset.DataObjects {
		if t == do.Name {
			return do
		}
	}
	for _, inf := range apiset.Interfaces {
		if t == inf.Name {
			return inf
		}
	}
	return
}

func (f *Field) Update(apiset *APISet, parentNode Node) {
	f.PkgName = apiset.Name
	f.Prefix = apiset.Prefix
	n := findDefiniationNode(f.Type, apiset)
	f.Primitive = true
	if n != nil {
		f.ImportName = apiset.Name
		f.Primitive = false
		parentNode.AddChild(n)
	}
}

func (f Field) ToLanguageField(language string) (r Field) {
	languageMap, ok := TypeMapping[language]
	if !ok {
		panic(language + " not supported.")
	}

	r.Name = f.Name
	r.IsArray = f.IsArray
	r.IsMap = f.IsMap
	r.Star = f.Star
	r.ImportName = f.ImportName
	r.Primitive = f.Primitive
	r.PkgName = f.PkgName
	r.Prefix = f.Prefix
	t := languageMap.TypeOf(f)
	r.Type = t.Type
	r.PropertyAnnotation = t.PropertyAnnotation
	r.SetPropertyConvertFormatter = t.SetPropertyConvertFormatter
	r.GetPropertyConvertFormatter = t.GetPropertyConvertFormatter
	r.ConstructorType = t.ConstructorType
	r.MapSpec = r.MapSpec
	return
}

package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

type Walker struct {
	fset              *token.FileSet
	APISet            *APISet
	currentName       string
	currentInterface  *Interface
	currentDataObject *DataObject
	currentMethod     *Method
	currentFieldList  *[]Field
}

type APISet struct {
	Name          string
	Prefix        string
	ImplPkg       string
	ServerImports []string
	Interfaces    []*Interface
	DataObjects   []*DataObject
}

func Parse(dir string, prefix string) (r *APISet) {
	fset := token.NewFileSet()

	astPkgs, _ := parser.ParseDir(fset, dir, nil, 0)

	var foundPkg *ast.Package
	for _, astPkg := range astPkgs {
		foundPkg = astPkg
	}

	w := &Walker{fset: fset, APISet: &APISet{}}

	// ast.Print(fset, foundPkg)
	ast.Walk(w, foundPkg)
	w.APISet.Prefix = prefix
	updateConstructors(w.APISet)
	updateFields(w.APISet)
	sortDataObjects(w.APISet)
	sortInterfaces(w.APISet)
	r = w.APISet

	return
}

// Visit implements the ast.Visitor interface.
func (w *Walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Package:
		w.APISet.Name = n.Name
	case *ast.TypeSpec:
		w.currentName = n.Name.Name
	case *ast.StructType:
		w.currentDataObject = &DataObject{Name: w.currentName}
		w.currentInterface = nil
		w.currentFieldList = nil
		w.APISet.DataObjects = append(w.APISet.DataObjects, w.currentDataObject)
	case *ast.InterfaceType:
		w.currentInterface = &Interface{Name: w.currentName}
		w.currentDataObject = nil
		w.currentFieldList = nil
		w.APISet.Interfaces = append(w.APISet.Interfaces, w.currentInterface)
	case *ast.FuncType:
		if w.currentInterface != nil {
			w.currentMethod = &Method{Name: w.currentName}
			for _, param := range n.Params.List {
				fs := parseField(param)
				w.currentMethod.Params = append(w.currentMethod.Params, fs...)
			}
			markIfStreamAndRemoveLastParam(w.currentMethod)
			if n.Results != nil {
				for _, result := range n.Results.List {
					fs := parseField(result)
					w.currentMethod.Results = append(w.currentMethod.Results, fs...)
				}
			}
			if len(w.currentMethod.Results) == 0 {
				die("method " + w.currentName + " must have return values in forms like this: (entry *Entry, err error)")
			} else {
				if w.currentMethod.Results[len(w.currentMethod.Results)-1].Type != "error" {
					die("method " + w.currentMethod.Name + " of " + w.currentInterface.Name + "'s must additionally return 'err error'")
				}
			}
			w.currentInterface.Methods = append(w.currentInterface.Methods, w.currentMethod)
			w.currentDataObject = nil
			w.currentFieldList = nil
		}
	case *ast.Field:
		if isImportedField(n) {
			break
		}

		if w.currentInterface != nil && len(n.Names) > 0 {
			w.currentName = n.Names[0].Name
		}

		if w.currentDataObject != nil && len(n.Names) > 0 {
			fs := parseField(n)
			w.currentDataObject.Fields = append(w.currentDataObject.Fields, fs...)
		}
	}

	return w
}

/*
If in the api packages, there are two interfaces like this:

	type Manager interface {
		login(name, password string) (api API, err error)
	}

	type API interface {
		getProfiile(name string) (profile string, err error)
	}

Then Manager will be the constructor of API, for Manage has functions that would return API.
This would be reflect in client packages, like in objective-c, accessing instances of API
would have to go through Manager. This is could be used as a kind of authentication mechanism.
*/
func updateConstructors(apiset *APISet) {
	for _, inf := range apiset.Interfaces {
		for _, inftarget := range apiset.Interfaces {
			for _, m := range inftarget.Methods {
				for _, f := range m.Results {
					if f.Type == inf.Name {
						m.ConstructorForInterface = inf
						if inf.Constructor != nil {
							die(inf.Name + "'s constructor already is " + inf.Constructor.Method.Name + ", can only one constructor exists for one service")
						}
						inf.Constructor = &Constructor{inftarget, m}
					}
				}
			}
		}
	}
}

func markIfStreamAndRemoveLastParam(m *Method) {
	if len(m.Params) == 0 {
		return
	}

	count := 0
	for _, p := range m.Params {
		if p.Type == "io.Reader" {
			count++
		}
	}

	if count == 1 && m.Params[len(m.Params)-1].Type == "io.Reader" {
		m.IsStreamInput = true
		m.StreamParamName = m.Params[len(m.Params)-1].Name
		m.Params = m.Params[0 : len(m.Params)-1]
	}

	if count > 1 {
		die(fmt.Sprintf("%s can NOT have more than one io.Reader parameters and it must be the last parameter", m.Name))
	}
	return

}

func die(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func updateFields(apiset *APISet) {
	for _, inf := range apiset.Interfaces {
		for _, m := range inf.Methods {
			for _, p := range m.Params {
				p.Update(apiset, inf)
			}
			for _, p := range m.Results {
				p.Update(apiset, inf)
			}
		}
	}
	for _, do := range apiset.DataObjects {
		for _, f := range do.Fields {
			f.Update(apiset, do)
		}
	}
}

type byDepthDataObjects struct {
	DataObjects []*DataObject
}

func (b byDepthDataObjects) Less(i, j int) bool {
	return depth(b.DataObjects[i], 20) < depth(b.DataObjects[j], 20)
}

func (b byDepthDataObjects) Len() int { return len(b.DataObjects) }

func (b byDepthDataObjects) Swap(i, j int) {
	b.DataObjects[i], b.DataObjects[j] = b.DataObjects[j], b.DataObjects[i]
}

type byDepthInterfaces struct {
	Interfaces []*Interface
}

func (b byDepthInterfaces) Less(i, j int) bool {
	return depth(b.Interfaces[i], 20) < depth(b.Interfaces[j], 20)
}

func (b byDepthInterfaces) Len() int { return len(b.Interfaces) }

func (b byDepthInterfaces) Swap(i, j int) {
	b.Interfaces[i], b.Interfaces[j] = b.Interfaces[j], b.Interfaces[i]
}

func depth(n Node, maxdepth int) (r int) {
	if maxdepth < 0 {
		panic("loop dependency: " + n.NodeName())
	}

	if len(n.Children()) == 0 {
		return 1
	}

	max := 0
	maxdepth = maxdepth - 1
	for _, c := range n.Children() {
		// ignore self reference
		if n.NodeName() == c.NodeName() {
			continue
		}
		d := depth(c, maxdepth)
		if d > max {
			max = d
		}
	}
	r = max + 1
	return
}

func sortDataObjects(apiset *APISet) {
	sort.Sort(byDepthDataObjects{apiset.DataObjects})
}

func sortInterfaces(apiset *APISet) {
	sort.Sort(byDepthInterfaces{apiset.Interfaces})
}

func parseField(n *ast.Field) (r []*Field) {
	for _, id := range n.Names {
		f := &Field{Name: id.Name}

		switch nt := n.Type.(type) {
		case *ast.Ident:
			f.Type = nt.Name
		case *ast.SelectorExpr:
			f.Type = nt.X.(*ast.Ident).Name + "." + nt.Sel.Name
		case *ast.StarExpr:
			f.Star = true
			switch xt := nt.X.(type) {
			case *ast.Ident:
				f.Type = xt.Name
			case *ast.SelectorExpr:
				f.Type = xt.X.(*ast.Ident).Name + "." + xt.Sel.Name
			}
		case *ast.ArrayType:
			switch vt := nt.Elt.(type) {
			case *ast.Ident:
				f.Type = nt.Elt.(*ast.Ident).Name
			case *ast.StarExpr:
				f.Star = true
				f.Type = vt.X.(*ast.Ident).Name
			case *ast.SelectorExpr:
				f.Type = vt.X.(*ast.Ident).Name + "." + vt.Sel.Name
			}
			f.IsArray = true
		case *ast.MapType:
			switch vt := nt.Key.(type) {
			case *ast.Ident:
				f.Type = `[` + vt.Name + `]`
				f.MapSpec[0] = vt.Name
			case *ast.SelectorExpr:
				f.Type = `[` + vt.X.(*ast.Ident).Name + "." + vt.Sel.Name + `]`
				f.MapSpec[0] = vt.X.(*ast.Ident).Name + "." + vt.Sel.Name
			}

			switch vt := nt.Value.(type) {
			case *ast.Ident:
				f.Type = f.Type + vt.Name
				f.MapSpec[1] = vt.Name
			case *ast.SelectorExpr:
				f.Type = f.Type + vt.X.(*ast.Ident).Name + "." + vt.Sel.Name
				f.MapSpec[1] = vt.X.(*ast.Ident).Name + "." + vt.Sel.Name
			}

			f.IsMap = true
		}

		r = append(r, f)
	}

	return
}

func isImportedField(f *ast.Field) bool {
	walker := &importedFieldWalker{}
	ast.Walk(walker, f.Type)
	return walker.containsImported
}

type importedFieldWalker struct {
	containsImported bool
}

const supportedSelectorExpr = "template.HTML,template.HTMLAttr,time.Time,govalidations.Validated,io.Reader,"

func (i *importedFieldWalker) Visit(node ast.Node) ast.Visitor {
	switch ftype := node.(type) {
	case *ast.SelectorExpr:
		if i.containsImported {
			break
		}

		name := ftype.X.(*ast.Ident).Name + "." + ftype.Sel.Name + ","
		i.containsImported = !strings.Contains(supportedSelectorExpr, name)
	}

	return i
}

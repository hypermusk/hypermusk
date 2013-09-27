package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
)

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

type Walker struct {
	fset              *token.FileSet
	APISet            *APISet
	currentName       string
	currentInterface  *Interface
	currentDataObject *DataObject
	currentMethod     *Method
	currentFieldList  *[]Field
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

			if n.Results != nil {
				for _, result := range n.Results.List {
					fs := parseField(result)
					w.currentMethod.Results = append(w.currentMethod.Results, fs...)
				}
			}
			if len(w.currentMethod.Results) == 0 {
				die("method " + w.currentName + " must have return values or must have return value names like (entry *Entry, err error)")
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
		if w.currentInterface != nil && len(n.Names) > 0 {
			w.currentName = n.Names[0].Name
			// fmt.Println(w.currentName)
		}

		if w.currentDataObject != nil && len(n.Names) > 0 {
			fs := parseField(n)
			w.currentDataObject.Fields = append(w.currentDataObject.Fields, fs...)
		}
	}
	return w
}

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
	// println(n.NodeName())

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
		// println("=> ", maxdepth, ":", c.NodeName())
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
		f := &Field{}
		r = append(r, f)
		f.Name = id.Name
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
			var tname *ast.Ident
			st, isstar := nt.Elt.(*ast.StarExpr)
			if isstar {
				f.Star = true
				tname = st.X.(*ast.Ident)
			} else {
				tname = nt.Elt.(*ast.Ident)
			}
			f.Type = tname.Name
			f.IsArray = true

		case *ast.MapType:
			k := nt.Key.(*ast.Ident)
			v := nt.Value.(*ast.Ident)
			f.Type = `[` + k.Name + `]` + v.Name
			f.IsMap = true

		}
	}
	return

}

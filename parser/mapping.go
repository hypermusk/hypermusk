package parser

type LanguageType struct {
	Type                        string
	PropertyAnnotation          string
	GetPropertyConvertFormatter string
	SetPropertyConvertFormatter string
	ConstructorType             string
}

var TypeMapping = map[string]TypeMap{
	"objc": TypeMap{
		KnownMapping: map[string]LanguageType{
			"string":                  {"NSString *", "(nonatomic, strong)", "%s", "%s", "NSString"},
			"int64":                   {"NSNumber *", "(nonatomic, strong)", "%s", "%s", "NSNumber"},
			"int32":                   {"NSNumber *", "(nonatomic, strong)", "%s", "%s", "NSNumber"},
			"int":                     {"NSNumber *", "(nonatomic, strong)", "%s", "%s", "NSNumber"},
			"float64":                 {"NSNumber *", "(nonatomic, strong)", "%s", "%s", "NSNumber"},
			"float32":                 {"NSNumber *", "(nonatomic, strong)", "%s", "%s", "NSNumber"},
			"float":                   {"NSNumber *", "(nonatomic, strong)", "%s", "%s", "NSNumber"},
			"bool":                    {"BOOL", "(nonatomic, assign)", "[NSNumber numberWithBool:%s]", "[%s boolValue]", "BOOL"},
			"error":                   {"NSError *", "(nonatomic, strong)", "%s", "%s", "NSError"},
			"template.HTML":           {"NSString *", "(nonatomic, strong)", "%s", "%s", "NSString"},
			"template.HTMLAttr":       {"NSString *", "(nonatomic, strong)", "%s", "%s", "NSString"},
			"time.Time":               {"NSDate *", "(nonatomic, strong)", "[%s stringFromDate:%s]", "[%s dateFromString:%s]", "NSDate"},
			"govalidations.Validated": {"Validated *", "(nonatomic, strong)", "%s", "%s", "Validated"},
		},
		UnknownFunc: func(f Field) (r LanguageType) {
			r = LanguageType{f.Type + " *", "(nonatomic, strong)", "%s", "%s", f.Type}
			return
		},
	},
	"java": TypeMap{
		KnownMapping: map[string]LanguageType{
			"string":                  {"String", "", "%s", "%s", "String"},
			"int64":                   {"int", "", "%s", "%s", "int"},
			"int32":                   {"int", "", "%s", "%s", "int"},
			"int":                     {"int", "", "%s", "%s", "int"},
			"float64":                 {"float", "", "%s", "%s", "float"},
			"float32":                 {"float", "", "%s", "%s", "float"},
			"float":                   {"float", "", "%s", "%s", "float"},
			"bool":                    {"boolean", "", "%s", "%s", "boolean"},
			"error":                   {"RemoteError", "", "%s", "%s", "RemoteError"},
			"template.HTML":           {"String", "", "%s", "%s", "String"},
			"template.HTMLAttr":       {"String", "", "%s", "%s", "String"},
			"time.Time":               {"Date", "", "Utils.stringFromDate(%s)]", "Utils.dateFromString(%s)", "Date"},
			"govalidations.Validated": {"Validated", "", "%s", "%s", "Validated"},
		},
		UnknownFunc: func(f Field) (r LanguageType) {
			r = LanguageType{f.Type + "", "", "%s", "%s", f.Type}
			return
		},
	},
}

type TypeMap struct {
	KnownMapping map[string]LanguageType
	UnknownFunc  func(f Field) (r LanguageType)
}

func (tm TypeMap) TypeOf(f Field) (r LanguageType) {
	t, ok := tm.KnownMapping[f.Type]
	if ok {
		r = t
		return
	}
	r = tm.UnknownFunc(f)
	return
}

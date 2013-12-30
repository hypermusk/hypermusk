package parser

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FieldSuite struct {
	Field
}

var _ = Suite(&FieldSuite{})

func (d *FieldSuite) TestFullObjcTypeName(c *C) {
	d.Field = Field{IsMap: true,
		Type:                        "[string]string *",
		Name:                        "Map",
		Star:                        false,
		PropertyAnnotation:          "(nonatomic, strong)",
		SetPropertyConvertFormatter: "%s",
		GetPropertyConvertFormatter: "%s",
		Primitive:                   true,
		ConstructorType:             "[string]string",
		PkgName:                     "api",
	}
}

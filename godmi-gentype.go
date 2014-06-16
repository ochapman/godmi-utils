/*
* generate function and string method
 */
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	//	"reflect"
)

type StructType struct {
	Name string
	Type string
}

type StructTypes []StructType

const prepart = `
package main

type SMBIOSStructureType byte

type SMBIOSStructureHandle uint16

type InfoCommon struct {
	Type   SMBIOSStructureType
	Length byte
	Handle SMBIOSStructureHandle
}
`

func gen(file string, typename string) (StructTypes, error) {
	fset := token.NewFileSet()
	template, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	src := prepart + string(template)
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	sts := make(StructTypes, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			for _, s := range x.Specs {
				switch k := s.(type) {
				case *ast.TypeSpec:
					if k.Name.Name == typename {
						//fmt.Println(k)
						switch t := k.Type.(type) {
						case *ast.StructType:
							for _, l := range t.Fields.List {
								var st StructType
								st.Name = fmt.Sprint(l.Names)
								st.Type = fmt.Sprint(l.Type)
								//fmt.Println(st.Name, st.Type)
								sts = append(sts, st)
							}
						}
					}
				}
			}
		}
		return true
	})
	return sts, nil
}

// split Capital
// BigBrotherIsWatchingYou -> Big Brother Is Watching You
func splitCap(s string) string {
	var last int
	var str string
	for i, ss := range s {
		if (i > 1 && (ss >= 'A' && ss <= 'Z') && (s[i+1] >= 'a' && s[i+1] <= 'z')) || i == len(s)-2 {
			if i == len(s)-2 {
				if last == 0 {
					str = s
				} else {
					str += " " + s[last:]
				}
				break
			}

			if str == "" {
				str = s[last:i]
			} else {
				str += " " + s[last:i]
			}
			last = i
		}
	}
	return str
}

func (s StructTypes) NewString(name string) string {
	return "NewString"
}

func (s StructTypes) TypeString(name string) string {
	lowername := strings.ToLower(name)
	for i, ss := range s {
		s[i].Name = strings.TrimFunc(ss.Name, func(r rune) bool {
			if r == '[' || r == ']' {
				return true
			}
			return false
		})
	}
	fmt.Printf("func (%c %s) String() string {\n", lowername[0], name)
	fmt.Printf("return fmt.Sprintf(\"%s:\\n\\t\\t\"+\n", splitCap(name))
	for i, ss := range s {
		if ss.Type == "InfoCommon" {
			continue
		}
		var fm string
		if ss.Type == "string" {
			fm = "%s"
		} else {
			fm = "%d"
		}
		fmt.Printf("\"%s: %s", splitCap(ss.Name), fm)
		if i != len(s)-1 {
			fmt.Printf("\\n\\t\\t\"+\n")
		} else {
			fmt.Printf("\\n\",\n")
		}
	}
	for _, ss := range s {
		if ss.Type == "InfoCommon" {
			continue
		}
		fmt.Printf("%c.%s,\n", lowername[0], ss.Name)
	}
	fmt.Println(")\n}")
	return "TypeString"
}

func main() {
	template := flag.String("template", "", "template for typename")
	typename := flag.String("typename", "", "typename in template")
	flag.Parse()
	if *template == "" || *typename == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	sts, err := gen(*template, *typename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen() failed: %s", err)
		os.Exit(-1)
	}
	sts.TypeString("PortableBattery")
}

/*
* godmi-gentype.go
* generate function and string method
*
* Chapman Ou <ochapman@ochapman.cn>
* 2014-06-17
 */
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
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

type infoCommon struct {
	SMType   SMBIOSStructureType
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

func (s StructTypes) DMIHeaderMethod(name string) string {
	nstring := bytes.NewBuffer([]byte{})
	fmt.Fprintf(nstring, "func (h DMIHeader) %s() *%s {\n"+
		"data := h.data\n"+
		"return &%s{\n",
		name, name, name)
	// Field
	var offset int
	for _, ss := range s {
		if ss.Type == "infoCommon" {
			offset += 4
			continue
		}
		switch ss.Type {
		case "byte":
			fmt.Fprintf(nstring, "%s: data[0x%02X],\n", ss.Name, offset)
			offset += 1
		case "uint16":
			fmt.Fprintf(nstring, "%s: u16(data[0x%02X:0x%02X]),\n", ss.Name, offset, offset+2)
			offset += 2
		case "uint32":
			fmt.Fprintf(nstring, "%s: u32(data[0x%02X:0x%02X]),\n", ss.Name, offset, offset+4)
			offset += 4
		case "uint64":
			fmt.Fprintf(nstring, "%s: u64(data[0x%02X:0x%02X]),\n", ss.Name, offset, offset+8)
			offset += 8
		case "string":
			fmt.Fprintf(nstring, "%s: h.FieldString(int(data[0x%02X])),\n", ss.Name, offset)
			offset += 1
		default:
			fmt.Fprintf(nstring, "%s: %s(data[0x%02X]),\n", ss.Name, ss.Type, offset)
			offset += 1
		}
	}
	// End of func
	fmt.Fprintf(nstring, "}\n}\n")
	fmtnstring, err := format.Source(nstring.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "format.Source() failed: %s", err)
		fmt.Println(nstring)
	}
	return string(fmtnstring)
}

func (s StructTypes) TypeString(name string) string {
	tstring := bytes.NewBuffer([]byte{})
	v := strings.ToLower(name)[0]
	for i, ss := range s {
		s[i].Name = strings.TrimFunc(ss.Name, func(r rune) bool {
			if r == '[' || r == ']' {
				return true
			}
			return false
		})
	}
	fmt.Fprintf(tstring, "func (%c %s) String() string {\n", v, name)
	fmt.Fprintf(tstring, "return fmt.Sprintf(\"%s\\n\"+\n", splitCap(name))
	for i, ss := range s {
		if ss.Type == "infoCommon" {
			continue
		}
		var fm string
		if ss.Type == "byte" || ss.Type == "uint16" || ss.Type == "uint32" || ss.Type == "uint64" {
			fm = "%d"
		} else {
			fm = "%s"
		}
		fmt.Fprintf(tstring, "\"\\t%s: %s", splitCap(ss.Name), fm)
		if i != len(s)-1 {
			fmt.Fprintf(tstring, "\\n\"+\n")
		} else {
			fmt.Fprintf(tstring, "\",\n")
		}
	}
	for i, ss := range s {
		if ss.Type == "infoCommon" {
			continue
		}
		if i != len(s)-1 {
			fmt.Fprintf(tstring, "%c.%s,\n", v, ss.Name)
		} else {
			fmt.Fprintf(tstring, "%c.%s)\n", v, ss.Name)
		}
	}
	fmt.Fprintf(tstring, "}\n")
	fmttstring, err := format.Source(tstring.Bytes())
	if err != nil {
		return "format Error: " + err.Error()
	}
	return string(fmttstring)
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
	s := sts.TypeString(*typename)
	fmt.Println(s)
	n := sts.DMIHeaderMethod(*typename)
	fmt.Println(n)
}

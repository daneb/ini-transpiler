package main

import (
	"github.com/alecthomas/repr"
	"io/ioutil"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var iniLexer = lexer.MustSimple([]lexer.Rule{
	{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`, nil},
	{`String`, `"(?:\\.|[^"])*"`, nil},
	{`Float`, `\d+(?:\.\d+)?`, nil},
	{`Punct`, `[][=]`, nil},
	{"comment", `[#;][^\n]*`, nil},
	{"whitespace", `\s+`, nil},
})

type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value *Value `@@`
}

type Value struct {
	Pos    lexer.Position
	String *string  ` @String`
	Number *float64 `| @Float`
}

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

var parser = participle.MustBuild(&INI{},
	participle.Lexer(iniLexer),
	participle.Unquote("String"),
)

func main() {
	ini := &INI{}
	content, _ := ioutil.ReadFile("sample.ini")

	err := parser.ParseString("", string(content), ini)
	if err != nil {
		panic(err)
	}
	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
}

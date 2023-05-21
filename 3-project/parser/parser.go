package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"os"
	"strings"
)

var (
	dfParser *participle.Parser[DepFile] = participle.MustBuild[DepFile](participle.Lexer(dfLexer))
	dfLexer                              = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},
		{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*([.][a-zA-Z][a-zA-Z0-9]*)?`},
		{"Punct", `<-`},
		{"EOL", `[;]`},
	})
)

type DepFile struct {
	Rules []*Rule `(@@)+ `
}

type Rule struct {
	Object string   `@Ident "<-"`
	Deps   []string `@Ident+ ";"`
}

func (df *DepFile) String() string {
	var res string
	for _, r := range df.Rules {
		res += r.String() + "\n"
	}
	return res
}

func (r *Rule) String() string {
	var deps string
	for _, d := range r.Deps {
		deps += d + " "

	}
	return r.Object + " <- " + strings.Join(r.Deps, " ")
}

func Parse(s string) (*DepFile, error) {
	return dfParser.ParseString("", s)

}

func ParseFile(file string) (*DepFile, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	ast, err := dfParser.Parse("", r)
	r.Close()
	return ast, err
}

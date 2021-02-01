package translation

import (
	"github.com/dangxia/google-translate-api/ctx"
	"strings"
	"testing"
)

func TestTranslation(t *testing.T) {
	ctx, _ := ctx.NewContext()

	input := `It’s not enough to know what you need to do—you’ve got to deliver, too. This one is all about keeping your nose to the grindstone and doing good work`

	ts := NewTranslation(ctx.DefaultSourceLang(), ctx.DefaultTargetLang(), input, ctx)
	s, _ := ts.Get()
	println(s)

	s, _ = ts.Get()
	println(strings.Join(s, ","))
}

func TestAnalyzeResult(t *testing.T) {
	input := `)]}'

511
[["wrb.fr","MkEWBc","[[null,null,null,[[[0,[[[null,5]\n]\n,[true]\n]\n]\n]\n,5]\n]\n,[[[null,\"Nín hǎo\",null,null,null,[[\"您好\",[\"您好\",\"你好\"]\n]\n]\n]\n]\n,\"zh-CN\",1,\"en\",[\"Hello\",\"en\",\"zh-CN\",true]\n]\n,null,[\"Hello!\",null,null,null,null,[[[\"感叹词\",[[\"你好!\",null,[\"Hello!\",\"Hi!\",\"Hallo!\"]\n,1,true]\n,[\"喂!\",null,[\"Hey!\",\"Hello!\"]\n,2,true]\n]\n,\"zh-CN\",\"en\"]\n]\n,2]\n,null,null,\"en\",1]\n]\n",null,null,null,"generic"]
,["di",21]
,["af.httprm",20,"8829770370350956435",64]
]
26
[["e",4,null,null,574]
]
`

	list, err := Analyze(input)
	if err != nil {
		t.Fatal(err)
	}

	if "您好,您好,你好" != strings.Join(list, ",") {
		t.Fatal()
	}
}

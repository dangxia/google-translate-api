package translation

import (
	"github.com/dangxia/google-translate-api/ctx"
	"strings"
	"testing"
)

func TestTranslation(t *testing.T) {
	ctx, _ := ctx.NewContext()

	input := `The DataStream API supports different runtime execution modes from which you can choose depending on the requirements of your use case and the characteristics of your job.

There is the “classic” execution behavior of the DataStream API, which we call STREAMING execution mode. This should be used for unbounded jobs that require continuous incremental processing and are expected to stay online indefinitely.

Additionally, there is a batch-style execution mode that we call BATCH execution mode. This executes jobs in a way that is more reminiscent of batch processing frameworks such as MapReduce. This should be used for bounded jobs for which you have a known fixed input and which do not run continuously.

Apache Flink’s unified approach to stream and batch processing means that a DataStream application executed over bounded input will produce the same final results regardless of the configured execution mode. It is important to note what final means here: a job executing in STREAMING mode might produce incremental updates (think upserts in a database) while a BATCH job would only produce one final result at the end. The final result will be the same if interpreted correctly but the way to get there can be different.

By enabling BATCH execution, we allow Flink to apply additional optimizations that we can only do when we know that our input is bounded. For example, different join/aggregation strategies can be used, in addition to a different shuffle implementation that allows more efficient task scheduling and failure recovery behavior. We will go into some of the details of the execution behavior below.`

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

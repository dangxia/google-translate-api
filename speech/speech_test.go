package speech

import (
	"github.com/dangxia/google-translate-api/ctx"
	"os"
	"path/filepath"
	"testing"
)

func TestTranslation(t *testing.T) {
	ctx, _ := ctx.NewContext()

	input := `The DataStream API supports different runtime execution modes from which you can choose depending on the requirements of your use case and the characteristics of your job.

There is the “classic” execution behavior of the DataStream API, which we call STREAMING execution mode. This should be used for unbounded jobs that require continuous incremental processing and are expected to stay online indefinitely.

Additionally, there is a batch-style execution mode that we call BATCH execution mode. This executes jobs in a way that is more reminiscent of batch processing frameworks such as MapReduce. This should be used for bounded jobs for which you have a known fixed input and which do not run continuously.

Apache Flink’s unified approach to stream and batch processing means that a DataStream application executed over bounded input will produce the same final results regardless of the configured execution mode. It is important to note what final means here: a job executing in STREAMING mode might produce incremental updates (think upserts in a database) while a BATCH job would only produce one final result at the end. The final result will be the same if interpreted correctly but the way to get there can be different.

By enabling BATCH execution, we allow Flink to apply additional optimizations that we can only do when we know that our input is bounded. For example, different join/aggregation strategies can be used, in addition to a different shuffle implementation that allows more efficient task scheduling and failure recovery behavior. We will go into some of the details of the execution behavior below.`
	ts := NewSpeech(ctx, input, ctx.DefaultSourceLang(), false)

	pwd, _ := os.Getwd()
	err := ts.Save(filepath.Join(pwd, "test.mp3"))
	if err != nil {
		t.Fatal(err)
	}
}

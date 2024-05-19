package hw06pipelineexecution

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func collectResults(out Out) []interface{} {
	var results []interface{}
	for v := range out {
		results = append(results, v)
	}
	return results
}

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("Empty Input", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		close(in) // Закрываем входной канал сразу

		out := ExecutePipeline(in, done, stages...)
		results := collectResults(out)

		if len(results) != 0 {
			t.Errorf("Expected empty results, got %v", results)
		}
	})

	t.Run("No Stages", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)

		go func() {
			defer close(in)
			for i := 0; i < 5; i++ {
				in <- i
			}
		}()

		out := ExecutePipeline(in, done)
		results := collectResults(out)
		expected := []interface{}{0, 1, 2, 3, 4}

		if len(results) != len(expected) {
			t.Errorf("Expected length %v, got %v", len(expected), len(results))
		}

		for i, v := range results {
			if v != expected[i] {
				t.Errorf("Expected value %v, got %v", expected[i], v)
			}
		}
	})

	t.Run("Pipeline With Strings", func(t *testing.T) {
		stringStage := g("String Processor", func(v interface{}) interface{} { return v.(string) + "_processed" })

		in := make(Bi)
		done := make(Bi)

		go func() {
			defer close(in)
			for _, s := range []string{"a", "b", "c"} {
				in <- s
			}
		}()

		out := ExecutePipeline(in, done, stringStage)
		results := collectResults(out)
		expected := []interface{}{"a_processed", "b_processed", "c_processed"}

		if len(results) != len(expected) {
			t.Errorf("Expected length %v, got %v", len(expected), len(results))
		}

		for i, v := range results {
			if v != expected[i] {
				t.Errorf("Expected value %v, got %v", expected[i], v)
			}
		}
	})
}

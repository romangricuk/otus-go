package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {

	for _, stage := range stages {
		stageChan := make(Bi)

		go func(in In, stageChan Bi) {
			defer close(stageChan)

			for {
				select {
				case <-done:
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					if data != nil {
						if !ok {
							return
						}
						select {
						case <-done:
							return
						case stageChan <- data:
						}
					}
				}
			}
		}(in, stageChan)

		in = stage(stageChan)
	}

	return in
}

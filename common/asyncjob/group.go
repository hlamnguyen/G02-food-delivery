package asyncjob

import (
	"context"
	"sync"
)

type group struct {
	isParallel bool
	jobs       []Job
	wg         *sync.WaitGroup
}

func NewGroup(isParallel bool, jobs ...Job) *group {
	g := &group{
		isParallel: isParallel,
		jobs:       jobs,
		wg:         new(sync.WaitGroup),
	}

	return g
}

//func (g *group) Run2(ctx context.Context) error {
//	errChan := make(chan error, len(g.jobs))
//
//	for i, _ := range g.jobs {
//		errChan <- g.runJob(ctx, g.jobs[i])
//	}
//
//	var err error
//
//	for i := 1; i <= len(g.jobs); i++ {
//		if v := <-errChan; v != nil {
//			err = v
//		}
//	}
//
//	return err
//}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isParallel {
			// Dont do that
			//go func() {
			//	// "i" is pointer to i of the loop outside
			//	errChan <- g.runJob(ctx, g.jobs[i]) // capture "i"
			//	g.wg.Done()
			//}()

			// Do this instead
			go func(aj Job) {
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}

		errChan <- g.runJob(ctx, g.jobs[i])
		g.wg.Done()
	}

	var err error

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}

	g.wg.Wait()
	return err
}

// Retry if needed
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}

	return nil
}

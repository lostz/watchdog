package machine

import "context"

type State func(*Runtime)

type item struct {
	state State
	ctx   context.Context
}

//Runtime ...
type Runtime struct {
	chi chan *item
	ctx context.Context
}

func (r *Runtime) NextState(ctx context.Context, state State) {
	r.chi <- &item{state, ctx}
}

func (r *Runtime) Context() context.Context {
	return r.ctx
}

func (r *Runtime) loop(start State) {
	var cancel context.CancelFunc
	r.ctx, cancel = context.WithCancel(r.ctx)

	go func() {
		var localItem *item
		ok := true

		rootCtx := r.ctx

		for ok {
			select {
			case localItem, ok = <-r.chi:
				if ok {
					if localItem.state != nil {
						r.ctx = localItem.ctx
						localItem.state(r)
					} else {
						//we are going to cancel if state sends a nil state as a next state
						//this implies that state reaches the end.
						cancel()
					}
				}
			case _, ok = <-rootCtx.Done():
			}
		}

		defer close(r.chi)
	}()
	start(r)
}

func Run(ctx context.Context, initialState State) *Runtime {
	runtime := Runtime{
		chi: make(chan *item, 1),
		ctx: ctx,
	}

	runtime.loop(initialState)

	return &runtime
}

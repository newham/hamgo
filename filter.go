package hamgo

type Filter interface {
	DoFilter(ctx Context) bool
}

type filter struct {
	Filter
	handler func(ctx Context) bool
}

func (f *filter) DoFilter(ctx Context) bool {
	return f.handler(ctx)
}

func newFilter(handler func(ctx Context) bool) Filter {
	return &filter{handler: handler}
}

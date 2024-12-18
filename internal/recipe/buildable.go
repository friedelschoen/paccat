package recipe

type Buildable interface {
	Build(*Context) (string, error)
	HasOutput() bool
}

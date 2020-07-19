package limter

type Limiter interface {
	Allow() bool
	Revert()
}

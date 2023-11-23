package interfaces

type JobBuilder interface {
	Build() (err error)
}

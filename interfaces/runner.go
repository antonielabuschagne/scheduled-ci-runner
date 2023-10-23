package interfaces

type IJobRunner interface {
	RunBuild() (err error)
}

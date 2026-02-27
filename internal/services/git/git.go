package git

type Info struct {
	Branch      string
	CommitHash  string
	Remote      string
	HasUnstaged bool
	HasStaged   bool
}

type Service interface {
	Detect(path string) (bool, error)
	GetInfo(path string) (*Info, error)
	GetModifiedFiles(path string) ([]string, error)
}

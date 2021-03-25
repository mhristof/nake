package precommit

type Repos struct {
	Repos []Repo
}

type Repo struct {
	Repo  string
	Rev   string
	Hooks []Hook
}

type Hook struct {
	ID string
}

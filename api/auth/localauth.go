package auth

type LocalAuth struct{}

func NewLocalAuth() *LocalAuth {
	return &LocalAuth{}
}

func (a *LocalAuth) LocalAuth(email, name string) (string, error) {
	return "", nil
}

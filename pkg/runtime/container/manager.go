package container

type Container struct {
	Id   string
	Name string
}

func newContainerManager() (*Container, error) {
	return &Container{}, nil
}

func createContainer(name string) (*Container, error) {
	return &Container{
		Id:   "123",
		Name: name,
	}, nil
}

package mongo

type mongo struct {
}

func New() *mongoClient {
	return &mongoClient{}
}

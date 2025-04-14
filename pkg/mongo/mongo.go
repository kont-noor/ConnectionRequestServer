package mongo

type mongoClient struct {
}

func New() *mongoClient {
	return &mongoClient{}
}

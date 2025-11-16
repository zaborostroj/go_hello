package warehouse

type ServiceConfig struct {
	KAFKA struct {
		Host    string
		Port    string
		Topic   string
		GroupId string
	}
}

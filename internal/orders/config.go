package orders

type ServiceConfig struct {
	APP struct {
		Host string
		Port string
	}
	DB struct {
		Prefix   string
		Username string
		Password string
		Host     string
		Port     string
		Dbname   string
	}
	KAFKA struct {
		Host    string
		Port    string
		Topic   string
		GroupId string
	}
}

package conf

type Config struct {
	Host string
	RoutinesNum uint64
	RoutineSize uint64
	Addr string
	WsAddr string
	Client int
}

func Init() *Config {
	return new(Config)
}
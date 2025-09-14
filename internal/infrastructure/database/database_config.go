package database

var (
	DB_NAME        = ""
	DB_HOST        = ""
	DB_PORT        = ""
	DB_USER        = ""
	DB_PASSWORD    = ""
	REDIS_ADDR     = ""
	REDIS_PORT     = ""
	REDIS_PASSWORD = ""
)

func SetDatabaseEnv(name, host, port, user, password string) {
	DB_NAME = name
	DB_HOST = host
	DB_PORT = port
	DB_USER = user
	DB_PASSWORD = password
}

func SetRedisEnv(address, port, pass string) {
	REDIS_ADDR = address
	REDIS_PORT = port
	REDIS_PASSWORD = pass
}

package config

type database struct {
	name        string `mapstructure:"DATABASE_NAME"`
	user        string `mapstructure:"DATABASE_USER"`
	password    string `mapstructure:"DATABASE_PASSWORD"`
	ssl_mode    string `mapstructure:"DATABASE_SSL_MODE"`
	port        string `mapstructure:"DATABASE_PORT"`
	sym_encrypt string `mapstructure:"DATABASE_SYM_ENCRYPT"`
}

type jwt struct {
	secret         string `mapstructure:"JWT_SECRET"`
	secret_temp    string `mapstructure:"JWT_SECRET_TEMP"`
	expiredIn      string `mapstructure:"JWT_EXPIREDIN"`
	expiredIn_temp string `mapstructure:"JWT_EXPIREDIN_TEMP"`
}

type rabbitMQ struct {
	user     string `mapstructure:"RABBITMQ_USER"`
	pass     string `mapstructure:"RABBITMQ_PASS"`
	host     string `mapstructure:"RABBITMQ_HOST"`
	port     string `mapstructure:"RABBITMQ_PORT"`
	vhost    string `mapstructure:"RABBITMQ_VHOST"`
	consumer string `mapstructure:"RABBITMQ_CONSUMER"`
	queue    string `mapstructure:"RABBITMQ_QUEUE"`
	exchange string `mapstructure:"RABBITMQ_EXCHANGE"`
	dlx      string `mapstructure:"RABBITMQ_DLX"`
	queues
	routingKey
}

type queues struct {
	queue string `mapstructure:"RABBITMQ_QUEUE_USER"`
}

type routingKey struct {
	userCreate       string `mapstructure:"ROUTINGKEY_USERCREATE"`
	userCreateReturn string `mapstructure:"ROUTINGKEY_USERCREATERETURN"`
}

// type grpc struct {
// 	port int `mapstructure:"GRPC_PORT"`
// }

type storage struct {
	localPath string `mapstructure:"LOCAL_STORAGE_PATH"`
	localFile string `mapstructure:"LOCAL"`
}

type authz struct {
	max_admin int    `mapstructure:"AUTHZ_MAX_ADMIN"`
	max_read  int    `mapstructure:"AUTHZ_MAX_READ"`
	admin_1   string `mapstructure:"AUTHZ_ADMIN_1"`
	read_1    string `mapstructure:"AUTHZ_READ_1"`
}

type security struct {
	aes_key  string `mapstructure:"SECURITY_AES_KEY"`
	hmac_key string `mapstructure:"SECURITY_HMAC_KEY"`
}

type password struct {
	redis string `mapstructure:"PASSWORD_REDIS"`
}

type concurrency struct {
	create_management string `mapstructure:"CONCURRENCY_CREATE_MANAGEMENT"`
}

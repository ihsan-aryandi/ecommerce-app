package config

type Config struct {
	ServerAddr string
	Database   *DatabaseConfig
	Redis      *RedisConfig
	Midtrans   *MidtransConfig
	RajaOngkir *RajaOngkirConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	DB                 string
	PoolSize           string
	MinIdleConnections string
}

type MidtransConfig struct {
	Host      string
	Endpoints *MidtransEndpointsConfig
}

type MidtransEndpointsConfig struct {
	Pay *RequestConfig
}

type RajaOngkirConfig struct {
	APIKey    string
	Host      string
	Endpoints *RajaOngkirEndpointsConfig
}

type RajaOngkirEndpointsConfig struct {
	CalculateShippingCost *RequestConfig
}

type RequestConfig struct {
	Method string
	Path   string
}

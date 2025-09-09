package config

type Config struct {
	ServerAddr string
	Database   *DatabaseConfig
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

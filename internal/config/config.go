package config

import (
	"os"
)

func Load() *Config {
	return &Config{
		ServerAddr: getEnv("HTTP_PORT", "8080"),
		Database: &DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "postgres"),
		},
		Midtrans: &MidtransConfig{
			Host: getEnv("MIDTRANS_HOST", ""),
		},
		RajaOngkir: &RajaOngkirConfig{
			APIKey: getEnv("RAJA_ONGKIR_API_KEY", ""),
			Host:   getEnv("RAJA_ONGKIR_HOST", ""),
			Endpoints: &RajaOngkirEndpointsConfig{
				CalculateShippingCost: &RequestConfig{
					Method: getEnv("RAJA_ONGKIR_METHOD_CALCULATE_SHIPPING_COST", ""),
					Path:   getEnv("RAJA_ONGKIR_PATH_CALCULATE_SHIPPING_COST", ""),
				},
			},
		},
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

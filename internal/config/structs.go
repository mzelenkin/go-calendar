package config

// Config это главная структура конфигурации
type Config struct {
	Listen Listen `config:"listen"`
	Log LoggingConfig `config:"log"`
}

// Listen настройки демона (для примера)
type Listen struct {
	Port int `config:"port"`
	Address string `config:"address"`
}

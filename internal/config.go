package internal

type Config struct {
	HTTP  HTTP  `yaml:"http"`
	SQL   SQL   `yaml:"sql"`
	Redis Redis `yaml:"redis"`
}

type HTTP struct {
	Addr string `yaml:"addr"`
}

type SQL struct {
	Driver  string `yaml:"driver"`
	URL     string `yaml:"url"`
	MaxOpen int    `yaml:"max_open"`
	MaxIdle int    `yaml:"max_idle"`
}

type Redis struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

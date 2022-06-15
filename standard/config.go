package standard

// Config 数据库配置
type Config struct {
	// MySQL/PostgresSQL/SQLite
	Driver          string `toml:"driver"`             // 数据库驱动
	Dsn             string `toml:"dsn"`                // 数据库源
	ConnMaxLifeTime int    `toml:"conn_max_life_time"` // 数据库最大连接时长
	MaxIdleConns    int    `toml:"max_idle_conns"`
	MaxOpenConns    int    `toml:"max_open_conns"`
	Slave           []struct {
		Dsn string `toml:"dsn"`
	} `toml:"slave"` // 数据库备用源
}

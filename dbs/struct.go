package dbs

type RedisDBStruct struct {
	Address  string
	Password string
}

type RedisClusterDBStruct struct {
	Addrs    []string
	Password string
}

type RedisRingDBStruct struct {
	Addrs    map[string]string
	Password string
}

type RedisConfig struct {
	Prefix  string
	Type    string
	Redis   RedisDBStruct
	Cluster RedisClusterDBStruct
	Ring    RedisRingDBStruct
}

// MySQLDSN 代表一个 mysql dsn 连接信息
type MySQLDSN struct {
	Name    string
	DSN     string
	Type    string
	SSHName string
}

// MySQLDB 代表一个 mysql 连接信息
type MySQLDB struct {
	Read     MySQLDSN
	Write    MySQLDSN
	Timezone string
	Region   string
	CoinType string
}

// AddDatabaseConfig 添加数据库配置
func AddDatabaseConfig(value *MySQLDB, configs map[string]MySQLDSN) {
	if value.Read.DSN != "" && value.Read.Name != "" {
		configs[value.Read.Name] = MySQLDSN{DSN: value.Read.DSN, SSHName: value.Read.SSHName, Type: value.Read.Type}
	}
	if value.Write.DSN != "" && value.Write.Name != "" {
		configs[value.Write.Name] = MySQLDSN{DSN: value.Write.DSN, SSHName: value.Read.SSHName, Type: value.Write.Type}
	}
}

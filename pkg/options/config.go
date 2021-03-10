package options

type Config struct {
	LogLevel       string
	HTTPPort       int
	GraphqlPort    int
	GRPCPort       int
	InfoPort       int
	PrometheusPort int
	// Postgres       postgres.Config
}

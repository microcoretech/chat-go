package configs

type Environment string

const (
	DevelopmentEnvironment Environment = "development"
	ProductionEnvironment  Environment = "production"
)

func (e Environment) String() string {
	return string(e)
}

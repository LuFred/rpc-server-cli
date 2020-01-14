package redis

type alias struct {
	Name         string
	DriverType   DriverType
	DB           Database
	MaxIdleConns int
	MaxOpenConns int
}

func addAliasWithCache(aliasName string, driverType DriverType, db Database) (*alias, error) {
	al := new(alias)
	al.Name = aliasName
	al.DriverType = driverType
	al.DB = db
	return al, nil
}

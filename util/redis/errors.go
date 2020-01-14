package redis

import (
	"errors"
)

var (
	ErrDbNotFound    = errors.New("Cache: The database to be closed does not exist")
	ErrAliasNotExist = errors.New("Cache: The specified alias does not exist")

	ErrDriverTypeNotFound = errors.New("Cache: Driver Type does not exist")

	ErrPoolDisconnected = errors.New("Cache: Connection pool has been disconnected")

	ErrNotHaveAvailableConn = errors.New("Cache: No connections available")

	ErrValueNotBeenNil = errors.New("Cache: Value cannot be empty")
	ErrKeyNotBeenNil   = errors.New("Cache: Key cannot be empty")

	ErrHashScan = errors.New("Cache: Wrong hash structure")

	ErrScanStructValue = errors.New("Cache: value must be non-nil pointer to a struct")
)

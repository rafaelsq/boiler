//go:generate ../../mock.sh storage
package iface

import "database/sql"

type Storage interface {
	SQL() *sql.DB
}

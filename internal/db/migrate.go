package db

type Migratable interface {
	Schema() string
}

func (c *Client) Migrate(migrationStructs ...Migratable) {
	for _, m := range migrationStructs {
		c.DB.MustExec(m.Schema())
	}
}

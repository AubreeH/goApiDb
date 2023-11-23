package driver

type Key struct {
	TableCatalog      string
	TableSchema       string
	TableName         string
	ColumnName        string
	ConstraintCatalog string
	ConstraintSchema  string
	ConstraintName    string
	OrdinalPosition   int
}

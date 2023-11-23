package driver

type TableDescription struct {
	Columns []Column
}

type Column struct {
	Name         string
	Type         string
	Nullable     string
	Key          string
	DefaultValue string
	Extra        string
}

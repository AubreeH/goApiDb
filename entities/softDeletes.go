package entities

type SoftDeletes struct {
	Deleted bool `json:"deleted" sql_name:"deleted" sql_type:"BOOLEAN" sql_nullable:"FALSE" sql_default:"FALSE" soft_deletes:"true"`
}

func (val SoftDeletes) OnDelete() (SoftDeletes, error) {
	val.Deleted = true
	return val, nil
}

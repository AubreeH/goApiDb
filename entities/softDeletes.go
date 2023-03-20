package entities

type SoftDeletes struct {
	Deleted bool `json:"deleted" db_type:"BOOLEAN" db_nullable:"TRUE" soft_deletes:"true"`
}

func (val SoftDeletes) OnDelete() (SoftDeletes, error) {
	val.Deleted = true
	return val, nil
}

func (_ SoftDeletes) GetPtrFunc(value *SoftDeletes) *bool {
	return &value.Deleted
}

package entities

type Dates struct {
	CreatedAt
	UpdatedAt
	DeletedAt
}

func (_ Dates) GetPtrFunc(value *Dates) []any {
	pointers := []any{
		value.UpdatedAt.GetPtrFunc(&value.UpdatedAt),
		value.CreatedAt.GetPtrFunc(&value.CreatedAt),
	}

	pointers = append(pointers, value.DeletedAt.GetPtrFunc(&value.DeletedAt)...)

	return pointers
}

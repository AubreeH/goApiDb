package entities

type Dates struct {
	CreatedAt
	UpdatedAt
	DeletedAt
}

func (_ Dates) GetPtrFunc(value *Dates) any {
	pointers := map[string]any{
		"updated_at": value.UpdatedAt.GetPtrFunc(&value.UpdatedAt),
		"created_at": value.CreatedAt.GetPtrFunc(&value.CreatedAt),
		"deleted_at": value.DeletedAt.GetPtrFunc(&value.DeletedAt),
	}

	return pointers
}

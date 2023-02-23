package entities

type IKey interface {
	string | int | int64
}

type ForeignKey[TModel any, TKey IKey] struct {
	Key   TKey
	model TModel
}

func (ForeignKey[TModel, TKey]) GetPtrFunc(val *ForeignKey[TModel, TKey]) *TKey {
	return &val.Key
}

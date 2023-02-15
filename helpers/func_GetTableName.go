package helpers

import "reflect"

func GetTableName(entity any) string {
	entityType := reflect.TypeOf(entity)
	return PascalToSnakeCase(entityType.Name()) + "s"
}

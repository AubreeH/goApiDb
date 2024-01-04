package structParsing

// FormatNullable takes the db_null struct tag string and determines if the column is nullable or not.
func FormatNullable(nullable string) string {
	if FormatBoolean(nullable) != 1 {
		return "NOT NULL"
	}

	return ""
}

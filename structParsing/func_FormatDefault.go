package structParsing

// FormatDefault takes the db_default struct tag string and determines the default values of the field.
func FormatDefault(def string) string {
	if def == "" {
		return ""
	}

	return "DEFAULT " + def
}

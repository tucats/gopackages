package cli

// GetParameter returns the ith parameter string parsed, or an
// empty string if not found.
func GetParameter(i int) string {
	if i < GetParameterCount() {
		return Parameters[i]
	}
	return ""
}

// GetParameterCount returns the number of parameters processed.
func GetParameterCount() int {
	return len(Parameters)
}

// WasFound reports if an entry in the grammar was found on
// the processed command line.
func WasFound(grammar Options, name string) bool {
	for _, entry := range grammar {

		if entry.OptionType == Subcommand && entry.Found {
			return WasFound(entry.Value.(Options), name)
		}
		if entry.Found && name == entry.LongName {
			return true
		}
	}
	return false
}

// GetInteger returns the value of a named integer from the
// parsed grammar, or a zero if not found. The boolean return
// value confirms if the value was specified on the command line.
func GetInteger(grammar Options, name string) (int, bool) {

	for _, entry := range grammar {

		if entry.OptionType == Subcommand && entry.Found {
			return GetInteger(entry.Value.(Options), name)
		}
		if entry.Found && entry.OptionType == IntType && name == entry.LongName {
			return entry.Value.(int), true
		}
	}
	return 0, false
}

// GetBool returns the value of a named integer from the
// parsed grammar, or a zero if not found.
func GetBool(grammar Options, name string) bool {

	for _, entry := range grammar {

		if entry.OptionType == Subcommand && entry.Found {
			return GetBool(entry.Value.(Options), name)
		}
		if entry.Found && (entry.OptionType == BooleanType || entry.OptionType == BooleanValueType) && name == entry.LongName {
			return entry.Value.(bool)
		}
	}
	return false
}

// GetString returns the value of a named integer from the
// parsed grammar, or a zero if not found. The second return value
// indicates if the value was explicitly specified.
func GetString(grammar Options, name string) (string, bool) {

	for _, entry := range grammar {

		if entry.OptionType == Subcommand && entry.Found {
			return GetString(entry.Value.(Options), name)
		}
		if entry.Found && entry.OptionType == StringType && name == entry.LongName {
			return entry.Value.(string), true
		}
	}
	return "", false
}

// GetStringList returns the array of strings that are the value of
// the named item. If the item is not found, an empty array is returned.
// The second value in the result indicates of the option was explicitly
// specified in the command line.
func GetStringList(grammar Options, name string) ([]string, bool) {
	for _, entry := range grammar {

		if entry.OptionType == Subcommand && entry.Found {
			return GetStringList(entry.Value.(Options), name)
		}
		if entry.Found && entry.OptionType == StringListType && name == entry.LongName {
			return entry.Value.([]string), true
		}
	}
	return make([]string, 0), false
}

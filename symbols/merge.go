package symbols

import "github.com/tucats/gopackages/app-cli/ui"

// Merge merges the contents of a table into the current table.
func (s *SymbolTable) Merge(st *SymbolTable) {

	ui.Debug("+++ Merging symbols from %s", st.Name)
	for k, v := range st.Symbols {

		// Is it a struct? If so we may need to merge to it...
		switch vv := v.(type) {

		case map[string]interface{}:

			// Does the old struct already exist in the compiler table?
			old, found := s.Get(k)
			if found {

				// Is the existing value also a struct?
				switch oldmap := old.(type) {
				case map[string]interface{}:

					// Copy the values into the existing map
					for newkeyword, newvalue := range vv {
						oldmap[newkeyword] = newvalue
						ui.Debug("    adding %v to old map at %s", newvalue, newkeyword)
					}
					// Rewrite the map back to the bytecode.
					s.SetAlways(k, oldmap)

				default:
					ui.Debug("    overwriting duplicate key %s with %v", k, old)
					s.SetAlways(k, v)
				}

			} else {
				ui.Debug("    creating new map %s with %v", k, v)
				s.SetAlways(k, v)
			}
		default:
			ui.Debug("    copying entry %s with %v", k, v)
			s.SetAlways(k, vv)
		}
	}

	// Do it again with the constants

	ui.Debug("+++ Merging constants from  %s", st.Name)
	for k, v := range st.Constants {

		// Is it a struct? If so we may need to merge to it...
		switch vv := v.(type) {

		case map[string]interface{}:

			// Does the old struct already exist in the compiler table?
			old, found := s.Get(k)
			if found {

				// Is the existing value also a struct?
				switch oldmap := old.(type) {
				case map[string]interface{}:

					// Copy the values into the existing map
					for newkeyword, newvalue := range vv {
						oldmap[newkeyword] = newvalue
						ui.Debug("    adding %v to old map at %s", newvalue, newkeyword)
					}
					// Rewrite the map back to the bytecode.
					s.SetConstant(k, oldmap)

				default:
					ui.Debug("    overwriting duplicate key %s with %v", k, old)
					s.SetConstant(k, v)
				}

			} else {
				ui.Debug("    creating new map %s with %v", k, v)
				s.SetConstant(k, v)
			}
		default:
			ui.Debug("    copying entry %s with %v", k, v)
			s.SetConstant(k, vv)
		}
	}
}

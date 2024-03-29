package builtins

// MaxDeepCopyDepth specifies the maximum depth that a recursive
// copy will go before failing. Setting this too small will
// prevent complex structures from copying correctly. Setting it
// too large can result in excessive memory consumption.
const MaxDeepCopyDepth = 100

// DeepCopy makes a deep copy of an Ego data type. It should be called with the
// maximum nesting depth permitted (i.e. array index->array->array...). Because
// it calls itself recursively, this is used to determine when to give up and
// stop traversing nested data. The default is MaxDeepCopyDepth.
func DeepCopy(source interface{}, depth int) interface{} {
	if depth < 0 {
		return nil
	}

	switch v := source.(type) {
	case bool:
		return v

	case byte:
		return v

	case int32:
		return v

	case int:
		return v

	case int64:
		return v

	case string:
		return v

	case float32:
		return v

	case float64:
		return v

	case []interface{}:
		r := make([]interface{}, 0)

		for _, d := range v {
			r = append(r, DeepCopy(d, depth-1))
		}

		return r

	default:
		return v
	}
}

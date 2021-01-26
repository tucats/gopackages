package datatypes

type EgoMap struct {
	data    map[interface{}]interface{}
	keyType int
}

func NewMap(keyType int) *EgoMap {

	m := &EgoMap{
		data:    map[interface{}]interface{}{},
		keyType: keyType,
	}
	return m
}

func (m *EgoMap) Get(key interface{}) (interface{}, bool) {
	if IsType(key, m.keyType) {
		v, found := m.data[key]
		return v, found
	}
	return nil, false
}

func (m *EgoMap) Set(key interface{}, value interface{}) bool {
	if !IsType(key, m.keyType) {
		return false
	}
	_, found := m.data[key]
	m.data[key] = value
	return found
}

func (m *EgoMap) Keys() []interface{} {
	r := []interface{}{}
	for k := range m.data {
		r = append(r, k)
	}
	return r
}

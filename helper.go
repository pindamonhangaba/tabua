package tabua

// Sq is a default Query implementation
type Sq struct {
	S string
	A []interface{}
}

// SQL implements Query
func (q Sq) SQL() string {
	return q.S
}

// Args implements Query
func (q Sq) Args() []interface{} {
	return q.A
}

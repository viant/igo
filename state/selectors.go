package state

//Selectors represents selectors
type Selectors []*Selector

//IDs returns selector IDs
func (s Selectors) IDs() []string {
	var result = make([]string, len(s))
	if len(s) == 0 {
		return result
	}
	for i, sel := range s {
		result[i] = sel.ID
	}
	return result
}

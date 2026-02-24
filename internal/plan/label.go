package plan

func (s *Scope) ensureLabelID(name string) uint32 {
	if id, ok := s.labels[name]; ok {
		return id
	}
	id := *s.nextLabel
	*s.nextLabel = *s.nextLabel + 1
	s.labels[name] = id
	return id
}

func (s *Scope) labelID(name string) (uint32, bool) {
	id, ok := s.labels[name]
	return id, ok
}

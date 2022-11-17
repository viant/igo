package option

//Option represents a planner option
type Option interface{}

//Options option slice
type Options []Option

//Tracker returns a tracker option or nil
func (o Options) Tracker() *Tracker {
	for _, item := range o {
		if t, ok := item.(*Tracker); ok {
			return t
		}
	}
	return nil
}

//StmtListener returns a stmt listener option or nil
func (o Options) StmtListener() StmtListener {
	for _, item := range o {
		if t, ok := item.(StmtListener); ok {
			return t
		}
	}
	return nil
}

//ExprListener returns an expr listener option or nil
func (o Options) ExprListener() ExprListener {
	for _, item := range o {
		if t, ok := item.(ExprListener); ok {
			return t
		}
	}
	return nil
}

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

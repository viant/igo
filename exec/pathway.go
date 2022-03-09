package exec

import "reflect"

//Pathway defines Selector leaf more complex pathway
type Pathway uint8

const (
	//PathwayUndefined represents non Selector component
	PathwayUndefined = Pathway(0)
	//PathwayDirect a path Selector has no pointer, nor slices nor map, precomputed Offset can be used for the fastest access
	PathwayDirect = Pathway(1)
	//PathwayRef a path Selector has at least on pointer or slice
	PathwayRef = Pathway(2)
)

//NewPathway returns a pathway for bitmas
func NewPathway(p Pathway) Pathway {
	if p&PathwayRef != 0 {
		return PathwayRef
	}
	return PathwayDirect
}

//IsDirect returns true if path can be access without pointer/slices
func (p Pathway) IsDirect() bool {
	return p == PathwayDirect
}

//IsRef returns true if path has ancestor with pointer/slice
func (p Pathway) IsRef() bool {
	return p == PathwayRef
}

//SelectorPathway sets  Selector pathway to pre computer fixed Offset as long it's not pointer,slice or map
func SelectorPathway(s *Selector) Pathway {
	if len(s.Ancestors) == 0 {
		return PathwayDirect
	}
	for _, f := range s.Ancestors {
		if f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice {
			return PathwayRef
		}
	}
	return PathwayDirect
}

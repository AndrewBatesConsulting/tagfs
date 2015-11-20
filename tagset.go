package tagfs

import (
	"code.google.com/p/go-uuid/uuid"
)

type tagset map[string]uuid.UUID

func loadTagSet(filename string) (tagset, error) {
	return nil, nil
}

func (t tagset) save(filename string) error {
	return nil
}

func intersection(sets ...tagset) tagset {
	intersection := make(tagset)

tag:
	for tag, value := range sets[0] {
		for _, set := range sets[1:] {
			if _, found := set[tag]; !found {
				continue tag
			}
		}
		intersection[tag] = value
	}
	return intersection
}

func union(sets ...tagset) tagset {
	union := make(tagset)
	for _, set := range sets {
		for tag, value := range set {
			union[tag] = value
		}
	}
	return union
}

func difference(sets ...tagset) tagset {
	difference := sets[0]
	for _, set := range sets[1:] {
		union := union(difference, set)
		intersection := intersection(difference, set)
		for tag, _ := range intersection {
			delete(union, tag)
		}
		difference = union
	}

	return difference
}

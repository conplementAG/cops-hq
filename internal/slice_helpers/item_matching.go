package slice_helpers

import "slices"

func FindItemsToAddAndRemove(current, expected []string) (add, remove []string) {
	for _, item := range current {
		if !slices.Contains(expected, item) {
			remove = append(remove, item)
		}
	}

	for _, item := range expected {
		if !slices.Contains(current, item) {
			add = append(add, item)
		}
	}

	return add, remove
}

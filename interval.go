package patchwork

import "fmt"

// interval is an inclusive range between two points, starting at zero
type interval struct {
	from, to int64
}

func (i interval) Contains(c interval) bool {
	return i.from <= c.from && i.to >= c.to
}

func (i interval) Intersects(c interval) bool {
	if i.from > c.from {
		return c.to >= i.from
	}
	return i.to >= c.from
}

func (i interval) Merge(c interval) interval {
	if c.to > i.to {
		i.to = c.to
	}
	if c.from < i.from {
		i.from = c.from
	}
	return i
}

func (i interval) String() string {
	return fmt.Sprintf("%d-%d", i.from, i.to)
}

// intervalSet is a collection of ordered, non-overlapping intervals
type intervalSet struct {
	intervals []interval
}

// Add adds an interval to the set, either merging it with an existing intersecting interval
// or inserting it in the correct
func (is *intervalSet) Add(i interval) {
	for n := 0; n < len(is.intervals); n++ {
		if is.intervals[n].Intersects(i) {
			is.intervals[n] = is.intervals[n].Merge(i)
			if len(is.intervals) >= n+2 && is.intervals[n].Intersects(is.intervals[n+1]) {
				is.intervals[n] = is.intervals[n].Merge(is.intervals[n+1])
				is.intervals = append(is.intervals[0:n+1], is.intervals[n+2:]...)
			}
			return
		}
	}

	is.intervals = append(is.intervals, i)
	for n := 1; n < len(is.intervals); n++ {
		if is.intervals[n].Intersects(is.intervals[n-1]) {
			is.intervals[n-1], is.intervals[n] = is.intervals[n], is.intervals[n-1]
		}
	}
}

func (is *intervalSet) Contains(c interval) bool {
	for _, i := range is.intervals {
		if i.Contains(c) {
			return true
		}
	}
	return false
}

func (is *intervalSet) Intervals() []interval {
	return is.intervals
}

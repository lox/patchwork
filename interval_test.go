package patchwork

import "testing"

func TestIntervalContainsAndIntersects(t *testing.T) {
	tests := []struct {
		i1, i2     interval
		contains   bool
		intersects bool
	}{
		{interval{0, 5}, interval{3, 10}, false, true},
		{interval{2, 5}, interval{5, 10}, false, true},
		{interval{0, 5}, interval{0, 5}, true, true},
		{interval{0, 5}, interval{1, 5}, true, true},
		{interval{0, 5}, interval{0, 4}, true, true},
		{interval{0, 5}, interval{6, 10}, false, false},
		{interval{0, 50}, interval{45, 52}, false, true},
	}

	for idx, st := range tests {
		if v := st.i1.Contains(st.i2); v != st.contains {
			t.Fatalf("Test #%d expected contains %#v, got %#v", idx, st.contains, v)
		}
		if v := st.i1.Intersects(st.i2); v != st.intersects {
			t.Fatalf("Test #%d expected intersects %#v, got %#v", idx, st.intersects, v)
		}
	}
}

func TestIntervalSet(t *testing.T) {
	ts := intervalSet{}
	ts.Add(interval{0, 2})
	ts.Add(interval{2, 5})
	ts.Add(interval{8, 10})
	ts.Add(interval{15, 25})
	ts.Add(interval{10, 12})
	ts.Add(interval{2, 4})
	ts.Add(interval{5, 8})
	ts.Add(interval{25, 26})
	ts.Add(interval{27, 28})
	ts.Add(interval{26, 27})
	ts.Add(interval{28, 29})

	if !ts.Contains(interval{0, 10}) {
		t.Fatalf("Expected to find 0-10")
	}

	if !ts.Contains(interval{6, 12}) {
		t.Fatalf("Expected to find 6-12")
	}

	if ts.Contains(interval{0, 25}) {
		t.Fatalf("Expected to NOT find 0-25")
	}

	if !ts.Contains(interval{15, 29}) {
		t.Fatalf("Expected to find 15-29")
	}
}

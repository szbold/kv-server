package datatypes

import "testing"

func TestStringResponse(t *testing.T) {
	var s KvString
	s = "example"
	got := s.Response()
	want := []byte("+7\r\nexample\r\n")
	failed := false

	for i := range got {
		if got[i] != want[i] {
			t.Logf("Results differs from want at position %v", i)
			failed = true
		}
	}

	if failed {
		t.Errorf("Want: %v Got: %v", want, got)
	}
}

func TestIntResponse(t *testing.T) {
	var s KvInt
	s = 123
	got := s.Response()
	want := []byte(":123\r\n")
	failed := false

	for i := range got {
		if got[i] != want[i] {
			t.Logf("Results differs from want at position %v", i)
			failed = true
		}
	}

	if failed {
		t.Errorf("Want: %v Got: %v", want, got)
	}
}

func TestListResponse(t *testing.T) {
	var l KvList
	l = []string{"1", "23", "456"}
	got := l.Response()
	want := []byte("*3\r\n$1\r\n1\r\n$2\r\n23\r\n$3\r\n456\r\n")
	failed := false

	for i := range got {
		if got[i] != want[i] {
			t.Logf("Results differs from want at position %v", i)
			failed = true
		}
	}

	if failed {
		t.Errorf("Want: %v Got: %v", want, got)
	}
}

func TestErrorResponse(t *testing.T) {
	var s KvError
	s = KvError{"Error"}
	got := s.Response()
	want := []byte("-5\r\nError\r\n")
	failed := false

	for i := range got {
		if got[i] != want[i] {
			t.Logf("Results differs from want at position %v", i)
			failed = true
		}
	}

	if failed {
		t.Errorf("Want: %v Got: %v", want, got)
	}
}

func TestSetResponse(t *testing.T) {
	var l KvSet
	l = NewKvSet()
  l.Insert("1")
  l.Insert("23")
  l.Insert("456")

	got := l.Response()
	want := []byte("*3\r\n$1\r\n1\r\n$2\r\n23\r\n$3\r\n456\r\n")
	failed := false

	for i := range got {
		if got[i] != want[i] {
			t.Logf("Results differs from want at position %v", i)
			failed = true
		}
	}

	if failed {
		t.Errorf("Want: %v Got: %v", want, got)
	}
}


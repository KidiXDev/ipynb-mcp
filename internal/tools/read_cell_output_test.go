package tools

import "testing"

func TestPaginateText(t *testing.T) {
	t.Parallel()

	chunk, start, end, total, hasMore, err := paginateText("abcdefghij", 2, 4)
	if err != nil {
		t.Fatalf("paginateText returned error: %v", err)
	}
	if chunk != "cdef" {
		t.Fatalf("chunk: got %q want %q", chunk, "cdef")
	}
	if start != 2 || end != 6 || total != 10 {
		t.Fatalf("unexpected pagination values: start=%d end=%d total=%d", start, end, total)
	}
	if !hasMore {
		t.Fatalf("expected hasMore to be true")
	}
}

func TestPaginateTextOffsetOutOfRange(t *testing.T) {
	t.Parallel()

	_, _, _, _, _, err := paginateText("abc", 4, 2)
	if err == nil {
		t.Fatalf("expected out-of-range error")
	}
}

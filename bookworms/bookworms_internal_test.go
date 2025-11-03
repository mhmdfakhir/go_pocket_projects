package main

import (
	"reflect"
	"testing"
)

var (
	handmaidsTale = Book{Author: "Margaret Atwood", Title: "The Handmaid's Tale"}
	oryxAndCrake  = Book{Author: "Margaret Atwood", Title: "Oryx and Crake"}
	theBellJar    = Book{Author: "Sylvia Plath", Title: "The Bell Jar"}
	janeEyre      = Book{Author: "Charlotte Brontë", Title: "Jane Eyre"}
	villette      = Book{Author: "Charlotte Brontë", Title: "Villette"}
	ilPrincipe    = Book{Author: "Niccolò Machiavelli", Title: "Il Principe"}
)
type testCase struct {
	bookwormsFile string
	want          []Bookworm
	wantErr       bool
}

func TestLoadBookworms_Success(t *testing.T) {
	tests := map[string]struct {
		bookwormsFile string
		want          []Bookworm
		wantErr       bool
	}{
		"file exists": {
			bookwormsFile: "testdata/bookworms.json",
			want: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			wantErr: false,
		},
		"file doesn't exist": {
			bookwormsFile: "testdata/no_file_here.json",
			want:          nil,
			wantErr:       true,
		},
		"invalid JSON": {bookwormsFile: "testdata/invalid.json",
			want:    nil,
			wantErr: true,
		},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := loadBookworms(testCase.bookwormsFile)
			if err != nil && !testCase.wantErr {
				t.Fatalf("expected an error %s, got none", err.Error())
			}
			if err == nil && testCase.wantErr {
				t.Fatalf("expected no error, got one %s", err.Error())
			}
			if !equalBookworms(t, got, testCase.want) {
				t.Fatalf("different result: got %v, expected %v",
					got, testCase.want)
			}

			// equivalent to above use of equalBookworms
			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("different result: got %v, expected %v", got, testCase.want)
			}
		})
	}
}

// equalBookworms is a helper to test the equality of two lists of Bookworms.
func equalBookworms(t *testing.T, bookworms, target []Bookworm) bool {
	t.Helper()
	if len(bookworms) != len(target) {
		return false
	}
	for i := range bookworms {
		if bookworms[i].Name != target[i].Name {
			return false
		}
		if !equalBooks(t,
			bookworms[i].Books, target[i].Books) {
			return false
		}
	}
	return true
}

// equalBooks is a helper to test the equality of two lists of Books.
func equalBooks(t *testing.T, books, target []Book) bool {
	t.Helper()

	if len(books) != len(target) {
		return false
	}
	for i := range books {
		if books[i] != target[i] {
			return false
		}
	}
	return true
}

func TestBooksCount(t *testing.T) {
	tests := map[string]struct {
		input []Bookworm
		want  map[Book]uint
	}{
		"nominal use case": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			want: map[Book]uint{
				handmaidsTale: 2,
				theBellJar:    1,
				oryxAndCrake:  1,
				janeEyre:      1},
		},
		"no bookworms": {
			input: []Bookworm{},
			want:  map[Book]uint{},
		},
		"bookworm without books": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{}},
				{Name: "Peggy", Books: []Book{}},
			},
			want: map[Book]uint{},
		},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			got := booksCount(testCase.input)
			if !equalBooksCount(t, got, testCase.want) {
				t.Fatalf("got a different list of books: got %v, expected %v",
					got, testCase.want)
			}

			// equivalent to above use of equalBookworms
			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("different result: got %v, expected %v", got, testCase.want)
			}
		})
	}
}

func equalBooksCount(t *testing.T, got, want map[Book]uint) bool {
	t.Helper()

	if len(got) != len(want) {
		return false
	}

	for book, targetCount := range want {
		count, ok := got[book]
		if !ok || targetCount != count {
			return false
		}
	}

	return true
}

func TestFindCommonBooks(t *testing.T) {
	tests := map[string]struct {
		input []Bookworm
		want  []Book
	}{
		"no common book": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, janeEyre}},
			},
			want: nil,
		},
		"one common book": {
			input: []Bookworm{
				{Name: "Fadi-1", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy-1", Books: []Book{oryxAndCrake, janeEyre, theBellJar}},
			},
			want: []Book{theBellJar},
		},
		"three bookworms have the same books on their shelves": {
			input: []Bookworm{
				{Name: "Fadi-1", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy-1", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Dom-1", Books: []Book{theBellJar, handmaidsTale}},
			},
			want: []Book{handmaidsTale, theBellJar},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := findCommonBooks(tc.input)
			if !equalBooks(t, tc.want, got) {
				t.Fatalf(
					"got a different list of books: %v, expected %v",
					got, tc.want)
			}
		})
	}
}

package helpers

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func TestJoinToNullString(t *testing.T) {
	nulStr123 := sql.NullString{String: "1;2;3", Valid: true}
	nulStrNull := sql.NullString{}
	testCases := []struct {
		slice []string

		expected sql.NullString
	}{
		{[]string{"1", "2", "3"}, nulStr123},
		{[]string{"1", "", "", "2", "3"}, nulStr123},
		{[]string{"", "1", "2", "3"}, nulStr123},
		{[]string{"", "1", "2", "3", ""}, nulStr123},
		{[]string{"", "", "1", "2", "3", ""}, nulStr123},
		{[]string{"", "", "1", "2", "", "", "3", ""}, nulStr123},
		{[]string{"", "   ", ""}, nulStrNull},
		{[]string{""}, nulStrNull},
		{[]string{}, nulStrNull},
		{nil, nulStrNull},
	}

	for i, tt := range testCases {
		res := JoinToNullString(tt.slice)
		if !reflect.DeepEqual(tt.expected, res) {
			t.Fatalf("case %d: expected %#v, got %#v\n", i, tt.expected, res)
		}
	}
}

func TestGetPostImgUrl(t *testing.T) {
	testCases := []struct {
		postID string
		
	}{
		{
			postID: "1",
		},
	}

	for i, test:= range testCases {
		res:= GetPostImgUrl(test.postID)

		fmt.Printf("Test #%d: res: %v\n", i, res)
	}
}

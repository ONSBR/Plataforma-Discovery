package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldParseSQLQuery(t *testing.T) {
	Convey("should parse sql query", t, func() {
		Convey("should parse simple query", func() {
			query := "id = :id"
			parsed := ParseQuery(query, map[string]interface{}{"id": "hello", "other": "bla"})
			So(parsed, ShouldEqual, "id = 'hello'")
		})

		Convey("should parse simple multiparams query", func() {
			query := "id = :id and name = :name"
			parsed := ParseQuery(query, map[string]interface{}{"id": "hello", "name": "my"})
			So(parsed, ShouldEqual, "id = 'hello' and name = 'my'")
		})

		Convey("should parse in list query", func() {
			query := "id in ($ids)"
			parsed := ParseQuery(query, map[string]interface{}{"ids": "1;2;3;4"})
			So(parsed, ShouldEqual, "id in (1,2,3,4)")

			parsed = ParseQuery(query, map[string]interface{}{"ids": "1;2;3;4!"})
			So(parsed, ShouldEqual, "id in ('1','2','3','4')")

			parsed = ParseQuery(query, map[string]interface{}{"ids": "a;b;c;d"})
			So(parsed, ShouldEqual, "id in ('a','b','c','d')")
		})

		Convey("should parse optional params query", func() {
			query := "id in ($ids) [and name = :name]"
			parsed := ParseQuery(query, map[string]interface{}{"ids": "1;2;3;4"})
			So(parsed, ShouldEqual, "id in (1,2,3,4)")
		})

		Convey("should parse with optional params query", func() {
			query := "id in ($ids) [and name = :name]"
			parsed := ParseQuery(query, map[string]interface{}{"ids": "1;2;3;4", "name": "test"})
			So(parsed, ShouldEqual, "id in (1,2,3,4) and name = 'test'")
		})

		Convey("should parse with two optional params query", func() {
			query := "id in ($ids) [and name = :name] [or lastName = :lastName]"
			parsed := ParseQuery(query, map[string]interface{}{"ids": "1;2;3;4", "name": "test", "lastName": "last"})
			So(parsed, ShouldEqual, "id in (1,2,3,4) and name = 'test' or lastName = 'last'")
		})

		Convey("should parse with two optional but just one pass params query", func() {
			query := "id in ($ids) [and name = :name][or lastName = :lastName]"
			parsed := ParseQuery(query, map[string]interface{}{"ids": "1;2;3;4", "lastName": "last"})
			So(parsed, ShouldEqual, "id in (1,2,3,4) or lastName = 'last'")
		})
	})
}

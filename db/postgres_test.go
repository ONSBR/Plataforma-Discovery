package db

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldQueryOnPostgres(t *testing.T) {
	Convey("should connect to postgres", t, func() {
		type conta struct {
			Id    string
			Saldo int
		}
		contas := make([]conta, 0)
		err := Query(func(scan Scan) {
			var c conta
			scan(&c.Id, &c.Saldo)
			contas = append(contas, c)
		}, "select id, saldo from conta where id=$1", "30696c2d-2ffc-4a2e-97d7-d5140534d3ec")

		So(len(contas), ShouldBeGreaterThan, 0)
		So(contas[0].Saldo, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})
}

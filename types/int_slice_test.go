package types

import (
	"testing"

	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func handleErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.Fail()
	}

}
func TestIntSlice(t *testing.T) {
	context, err := utils.NewTestPatrolContext()
	handleErr(t, err)

	_, err = context.DB.Exec(`
		DROP TABLE IF EXISTS testtable;
	`)

	handleErr(t, err)

	_, err = context.DB.Exec(`
		CREATE TABLE IF NOT EXISTS testtable (
			id serial PRIMARY KEY,
			ints integer array
		)
	`)

	handleErr(t, err)

	Convey("test get", t, func() {
		var err error
		sl := IntSlice{}
		sl.Add(1).Add(2).Add(4)

		var id int
		err = context.DB.QueryRowx("INSERT into testtable (ints) VALUES ($1) RETURNING id", sl).Scan(&id)
		So(err, ShouldBeNil)

		sl2 := IntSlice{}
		var id2 int
		err = context.DB.QueryRowx("SELECT id, ints FROM testtable WHERE id = $1", id).Scan(&id2, &sl2)
		So(err, ShouldBeNil)
		So(len(sl2), ShouldEqual, 3)

	})

	Convey("test add unique", t, func() {
		sl := IntSlice{}
		sl.Add(1).Add(2).Add(4)
		sl.AddUnique(1)
		So(len(sl), ShouldEqual, 3)
		sl.AddUnique(55)
		So(len(sl), ShouldEqual, 4)
	})

	Convey("test remove", t, func() {
		sl := IntSlice{}
		sl.Add(1).Add(2).Add(4)
		sl.Remove(1)
		So(len(sl), ShouldEqual, 2)

		// remove non existing
		sl.Remove(66)
		So(len(sl), ShouldEqual, 2)
	})

	Convey("test compact", t, func() {
		sl := IntSlice{}
		sl.Add(1).Add(1).Add(2).Add(2).Add(4)
		sl.Compact()

		sl2 := IntSlice{}
		sl2.Add(1).Add(2).Add(4)
		So(len(sl), ShouldEqual, 3)
		So(sl, ShouldResemble, sl2)
	})

	Convey("test has", t, func() {
		sl := IntSlice{}
		sl.Add(1).Add(2).Add(4)
		So(sl.Has(1), ShouldBeTrue)
		So(sl.Has(77), ShouldBeFalse)

	})

	_, err = context.DB.Exec(`
		DROP TABLE IF EXISTS testtable;
	`)
	handleErr(t, err)
}

package importer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestImporter_New(t *testing.T) {
	Convey("Test Importer", t, func() {
		imp := Importer{}
		Convey("Should return same Alias for same Path", func() {
			Convey("With slash ended", func() {
				So(imp.New("a/b/c"), ShouldEqual, imp.New("a/b/c/"))
				So(imp.Imports(), ShouldResemble, []Import{{Path: "a/b/c", Alias: "c"}})
			})
			Convey("Without slash ended", func() {
				So(imp.New("a/b/c"), ShouldEqual, imp.New("a/b/c"))
				So(imp.Imports(), ShouldResemble, []Import{{Path: "a/b/c", Alias: "c"}})
			})
			Convey("With number ended", func() {
				So(imp.New("a/b/1"), ShouldEqual, imp.New("a/b/1"))
				So(imp.Imports(), ShouldResemble, []Import{{Path: "a/b/1", Alias: "imp1"}})
			})
		})
		Convey("Should return another Alias for same package", func() {
			Convey("With slash ended", func() {
				So(imp.New("a/b/c")+"_1", ShouldEqual, imp.New("a/bb/c/"))
				So(imp.Imports(), ShouldResemble, []Import{{Path: "a/b/c", Alias: "c"}, {Path: "a/bb/c", Alias: "c_1"}})
			})
			Convey("Without slash ended", func() {
				So(imp.New("a/b/c")+"_1", ShouldEqual, imp.New("a/bb/c/"))
				So(imp.Imports(), ShouldResemble, []Import{{Path: "a/b/c", Alias: "c"}, {Path: "a/bb/c", Alias: "c_1"}})
			})
			Convey("With number ended", func() {
				So(imp.New("a/b/1")+"_1", ShouldEqual, imp.New("a/bb/1/"))
				So(imp.Imports(), ShouldResemble, []Import{{Path: "a/b/1", Alias: "imp1"}, {Path: "a/bb/1", Alias: "imp1_1"}})
			})
		})
	})

}

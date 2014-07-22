package goinis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfigFile(t *testing.T) {
	Convey("Load a single configuration file that does exists", t, func() {
		c, err := NewConfigFile("testdata/config_1.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		Convey("test get a section", func() {
			s, err := c.GetSection("parent")
			So(err, ShouldBeNil)
			mockSection := NewSection(c, "parent")
			mockSection.SetKeyValue(newKeyValue("name", "johnnihaoyaa   asfasfhahah"))
			mockSection.SetKeyValue(newKeyValue("relation", "father"))
			mockSection.SetKeyValue(newKeyValue("boolean", "true"))
			mulitKeyValue := newKeyValue("sex[]", "maleqweqw  999")
			mulitKeyValue.addValue("zhangming1 888")
			mockSection.SetKeyValue(mulitKeyValue)
			mockSection.SetKeyValue(newKeyValue("age", "32"))
			nameKeyValue, _ := s.GetKeyValue("name")
			mockNameKeyValue, _ := s.GetKeyValue("name")
			So(nameKeyValue, ShouldResemble, mockNameKeyValue)
			So(s, ShouldNotResemble, mockSection)

			sexKeyValue, _ := s.GetKeyValue("sex[]")
			mockSexKeyValue, _ := s.GetKeyValue("sex[]")
			So(sexKeyValue, ShouldResemble, mockSexKeyValue)

			relationKeyValue, _ := s.GetKeyValue("relation")
			mockRelationKeyValue, _ := s.GetKeyValue("relation")
			So(relationKeyValue, ShouldResemble, mockRelationKeyValue)

			ageInt64, _ := s.Int64("age")
			mockAgeInt64, _ := mockSection.Int64("age")
			So(ageInt64, ShouldEqual, mockAgeInt64)

			boolean, _ := s.Bool("boolean")
			mockBoolean, _ := mockSection.Bool("boolean")
			_, er := mockSection.Bool("boolean_1")
			So(er, ShouldNotBeNil)
			So(boolean, ShouldEqual, mockBoolean)

			notExist1 := s.MustInt64("haha", 9)
			So(notExist1, ShouldEqual, 9)

			notExist2 := s.MustBool("haha", true)
			So(notExist2, ShouldBeTrue)

			notExist3 := s.MustFloat64("haha", 16.0)
			So(notExist3, ShouldEqual, 16.0)

			notExist4 := s.MustStringValue("haha", "九妹九妹漂亮的妹妹")
			So(notExist4, ShouldEqual, "九妹九妹漂亮的妹妹")

			notExist5 := s.MustStringValueRange("relation", "haha", []string{"haha1"})
			So(notExist5, ShouldEqual, "haha")

			notExist6 := s.MustStringValueRange("relation", "ah", []string{"haha1"})
			So(notExist6, ShouldEqual, "ah")
		})

		Convey("test get a section is not exist.", func() {
			s, err := c.GetSection("parent_1")
			So(s, ShouldBeNil)
			So(err, ShouldResemble, &getError{ErrSectionNotFound, "parent_1"})
		})

		Convey("test get subsection", func() {
			s, _ := c.GetSection("hasChildren")
			So(s, ShouldNotBeNil)
			sub1, err := s.GetSubSection("child1")
			So(err, ShouldBeNil)
			So(sub1, ShouldNotBeNil)
			mockSub1 := NewSection(c, "child1")
			mockSub1.SetKeyValue(newKeyValue("name", "child1name"))
			So(sub1, ShouldResemble, mockSub1)
			name, _ := sub1.GetValue("name")
			mockName, _ := mockSub1.GetValue("name")
			So(name, ShouldEqual, mockName)

			subValue, err := s.GetValue("child1.name")
			So(err, ShouldBeNil)
			So("child1name", ShouldEqual, subValue)

		})
	})
}

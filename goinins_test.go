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
			mockSection.SetKeyValue(NewKeyValue("name", "johnnihaoyaa   asfasfhahah"))
			mockSection.SetKeyValue(NewKeyValue("relation", "father"))
			mockSection.SetKeyValue(NewKeyValue("boolean", "true"))
			mulitKeyValue := NewKeyValue("sex[]", "maleqweqw  999")
			mulitKeyValue.AddValue("zhangming1 888")
			mockSection.SetKeyValue(mulitKeyValue)
			mockSection.SetKeyValue(NewKeyValue("age", "32"))
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
		})

		Convey("test get a section is not exist.", func() {
			s, err := c.GetSection("parent_1")
			So(s, ShouldBeNil)
			So(err, ShouldResemble, &getError{ErrSectionNotFound, "parent_1"})
		})
	})
}

package main_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/chiselstrike/libsql-shell/testing/utils"
)

func TestRootCommandExec_GivenEmptyDb_WhenCreateTable_ExpectEmptyResult(t *testing.T) {
	tc := utils.NewTestContext(t)

	result, err := tc.Execute("CREATE TABLE test (id INTEGER PRIMARY KEY, value TEXT)")

	tc.Assert(err, qt.IsNil)
	tc.Assert(result, qt.Equals, "")
}

func TestRootCommandExec_GivenSimpleTableCreated_WhenInsertValue_ExpectEmptyResult(t *testing.T) {
	tc := utils.NewTestContext(t)
	tc.CreateEmptySimpleTable("simple_table")

	result, err := tc.Execute("INSERT INTO simple_table(textField, intField) VALUES ('textValue', 1)")

	tc.Assert(err, qt.IsNil)
	tc.Assert(result, qt.Equals, "")
}

func TestRootCommandExec_GivenSimpleTableCreated_WhenSelectEntireTable_ExpectFirstLineToBeTheHeader(t *testing.T) {
	tc := utils.NewTestContext(t)
	tc.CreateEmptySimpleTable("simple_table")

	result, err := tc.Execute("SELECT * FROM simple_table")
	tc.Assert(err, qt.IsNil)

	headerLine := strings.Split(result, "\n")[0]

	tc.Assert(headerLine, qt.Equals, utils.GetPrintTableOutput([]string{"id", "textField", "intField"}, [][]string{}))
}

func TestRootCommandExec_GivenPopulatedSimpleTable_WhenSelectEntireTable_ExpectAllEntries(t *testing.T) {
	tc := utils.NewTestContext(t)
	tc.CreateSimpleTable("simple_table", []utils.SimpleTableEntry{{TextField: "value1", IntField: 1}, {TextField: "value2", IntField: 2}})

	result, err := tc.Execute("SELECT * FROM simple_table")
	tc.Assert(err, qt.IsNil)

	tc.Assert(result, qt.Equals, utils.GetPrintTableOutput([]string{"id", "textField", "intField"}, [][]string{{"1", "value1", "1"}, {"2", "value2", "2"}}))
}

func TestRootCommandExec_GivenPopulatedSimpleTable_WhenSelectEntireTableTwice_ExpectTwoResults(t *testing.T) {
	tc := utils.NewTestContext(t)
	tc.CreateSimpleTable("simple_table", []utils.SimpleTableEntry{{TextField: "value1", IntField: 1}, {TextField: "value2", IntField: 2}})

	result, err := tc.Execute("SELECT * FROM simple_table; SELECT * FROM simple_table")
	tc.Assert(err, qt.IsNil)

	resultText := utils.GetPrintTableOutput([]string{"id", "textField", "intField"}, [][]string{{"1", "value1", "1"}, {"2", "value2", "2"}})
	resultLines := resultText + "            \n" + resultText
	tc.Assert(result, qt.ContentEquals, resultLines)
}

func TestRootCommandExec_GivenEmptyDb_WhenCreateInsertAndSelectTableAtOnce_ExpectSelectResult(t *testing.T) {
	tc := utils.NewTestContext(t)

	result, err := tc.Execute("CREATE TABLE simple_table (id INTEGER PRIMARY KEY, textField TEXT, intField INTEGER); INSERT INTO simple_table(textField, intField) VALUES ('value1', 1), ('value2', 2); SELECT * FROM simple_table")
	tc.Assert(err, qt.IsNil)

	tc.Assert(result, qt.Equals, utils.GetPrintTableOutput([]string{"id", "textField", "intField"}, [][]string{{"1", "value1", "1"}, {"2", "value2", "2"}}))
}

func TestRootCommandExec_WhenSendStatementWithSemicolonAtEnd_ExpectNoError(t *testing.T) {
	tc := utils.NewTestContext(t)
	tc.CreateSimpleTable("simple_table", []utils.SimpleTableEntry{{TextField: "value1", IntField: 1}, {TextField: "value2", IntField: 2}})

	result, err := tc.Execute("SELECT * FROM simple_table;;;;;;;")
	tc.Assert(err, qt.IsNil)

	tc.Assert(result, qt.Equals, utils.GetPrintTableOutput([]string{"id", "textField", "intField"}, [][]string{{"1", "value1", "1"}, {"2", "value2", "2"}}))

	tc.Assert(err, qt.IsNil)
}

func TestRootCommandExec_GivenSimpleTableCreated_WhenInsertValueWithSemiColumnAndSelectIt_ExpectNoError(t *testing.T) {
	tc := utils.NewTestContext(t)
	tc.CreateEmptySimpleTable("simple_table")

	result, err := tc.Execute("INSERT INTO simple_table(textField, intField) VALUES ('text;Value', 1)")
	tc.Assert(err, qt.IsNil)
	tc.Assert(result, qt.Equals, "")

	result, err = tc.Execute("SELECT * FROM simple_table")
	tc.Assert(err, qt.IsNil)

	tc.Assert(result, qt.Equals, utils.GetPrintTableOutput([]string{"id", "textField", "intField"}, [][]string{{"1", "text;Value", "1"}}))
}
//******************************************************************************************************
//  DataRow.go - Gbtc
//
//  Copyright © 2021, Grid Protection Alliance.  All Rights Reserved.
//
//  Licensed to the Grid Protection Alliance (GPA) under one or more contributor license agreements. See
//  the NOTICE file distributed with this work for additional information regarding copyright ownership.
//  The GPA licenses this file to you under the MIT License (MIT), the "License"; you may not use this
//  file except in compliance with the License. You may obtain a copy of the License at:
//
//      http://opensource.org/licenses/MIT
//
//  Unless agreed to in writing, the subject software distributed under the License is distributed on an
//  "AS-IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. Refer to the
//  License for the specific language governing permissions and limitations.
//
//  Code Modification History:
//  ----------------------------------------------------------------------------------------------------
//  09/23/2021 - J. Ritchie Carroll
//       Generated original version of source code.
//
//******************************************************************************************************

package metadata

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sttp/goapi/sttp/guid"
)

// DataRow represents a row, i.e., a record, in a DataTable defining a set of values for each
// defined DataColumn field in the DataTable columns collection.
type DataRow struct {
	parent *DataTable
	values []interface{}
}

func newDataRow(parent *DataTable) *DataRow {
	return &DataRow{
		parent: parent,
		values: make([]interface{}, parent.ColumnCount()),
	}
}

// Parent gets the parent DataTable of the DataRow.
func (dr *DataRow) Parent() *DataTable {
	return dr.parent
}

func (dr *DataRow) getColumnIndex(columnName string) (int, error) {
	column := dr.parent.ColumnByName(columnName)

	if column == nil {
		return -1, errors.New("Column name \"" + columnName + "\" was not found in table \"" + dr.parent.Name() + "\"")
	}

	return column.Index(), nil
}

func (dr *DataRow) validateColumnType(columnIndex int, targetType int, read bool) (*DataColumn, error) {
	column := dr.parent.Column(columnIndex)

	if column == nil {
		return nil, errors.New("Column index " + strconv.Itoa(columnIndex) + " is out of range for table \"" + dr.parent.Name() + "\"")
	}

	if targetType > -1 && column.Type() != DataTypeEnum(targetType) {
		var action string
		var preposition string

		if read {
			action = "read"
			preposition = "from"
		} else {
			action = "assign"
			preposition = "to"
		}

		panic(fmt.Sprintf("Cannot %s \"%s\" value %s DataColumn \"%s\" for table \"%s\", column data type is \"%s\"", action, DataTypeEnum(targetType).Name(), preposition, column.Name(), dr.parent.Name(), column.Type().Name()))
	}

	if !read && column.Computed() {
		panic("Cannot assign value to DataColumn \"" + column.Name() + "\" for table \"" + dr.parent.Name() + "\", column is computed with an expression")
	}

	return column, nil
}

// func (dr *DataRow) getExpressionTree(column *DataColumn) (*ExpressionTree, error) {
// 	columnIndex := column.Index()

// 	if dr.values[columnIndex] == nil {
// 		dataTable := column.Parent()
// 		parser := NewFilterExpressionParser(column.Expression())

// 		parser.SetDataSet(dataTable.Parent())
// 		parser.SetPrimaryTableName(dataTable.Name())
// 		parser.SetTrackFilteredSignalIDs(false)
// 		parser.SetTrackFilteredRows(false)

// 		expressionTrees := parser.GetExpressionTrees()

// 		if len(expressionTrees) == 0 {
// 			return nil, errors.New("Expression defined for computed DataColumn \"" + column.Name() + "\" for table \"" + dr.parent.Name() + "\" cannot produce a value")
// 		}

// 		dr.values[columnIndex] = parser
// 		return expressionTrees[0]
// 	}

// 	return dr.values[columnIndex].(*FilterExpressionParser).GetExpressionTrees()[0]
// }

func (dr *DataRow) getComputedValue(column *DataColumn) (interface{}, error) {
	// TODO: Evaluate expression using ANTLR grammar:
	// https://github.com/sttp/cppapi/blob/master/src/lib/filterexpressions/FilterExpressionSyntax.g4
	// expressionTree, err := dr.getExpressionTree(column)
	// sourceValue = expressionTree.Evaluate()

	// switch sourceValue.ValueType {
	// case ExpressionValueType.Boolean:
	// }

	return nil, nil
}

// Value reads the record value at the specified columnIndex.
func (dr *DataRow) Value(columnIndex int) (interface{}, error) {
	column, err := dr.validateColumnType(columnIndex, -1, true)

	if err != nil {
		return nil, err
	}

	if column.Computed() {
		return dr.getComputedValue(column)
	}

	return dr.values[columnIndex], nil
}

// ValueByName reads the record value for the specified columnName.
func (dr *DataRow) ValueByName(columnName string) (interface{}, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return nil, err
	}

	return dr.values[index], nil
}

// SetValue assigns the record value at the specified columnIndex.
func (dr *DataRow) SetValue(columnIndex int, value interface{}) error {
	_, err := dr.validateColumnType(columnIndex, -1, false)

	if err != nil {
		return err
	}

	dr.values[columnIndex] = value
	return nil
}

// SetValueByName assigns the record value for the specified columnName.
func (dr *DataRow) SetValueByName(columnName string, value interface{}) error {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return err
	}

	return dr.SetValue(index, value)
}

// ColumnValueAsString reads the record value for the specified data column
// converted to a string. For any errors, an empty string will be returned.
func (dr *DataRow) ColumnValueAsString(column *DataColumn) string {
	if column == nil {
		return ""
	}

	index := column.Index()

	switch column.Type() {
	case DataType.String:
		value, err := dr.StringValue(index)

		if err != nil {
			return ""
		}

		return value
	case DataType.Boolean:
		value, err := dr.BoolValue(index)

		if err != nil {
			return ""
		}

		return strconv.FormatBool(value)
	case DataType.DateTime:
		value, err := dr.DateTimeValue(index)

		if err != nil {
			return ""
		}
		return value.String()
	case DataType.Single:
		value, err := dr.SingleValue(index)

		if err != nil {
			return ""
		}
		return strconv.FormatFloat(float64(value), 'f', 6, 32)
	case DataType.Decimal:
		fallthrough
	case DataType.Double:
		value, err := dr.DoubleValue(index)

		if err != nil {
			return ""
		}

		return strconv.FormatFloat(value, 'f', 6, 64)
	case DataType.Guid:
		value, err := dr.GuidValue(index)

		if err != nil {
			return ""
		}

		return value.String()
	case DataType.Int8:
		value, err := dr.Int8Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatInt(int64(value), 10)
	case DataType.Int16:
		value, err := dr.Int16Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatInt(int64(value), 10)
	case DataType.Int32:
		value, err := dr.Int32Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatInt(int64(value), 10)
	case DataType.Int64:
		value, err := dr.Int64Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatInt(value, 10)
	case DataType.UInt8:
		value, err := dr.UInt8Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatUint(uint64(value), 10)
	case DataType.UInt16:
		value, err := dr.UInt16Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatUint(uint64(value), 10)
	case DataType.UInt32:
		value, err := dr.UInt32Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatUint(uint64(value), 10)
	case DataType.UInt64:
		value, err := dr.UInt64Value(index)

		if err != nil {
			return ""
		}

		return strconv.FormatUint(value, 10)
	default:
		return ""
	}
}

// ValueAsString reads the record value at the specified columnIndex converted to a string.
// For columnIndex out of range or any other errors, an empty string will be returned.
func (dr *DataRow) ValueAsString(columnIndex int) string {
	return dr.ColumnValueAsString(dr.parent.Column(columnIndex))
}

// ValueAsStringByName reads the record value for the specified columnName converted to a string.
// For columnName not found or any other errors, an empty string will be returned.
func (dr *DataRow) ValueAsStringByName(columnName string) string {
	return dr.ColumnValueAsString(dr.parent.ColumnByName(columnName))
}

// StringValue gets the record value at the specified columnIndex cast as a string.
// An error will be returned if column type is not DataType.String.
func (dr *DataRow) StringValue(columnIndex int) (string, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.String), true)

	if err != nil {
		return "", err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return "", err
		}

		return value.(string), nil
	}

	return dr.values[columnIndex].(string), nil
}

// StringValueByName gets the record value for the specified columnName cast as a string.
// An error will be returned if column type is not DataType.String.
func (dr *DataRow) StringValueByName(columnName string) (string, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return "", err
	}

	return dr.StringValue(index)
}

// BoolValue gets the record value at the specified columnIndex cast as a bool.
// An error will be returned if column type is not DataType.Boolean.
func (dr *DataRow) BoolValue(columnIndex int) (bool, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Boolean), true)

	if err != nil {
		return false, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return false, err
		}

		return value.(bool), nil
	}

	return dr.values[columnIndex].(bool), nil
}

// BoolValueByName gets the record value for the specified columnName cast as a bool.
// An error will be returned if column type is not DataType.Boolean.
func (dr *DataRow) BoolValueByName(columnName string) (bool, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return false, err
	}

	return dr.BoolValue(index)
}

// DateTimeValue gets the record value at the specified columnIndex cast as a time.Time.
// An error will be returned if column type is not DataType.DateTime.
func (dr *DataRow) DateTimeValue(columnIndex int) (time.Time, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.DateTime), true)

	if err != nil {
		return time.Time{}, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return time.Time{}, err
		}

		return value.(time.Time), nil
	}

	return dr.values[columnIndex].(time.Time), nil
}

// DateTimeValueByName gets the record value for the specified columnName cast as a time.Time.
// An error will be returned if column type is not DataType.DateTime.
func (dr *DataRow) DateTimeValueByName(columnName string) (time.Time, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return time.Time{}, err
	}

	return dr.DateTimeValue(index)
}

// SingleValue gets the record value at the specified columnIndex cast as a float32.
// An error will be returned if column type is not DataType.Single.
func (dr *DataRow) SingleValue(columnIndex int) (float32, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Single), true)

	if err != nil {
		return 0.0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0.0, err
		}

		return value.(float32), nil
	}

	return dr.values[columnIndex].(float32), nil
}

// SingleValueByName gets the record value for the specified columnName cast as a float32.
// An error will be returned if column type is not DataType.Single.
func (dr *DataRow) SingleValueByName(columnName string) (float32, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0.0, err
	}

	return dr.SingleValue(index)
}

// DoubleValue gets the record value at the specified columnIndex cast as a float64.
// An error will be returned if column type is not DataType.Double.
func (dr *DataRow) DoubleValue(columnIndex int) (float64, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Double), true)

	if err != nil {
		return 0.0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0.0, err
		}

		return value.(float64), nil
	}

	return dr.values[columnIndex].(float64), nil
}

// DoubleValueByName gets the record value for the specified columnName cast as a float64.
// An error will be returned if column type is not DataType.Double.
func (dr *DataRow) DoubleValueByName(columnName string) (float64, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0.0, err
	}

	return dr.DoubleValue(index)
}

// DecimalValue gets the record value at the specified columnIndex cast as a float64.
// An error will be returned if column type is not DataType.Decimal.
func (dr *DataRow) DecimalValue(columnIndex int) (float64, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Decimal), true)

	if err != nil {
		return 0.0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0.0, err
		}

		return value.(float64), nil
	}

	return dr.values[columnIndex].(float64), nil
}

// DecimalValueByName gets the record value for the specified columnName cast as a float64.
// An error will be returned if column type is not DataType.Decimal.
func (dr *DataRow) DecimalValueByName(columnName string) (float64, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0.0, err
	}

	return dr.DecimalValue(index)
}

// GuidValue gets the record value at the specified columnIndex cast as a guid.Guid.
// An error will be returned if column type is not DataType.Guid.
func (dr *DataRow) GuidValue(columnIndex int) (guid.Guid, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Guid), true)

	if err != nil {
		return guid.Guid{}, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return guid.Guid{}, err
		}

		return value.(guid.Guid), nil
	}

	return dr.values[columnIndex].(guid.Guid), nil
}

// GuidValueByName gets the record value for the specified columnName cast as a guid.Guid.
// An error will be returned if column type is not DataType.Guid.
func (dr *DataRow) GuidValueByName(columnName string) (guid.Guid, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return guid.Guid{}, err
	}

	return dr.GuidValue(index)
}

// Int8Value gets the record value at the specified columnIndex cast as a int8.
// An error will be returned if column type is not DataType.Int8.
func (dr *DataRow) Int8Value(columnIndex int) (int8, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Int8), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(int8), nil
	}

	return dr.values[columnIndex].(int8), nil
}

// Int8ValueByName gets the record value for the specified columnName cast as a int8.
// An error will be returned if column type is not DataType.Int8.
func (dr *DataRow) Int8ValueByName(columnName string) (int8, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.Int8Value(index)
}

// Int16Value gets the record value at the specified columnIndex cast as a int16.
// An error will be returned if column type is not DataType.Int16.
func (dr *DataRow) Int16Value(columnIndex int) (int16, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Int16), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(int16), nil
	}

	return dr.values[columnIndex].(int16), nil
}

// Int16ValueByName gets the record value for the specified columnName cast as a int16.
// An error will be returned if column type is not DataType.Int16.
func (dr *DataRow) Int16ValueByName(columnName string) (int16, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.Int16Value(index)
}

// Int32Value gets the record value at the specified columnIndex cast as a int32.
// An error will be returned if column type is not DataType.Int32.
func (dr *DataRow) Int32Value(columnIndex int) (int32, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Int32), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(int32), nil
	}

	return dr.values[columnIndex].(int32), nil
}

// Int32ValueByName gets the record value for the specified columnName cast as a int32.
// An error will be returned if column type is not DataType.Int32.
func (dr *DataRow) Int32ValueByName(columnName string) (int32, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.Int32Value(index)
}

// Int64Value gets the record value at the specified columnIndex cast as a int64.
// An error will be returned if column type is not DataType.Int64.
func (dr *DataRow) Int64Value(columnIndex int) (int64, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.Int64), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(int64), nil
	}

	return dr.values[columnIndex].(int64), nil
}

// Int64ValueByName gets the record value for the specified columnName cast as a int64.
// An error will be returned if column type is not DataType.Int64.
func (dr *DataRow) Int64ValueByName(columnName string) (int64, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.Int64Value(index)
}

// UInt8Value gets the record value at the specified columnIndex cast as a uint8.
// An error will be returned if column type is not DataType.UInt8.
func (dr *DataRow) UInt8Value(columnIndex int) (uint8, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.UInt8), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(uint8), nil
	}

	return dr.values[columnIndex].(uint8), nil
}

// UInt8ValueByName gets the record value for the specified columnName cast as a uint8.
// An error will be returned if column type is not DataType.UInt8.
func (dr *DataRow) UInt8ValueByName(columnName string) (uint8, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.UInt8Value(index)
}

// UInt16Value gets the record value at the specified columnIndex cast as a uint16.
// An error will be returned if column type is not DataType.UInt16.
func (dr *DataRow) UInt16Value(columnIndex int) (uint16, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.UInt16), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(uint16), nil
	}

	return dr.values[columnIndex].(uint16), nil
}

// UInt16ValueByName gets the record value for the specified columnName cast as a uint16.
// An error will be returned if column type is not DataType.UInt16.
func (dr *DataRow) UInt16ValueByName(columnName string) (uint16, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.UInt16Value(index)
}

// UInt32Value gets the record value at the specified columnIndex cast as a uint32.
// An error will be returned if column type is not DataType.UInt32.
func (dr *DataRow) UInt32Value(columnIndex int) (uint32, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.UInt32), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(uint32), nil
	}

	return dr.values[columnIndex].(uint32), nil
}

// UInt32ValueByName gets the record value for the specified columnName cast as a uint32.
// An error will be returned if column type is not DataType.UInt32.
func (dr *DataRow) UInt32ValueByName(columnName string) (uint32, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.UInt32Value(index)
}

// UInt64Value gets the record value at the specified columnIndex cast as a uint64.
// An error will be returned if column type is not DataType.UInt64.
func (dr *DataRow) UInt64Value(columnIndex int) (uint64, error) {
	column, err := dr.validateColumnType(columnIndex, int(DataType.UInt64), true)

	if err != nil {
		return 0, err
	}

	if column.Computed() {
		value, err := dr.getComputedValue(column)

		if err != nil {
			return 0, err
		}

		return value.(uint64), nil
	}

	return dr.values[columnIndex].(uint64), nil
}

// UInt64ValueByName gets the record value for the specified columnName cast as a uint64.
// An error will be returned if column type is not DataType.UInt64.
func (dr *DataRow) UInt64ValueByName(columnName string) (uint64, error) {
	index, err := dr.getColumnIndex(columnName)

	if err != nil {
		return 0, err
	}

	return dr.UInt64Value(index)
}
package spec

import (
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/timestamp-model/operation/timestamp"
	"github.com/imfact-labs/timestamp-model/state"
	"github.com/imfact-labs/timestamp-model/types"
)

var AddedHinters = []encoder.DecodeDetail{
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.ItemHint, Instance: types.Item{}},
	{Hint: timestamp.IssueHint, Instance: timestamp.Issue{}},
	{Hint: timestamp.RegisterModelHint, Instance: timestamp.RegisterModel{}},
	{Hint: state.ItemStateValueHint, Instance: state.ItemStateValue{}},
	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.LastIdxStateValueHint, Instance: state.LastIdxStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: timestamp.IssueFactHint, Instance: timestamp.IssueFact{}},
	{Hint: timestamp.RegisterModelFactHint, Instance: timestamp.RegisterModelFact{}},
}

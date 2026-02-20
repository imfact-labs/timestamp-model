package cmds

import (
	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/timestamp-model/operation/timestamp"
	"github.com/imfact-labs/timestamp-model/state"
	"github.com/imfact-labs/timestamp-model/types"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit

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

func init() {
	Hinters = append(Hinters, ccmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}

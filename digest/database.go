package digest

import (
	cdigest "github.com/imfact-labs/currency-model/digest"
	utilc "github.com/imfact-labs/currency-model/digest/util"
	"github.com/imfact-labs/mitum2/base"
	utilm "github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/timestamp-model/state"
	"github.com/imfact-labs/timestamp-model/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	DefaultColNameTimeStamp = "digest_ts"
)

func TimestampDesign(st *cdigest.Database, contract string) (types.Design, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("isItem", false)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var sta base.State
	if err := st.MongoClient().GetByFilter(
		DefaultColNameTimeStamp,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return types.Design{}, nil, utilm.ErrNotFound.WithMessage(err, "timestamp design by contract account %v", contract)
	}

	if sta != nil {
		de, err := state.GetDesignFromState(sta)
		if err != nil {
			return types.Design{}, nil, err
		}
		return de, sta, nil
	} else {
		return types.Design{}, nil, errors.Errorf("state is nil")
	}
}

func TimestampItem(db *cdigest.Database, contract, project string, idx uint64) (types.Item, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("project_id", project)
	filter = filter.Add("timestamp_idx", idx)
	filter = filter.Add("isItem", true)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var st base.State
	if err := db.MongoClient().GetByFilter(
		DefaultColNameTimeStamp,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			st = i
			return nil
		},
		opt,
	); err != nil {
		return types.Item{}, nil, utilm.ErrNotFound.WithMessage(err, "timestamp item by contract account %s, project %s, timestamp idx %d", contract, project, idx)
	}

	if st != nil {
		it, err := state.GetItemFromState(st)
		if err != nil {
			return types.Item{}, nil, err
		}
		return it, st, nil
	} else {
		return types.Item{}, nil, errors.Errorf("state is nil")
	}
}

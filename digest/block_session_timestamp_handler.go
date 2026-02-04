package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func PrepareTimeStamps(bs *currencydigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
	switch {
	case state.IsDesignStateKey(st.Key()):
		j, err := handleTimeStampDesignState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameTimeStamp, j, nil
	case state.IsItemStateKey(st.Key()):
		j, err := handleTimeStampItemState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameTimeStamp, j, nil
	}

	return "", nil, nil
}

func handleTimeStampDesignState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if serviceDesignDoc, err := NewDesignDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(serviceDesignDoc),
		}, nil
	}
}

func handleTimeStampItemState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if TimeStampItemDoc, err := NewItemDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(TimeStampItemDoc),
		}, nil
	}
}

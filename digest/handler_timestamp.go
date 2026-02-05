package digest

import (
	"net/http"
	"strconv"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
)

var (
	HandlerPathTimeStampDesign = `/timestamp/{contract:(?i)` + ctypes.REStringAddressString + `}`
	HandlerPathTimeStampItem   = `/timestamp/{contract:(?i)` + ctypes.REStringAddressString + `}/project/{project_id:` + ctypes.ReSpecialCh + `}/idx/{timestamp_idx:[0-9]+}`
)

func SetHandlers(hd *cdigest.Handlers) {
	get := 1000
	_ = hd.SetHandler(HandlerPathTimeStampItem, HandleTimeStampItem, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathTimeStampDesign, HandleTimeStampDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func HandleTimeStampDesign(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cacheKey, func() (interface{}, error) {
		return handleTimeStampDesignInGroup(hd, contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.ExpireShortLived())
		}
	}
}

func handleTimeStampDesignInGroup(hd *cdigest.Handlers, contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := TimestampDesign(hd.Database(), contract)
	if err != nil {
		return nil, err
	}

	i, err := buildTimeStampDesign(hd, contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.Encoder().Marshal(i)
}

func buildTimeStampDesign(hd *cdigest.Handlers, contract string, de types.Design, st base.State) (cdigest.Hal, error) {
	h, err := hd.CombineURL(HandlerPathTimeStampDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(de, cdigest.NewHalLink(h, nil))

	h, err = hd.CombineURL(cdigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.CombineURL(cdigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func HandleTimeStampItem(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	project, err, status := cdigest.ParseRequest(w, r, "project_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	s, found := mux.Vars(r)["timestamp_idx"]
	if !found {
		cdigest.HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}
	idx, err := parseIdxFromPath(s)
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleTimeStampItemInGroup(hd, contract, project, idx)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, hd.ExpireLongLived())
		}
	}
}

func handleTimeStampItemInGroup(hd *cdigest.Handlers, contract, project string, idx uint64) ([]byte, error) {
	var it types.Item
	var st base.State

	it, st, err := TimestampItem(hd.Database(), contract, project, idx)
	if err != nil {
		return nil, err
	}

	i, err := buildTimeStampItem(hd, contract, it, st)
	if err != nil {
		return nil, err
	}
	return hd.Encoder().Marshal(i)
}

func buildTimeStampItem(hd *cdigest.Handlers, contract string, it types.Item, st base.State) (cdigest.Hal, error) {
	h, err := hd.CombineURL(
		HandlerPathTimeStampItem,
		"contract", contract, "project_id", it.ProjectID(), "timestamp_idx",
		strconv.FormatUint(it.TimestampID(), 10))
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(it, cdigest.NewHalLink(h, nil))

	h, err = hd.CombineURL(cdigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.CombineURL(cdigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}

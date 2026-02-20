package api

import (
	"net/http"
	"strconv"
	"strings"

	apic "github.com/imfact-labs/currency-model/api"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/timestamp-model/digest"
	"github.com/imfact-labs/timestamp-model/types"
	"github.com/pkg/errors"
)

var (
	HandlerPathTimeStampDesign = `/timestamp/{contract:(?i)` + ctypes.REStringAddressString + `}`
	HandlerPathTimeStampItem   = `/timestamp/{contract:(?i)` + ctypes.REStringAddressString + `}/project/{project_id:` + ctypes.ReSpecialCh + `}/idx/{timestamp_idx:[0-9]+}`
)

func SetHandlers(hd *apic.Handlers) {
	get := 1000
	_ = hd.SetHandler(HandlerPathTimeStampItem, HandleTimeStampItem, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathTimeStampDesign, HandleTimeStampDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func HandleTimeStampDesign(hd *apic.Handlers, w http.ResponseWriter, r *http.Request) {
	cacheKey := apic.CacheKeyPath(r)
	if err := apic.LoadFromCache(hd.Cache(), cacheKey, w); err == nil {
		return
	}

	contract, err, status := apic.ParseRequest(w, r, "contract")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cacheKey, func() (interface{}, error) {
		return handleTimeStampDesignInGroup(hd, contract)
	}); err != nil {
		apic.HTTP2HandleError(w, err)
	} else {
		apic.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)

		if !shared {
			apic.HTTP2WriteCache(w, cacheKey, hd.ExpireShortLived())
		}
	}
}

func handleTimeStampDesignInGroup(hd *apic.Handlers, contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := digest.TimestampDesign(hd.Database(), contract)
	if err != nil {
		return nil, err
	}

	i, err := buildTimeStampDesign(hd, contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.Encoder().Marshal(i)
}

func buildTimeStampDesign(hd *apic.Handlers, contract string, de types.Design, st base.State) (apic.Hal, error) {
	h, err := hd.CombineURL(HandlerPathTimeStampDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal apic.Hal
	hal = apic.NewBaseHal(de, apic.NewHalLink(h, nil))

	h, err = hd.CombineURL(apic.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", apic.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.CombineURL(apic.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", apic.NewHalLink(h, nil))
	}

	return hal, nil
}

func HandleTimeStampItem(hd *apic.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := apic.CacheKeyPath(r)
	if err := apic.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := apic.ParseRequest(w, r, "contract")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	project, err, status := apic.ParseRequest(w, r, "project_id")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	s, err, status := apic.ParseRequest(w, r, "timestamp_idx")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}
	idx, err := parseIdxFromPath(s)
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleTimeStampItemInGroup(hd, contract, project, idx)
	}); err != nil {
		apic.HTTP2HandleError(w, err)
	} else {
		apic.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)

		if !shared {
			apic.HTTP2WriteCache(w, cachekey, hd.ExpireLongLived())
		}
	}
}

func handleTimeStampItemInGroup(hd *apic.Handlers, contract, project string, idx uint64) ([]byte, error) {
	var it types.Item
	var st base.State

	it, st, err := digest.TimestampItem(hd.Database(), contract, project, idx)
	if err != nil {
		return nil, err
	}

	i, err := buildTimeStampItem(hd, contract, it, st)
	if err != nil {
		return nil, err
	}
	return hd.Encoder().Marshal(i)
}

func buildTimeStampItem(hd *apic.Handlers, contract string, it types.Item, st base.State) (apic.Hal, error) {
	h, err := hd.CombineURL(
		HandlerPathTimeStampItem,
		"contract", contract, "project_id", it.ProjectID(), "timestamp_idx",
		strconv.FormatUint(it.TimestampID(), 10))
	if err != nil {
		return nil, err
	}

	var hal apic.Hal
	hal = apic.NewBaseHal(it, apic.NewHalLink(h, nil))

	h, err = hd.CombineURL(apic.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", apic.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.CombineURL(apic.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", apic.NewHalLink(h, nil))
	}

	return hal, nil
}

func parseIdxFromPath(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return 0, errors.Errorf("empty idx")
	} else if len(s) > 1 && strings.HasPrefix(s, "0") {
		return 0, errors.Errorf("invalid idx, %q", s)
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

package keeper

import (
	"fmt"
	"math"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/query"
)

// Customized version of query.Paginate to paginate array instead of KVStore
func PaginateArray(
	length int,
	pageReq *query.PageRequest,
	onResult func(i int) error,
	defaultLimit int,
	maxLimit int,
) (pageRes *query.PageResponse, err error) {
	if pageReq == nil {
		pageReq = &query.PageRequest{}
	}

	var limit int
	var offset int
	var key int
	reverse := pageReq.Reverse
	// note we always return total count, we have this info anyways (unlike KVStore)

	if pageReq.Offset > math.MaxInt {
		return nil, fmt.Errorf("offset out of range")
	}
	offset = int(pageReq.Offset)

	if pageReq.Limit > uint64(maxLimit) || pageReq.Limit > math.MaxInt {
		return nil, fmt.Errorf("limit out of range")
	}

	limit = int(pageReq.Limit)

	if limit == 0 {
		limit = defaultLimit
	}

	_key := pageReq.Key
	if offset > 0 && _key != nil {
		return nil, fmt.Errorf("either offset or key is expected, got both")
	}

	if _key != nil {
		// using key
		key, err = strconv.Atoi(string(_key))
		if err != nil {
			return nil, fmt.Errorf("invalid key: %s", err.Error())
		}
		if key < 0 || key >= length {
			return nil, fmt.Errorf("key out of range")
		}
	} else {
		// using offset, adopt to key
		if reverse {
			key = length - 1 - offset
		} else { // normal
			key = offset
		}
		if key < 0 || key >= length { // no more items to return
			return &query.PageResponse{
				Total: uint64(length),
			}, nil
		}
	}

	count := 0
	_next := -1 // temp value, default to out-of-range index
	if reverse {
		for i := key; i >= 0 && count < limit; i-- {
			err := onResult(i)
			if err != nil {
				return nil, err
			}
			count += 1
			_next = i - 1
		}
	} else { // normal
		for i := key; i < length && count < limit; i++ {
			err := onResult(i)
			if err != nil {
				return nil, err
			}
			count += 1
			_next = i + 1
		}
	}

	res := query.PageResponse{
		Total: uint64(length),
	}

	if _next >= 0 && _next < length { // only set value if valid
		res.NextKey = []byte(strconv.Itoa(_next))
	}

	return &res, nil
}

func PaginateStringArray(
	arr []string,
	pageReq *query.PageRequest,
	onResult func(i int, val string) error,
	defaultLimit int,
	maxLimit int,
) (pageRes *query.PageResponse, err error) {
	return PaginateArray(len(arr), pageReq, func(i int) error {
		return onResult(i, arr[i])
	}, defaultLimit, maxLimit)
}

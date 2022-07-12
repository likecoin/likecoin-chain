package keeper_test

import (
	"math"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/keeper"
	"github.com/stretchr/testify/require"
)

func TestPaginationNormalOffset(t *testing.T) {
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit: 3,
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("3"),
		Total:   uint64(10),
	}, res1)
	require.Equal(t, []int{0, 1, 2}, actualPage1)

	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:  3,
		Offset: 3,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("6"),
		Total:   uint64(10),
	}, res2)
	require.Equal(t, []int{3, 4, 5}, actualPage2)

	var actualPage3 []int
	res3, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:  3,
		Offset: 6,
	}, func(i int) error {
		actualPage3 = append(actualPage3, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("9"),
		Total:   uint64(10),
	}, res3)
	require.Equal(t, []int{6, 7, 8}, actualPage3)

	var actualPage4 []int
	res4, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:  3,
		Offset: 9,
	}, func(i int) error {
		actualPage4 = append(actualPage4, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res4)
	require.Equal(t, []int{9}, actualPage4)
}

func TestPaginationNormalKey(t *testing.T) {
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit: 4,
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("4"),
		Total:   uint64(10),
	}, res1)
	require.Equal(t, []int{0, 1, 2, 3}, actualPage1)

	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit: 4,
		Key:   []byte("4"),
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("8"),
		Total:   uint64(10),
	}, res2)
	require.Equal(t, []int{4, 5, 6, 7}, actualPage2)

	var actualPage3 []int
	res3, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit: 4,
		Key:   []byte("8"),
	}, func(i int) error {
		actualPage3 = append(actualPage3, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res3)
	require.Equal(t, []int{8, 9}, actualPage3)
}

func TestPaginationReverseOffset(t *testing.T) {
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   3,
		Reverse: true,
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("6"),
		Total:   uint64(10),
	}, res1)
	require.Equal(t, []int{9, 8, 7}, actualPage1)

	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   3,
		Offset:  3,
		Reverse: true,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("3"),
		Total:   uint64(10),
	}, res2)
	require.Equal(t, []int{6, 5, 4}, actualPage2)

	var actualPage3 []int
	res3, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   3,
		Offset:  6,
		Reverse: true,
	}, func(i int) error {
		actualPage3 = append(actualPage3, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("0"),
		Total:   uint64(10),
	}, res3)
	require.Equal(t, []int{3, 2, 1}, actualPage3)

	var actualPage4 []int
	res4, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   3,
		Offset:  9,
		Reverse: true,
	}, func(i int) error {
		actualPage4 = append(actualPage4, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res4)
	require.Equal(t, []int{0}, actualPage4)
}

func TestPaginationReverseKey(t *testing.T) {
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   4,
		Reverse: true,
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("5"),
		Total:   uint64(10),
	}, res1)
	require.Equal(t, []int{9, 8, 7, 6}, actualPage1)

	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   4,
		Key:     []byte("5"),
		Reverse: true,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("1"),
		Total:   uint64(10),
	}, res2)
	require.Equal(t, []int{5, 4, 3, 2}, actualPage2)

	var actualPage3 []int
	res3, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit:   4,
		Key:     []byte("1"),
		Reverse: true,
	}, func(i int) error {
		actualPage3 = append(actualPage3, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res3)
	require.Equal(t, []int{1, 0}, actualPage3)
}

func TestPaginationOutOfRangeOffset(t *testing.T) {
	// Normal
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Offset: 10,
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res1)
	require.Equal(t, []int(nil), actualPage1)

	// Reverse
	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Offset:  10,
		Reverse: true,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res2)
	require.Equal(t, []int(nil), actualPage2)

	// Golang array size limit
	var actualPage3 []int
	res3, err := keeper.PaginateArray(10, &query.PageRequest{
		Offset: math.MaxUint64,
	}, func(i int) error {
		actualPage3 = append(actualPage3, i)
		return nil
	}, 5, 10)
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset out of range")
	require.Nil(t, res3)
	require.Equal(t, []int(nil), actualPage3)

	// No more pages
	var actualPage4 []int
	res4, err := keeper.PaginateArray(10, &query.PageRequest{
		Offset: 10,
	}, func(i int) error {
		actualPage4 = append(actualPage4, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res4)
	require.Equal(t, []int(nil), actualPage4)
}

func TestPaginationOutOfRangeLimit(t *testing.T) {
	// Golang array size limit
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit: math.MaxUint64,
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit out of range")
	require.Nil(t, res1)
	require.Equal(t, []int(nil), actualPage1)

	// App defined limit
	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Limit: 6,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 5)
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit out of range")
	require.Nil(t, res2)
	require.Equal(t, []int(nil), actualPage2)
}

func TestPaginationOutOfRangeKey(t *testing.T) {
	// Normal
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, &query.PageRequest{
		Key: []byte("10"),
	}, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.Error(t, err)
	require.Contains(t, err.Error(), "key out of range")
	require.Nil(t, res1)
	require.Equal(t, []int(nil), actualPage1)

	// Reverse
	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Key:     []byte("10"),
		Reverse: true,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.Error(t, err)
	require.Contains(t, err.Error(), "key out of range")
	require.Nil(t, res2)
	require.Equal(t, []int(nil), actualPage2)
}

func TestPaginationDefaults(t *testing.T) {
	// all default
	var actualPage1 []int
	res1, err := keeper.PaginateArray(10, nil, func(i int) error {
		actualPage1 = append(actualPage1, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("5"),
		Total:   uint64(10),
	}, res1)
	require.Equal(t, []int{0, 1, 2, 3, 4}, actualPage1)

	// reverse only
	var actualPage2 []int
	res2, err := keeper.PaginateArray(10, &query.PageRequest{
		Reverse: true,
	}, func(i int) error {
		actualPage2 = append(actualPage2, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: []byte("4"),
		Total:   uint64(10),
	}, res2)
	require.Equal(t, []int{9, 8, 7, 6, 5}, actualPage2)

	// key only
	var actualPage3 []int
	res3, err := keeper.PaginateArray(10, &query.PageRequest{
		Key: []byte("5"),
	}, func(i int) error {
		actualPage3 = append(actualPage3, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res3)
	require.Equal(t, []int{5, 6, 7, 8, 9}, actualPage3)

	// offset only
	var actualPage4 []int
	res4, err := keeper.PaginateArray(10, &query.PageRequest{
		Offset: 5,
	}, func(i int) error {
		actualPage4 = append(actualPage4, i)
		return nil
	}, 5, 10)
	require.NoError(t, err)
	require.Equal(t, &query.PageResponse{
		NextKey: nil,
		Total:   uint64(10),
	}, res4)
	require.Equal(t, []int{5, 6, 7, 8, 9}, actualPage4)

	// limit only covered by normal case
}

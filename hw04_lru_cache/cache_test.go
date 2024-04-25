package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type item struct {
	key Key
	val string
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

type TestSuite struct {
	suite.Suite
	Cache Cache
}

func (tS *TestSuite) SetupTest() {
	tS.Cache = NewCache(5)
}

func (tS *TestSuite) TestEmptyCache() {
	_, ok := tS.Cache.Get("aaa")
	tS.False(ok)

	_, ok = tS.Cache.Get("bbb")
	tS.False(ok)
}

func (tS *TestSuite) TestSimpleFilling() {
	wasInCache := tS.Cache.Set("aaa", 100)
	tS.Require().False(wasInCache)

	wasInCache = tS.Cache.Set("bbb", 200)
	tS.Require().False(wasInCache)

	val, ok := tS.Cache.Get("aaa")
	tS.Require().True(ok)
	tS.Require().Equal(100, val)

	val, ok = tS.Cache.Get("bbb")
	tS.Require().True(ok)
	tS.Require().Equal(200, val)

	wasInCache = tS.Cache.Set("aaa", 300)
	tS.Require().True(wasInCache)

	val, ok = tS.Cache.Get("aaa")
	tS.Require().True(ok)
	tS.Require().Equal(300, val)

	val, ok = tS.Cache.Get("ccc")
	tS.Require().False(ok)
	tS.Require().Nil(val)
}

func (tS *TestSuite) TestPurge() {
	testCases := []struct {
		Item     item
		Expected []item
	}{
		{
			Item: item{key: "a", val: "10"},
			Expected: []item{
				{key: "a", val: "10"},
			},
		},
		{
			Item: item{key: "b", val: "20"},
			Expected: []item{
				{key: "b", val: "20"},
				{key: "a", val: "10"},
			},
		},
		{
			Item: item{key: "c", val: "30"},
			Expected: []item{
				{key: "c", val: "30"},
				{key: "b", val: "20"},
				{key: "a", val: "10"},
			},
		},
		{
			Item: item{key: "d", val: "40"},
			Expected: []item{
				{key: "d", val: "40"},
				{key: "c", val: "30"},
				{key: "b", val: "20"},
				{key: "a", val: "10"},
			},
		},
		{
			Item: item{key: "e", val: "50"},
			Expected: []item{
				{key: "e", val: "50"},
				{key: "d", val: "40"},
				{key: "c", val: "30"},
				{key: "b", val: "20"},
				{key: "a", val: "10"},
			},
		},
		{
			Item: item{key: "f", val: "60"},
			Expected: []item{
				{key: "f", val: "60"},
				{key: "e", val: "50"},
				{key: "d", val: "40"},
				{key: "c", val: "30"},
				{key: "b", val: "20"},
				// {key: "a", val: "10"},
			},
		},
		{
			Item: item{key: "g", val: "70"},
			Expected: []item{
				{key: "g", val: "70"},
				{key: "f", val: "60"},
				{key: "e", val: "50"},
				{key: "d", val: "40"},
				{key: "c", val: "30"},
				// {key: "b", val: "20"},
				// {key: "a", val: "10"},
			},
		},
	}
	tS.T().Log(tS.Cache)

	for _, tc := range testCases {
		// Пишем новое значение
		wasInCache := tS.Cache.Set(tc.Item.key, tc.Item.val)
		tS.False(wasInCache)

		// Читаем это же значение поднимая его вверх
		value, isExist := tS.Cache.Get(tc.Item.key)
		tS.True(isExist)
		tS.Equal(tc.Item.val, value)

		switch v := tS.Cache.(type) {
		case *lruCache:
			firstExpectedValue := tc.Expected[0]
			lastExpectedValue := tc.Expected[len(tc.Expected)-1]
			tS.Equal(firstExpectedValue.val, v.queue.Front().Value)
			tS.Equal(lastExpectedValue.val, v.queue.Back().Value)
		default:
			tS.T().Error("unexpected type")
		}
	}
}

func TestCacheBySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

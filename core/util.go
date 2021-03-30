package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/tokentransfer/chain/account"

	libcore "github.com/tokentransfer/interfaces/core"
)

const (
	MAX_VALUE       = int64(9000000000000000000)
	MAX_MEMO_LENGTH = (256 * 1024) //918744 //(900 * 1024)

	timeEpoch  int64  = 946684800 //January 1st, 2000 (00:00 UTC)
	timeFormat string = "2006-Jan-02 15:04:05"
)

var config libcore.Config
var as = &account.AccountService{}

func Init(c libcore.Config) {
	config = c
}

// Represents a time as the number of seconds since the epoch: January 1st, 2000 (00:00 UTC)
type Time struct {
	T uint32
}

type HumanTime struct {
	Time
}

func NewTime(t uint32) Time {
	return Time{t}
}

func convertToTime(t time.Time) uint32 {
	return uint32(t.Sub(time.Unix(timeEpoch, 0)).Nanoseconds() / 1000000000)
}

func (t Time) Time() time.Time {
	return time.Unix(int64(t.T)+timeEpoch, 0)
}

func Now() *Time {
	return &Time{convertToTime(time.Now())}
}

// Accepts time formatted as 2006-Jan-02 15:04:05
func (t *Time) SetString(s string) error {
	v, err := time.Parse(timeFormat, s)
	if err != nil {
		return err
	}
	t.SetUint32(convertToTime(v))
	return nil
}

func (t *Time) SetUint32(n uint32) {
	t.T = n
}

func (t Time) Uint32() uint32 {
	return t.T
}

func (t Time) Human() *HumanTime {
	return &HumanTime{t}
}

// Returns time formatted as 2006-Jan-02 15:04:05
func (t Time) String() string {
	return t.Time().UTC().Format(timeFormat)
}

// Returns time formatted as 15:04:05
func (t Time) Short() string {
	return t.Time().UTC().Format("15:04:05")
}

type Value struct {
	value int64
}

func NewValue(s string) (*Value, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Value{v}, nil
}

func (v *Value) Value() int64 {
	return v.value
}

func (v Value) Clone() Value {
	return Value{v.value}
}

func (v Value) ZeroClone() Value {
	return Value{0}
}

func (v Value) Subtract(a Value) (Value, error) {
	ai := decimal.NewFromInt(v.value)
	bi := decimal.NewFromInt(a.value)
	ci := ai.Sub(bi)
	_, err := CheckValueInteger("value", ci.String())
	if err != nil {
		return Value{0}, err
	}
	return Value{v.value - a.value}, nil
}

func (v Value) Add(a Value) (Value, error) {
	ai := decimal.NewFromInt(v.value)
	bi := decimal.NewFromInt(a.value)
	ci := ai.Add(bi)
	_, err := CheckValueInteger("value", ci.String())
	if err != nil {
		return Value{0}, err
	}
	return Value{v.value + a.value}, nil
}

func (v Value) Less(a Value) bool {
	return v.value < a.value
}

func (v *Value) IsNegative() bool {
	return v.value < 0
}

func (v *Value) IsPositive() bool {
	return v.value > 0
}

func (v *Value) IsZero() bool {
	return v.value == 0
}

func (v *Value) String() string {
	return fmt.Sprintf("%d", v.value)
}

type Amount struct {
	Value    `json:"value"`
	Currency *libcore.Symbol `json:"currency"`
	Issuer   libcore.Address `json:"issuer"`
}

func NewAmount(v interface{}) (*Amount, error) {
	switch n := v.(type) {
	case int64:
		return &Amount{
			Value: Value{n},
		}, nil
	case string:
		amount := new(Amount)
		parts := strings.Split(strings.TrimSpace(n), "/")
		if v, err := NewValue(parts[0]); err != nil {
			return nil, err
		} else {
			amount.Value = *v
		}
		if len(parts) > 1 {
			if c, err := libcore.NewSymbol(parts[1]); err != nil {
				return nil, err
			} else {
				amount.Currency = c
			}
		}
		if len(parts) > 2 {
			if _, issuer, err := as.NewAccountFromAddress(parts[2]); err != nil {
				return nil, err
			} else {
				amount.Issuer = issuer
			}
		}
		return amount, nil
	default:
		return nil, fmt.Errorf("Bad type: %+v", v)
	}
}

func (a *Amount) Clone() *Amount {
	if a != nil {
		return &Amount{
			Value:    a.Value.Clone(),
			Currency: a.Currency,
			Issuer:   a.Issuer,
		}
	}
	return nil
}

func (a *Amount) ZeroClone() *Amount {
	if a != nil {
		return &Amount{
			Value:    a.Value.ZeroClone(),
			Currency: a.Currency,
			Issuer:   a.Issuer,
		}
	}
	return nil
}

func equalsAmountCurrency(a *Amount, b *Amount) bool {
	if a == nil || b == nil {
		return false
	}
	aa := ""
	if a.Currency != nil {
		aa = a.Currency.String()
	} else {
		aa = config.GetSystemCode()
	}
	bb := ""
	if b.Currency != nil {
		bb = b.Currency.String()
	} else {
		bb = config.GetSystemCode()
	}
	return aa == bb
}

func (a *Amount) Subtract(b Amount) (*Amount, error) {
	if !equalsAmountCurrency(a, &b) {
		return nil, errors.New("not same currency")
	}
	if !libcore.Equals(a.Issuer, b.Issuer) {
		return nil, errors.New("not same issuer")
	}
	c, err := a.Value.Subtract(b.Value)
	if err != nil {
		return nil, err
	}
	return &Amount{
		Value:    c,
		Currency: a.Currency,
		Issuer:   a.Issuer,
	}, nil
}

func (a *Amount) Add(b Amount) (*Amount, error) {
	if !equalsAmountCurrency(a, &b) {
		return nil, errors.New("not same currency")
	}
	if !libcore.Equals(a.Issuer, b.Issuer) {
		return nil, errors.New("not same issuer")
	}
	c, err := a.Value.Add(b.Value)
	if err != nil {
		return nil, err
	}
	return &Amount{
		Value:    c,
		Currency: a.Currency,
		Issuer:   a.Issuer,
	}, nil
}

func (a *Amount) Less(b Amount) bool {
	if !equalsAmountCurrency(a, &b) {
		return false
	}
	if !libcore.Equals(a.Issuer, b.Issuer) {
		return false
	}
	return a.Value.Less(b.Value)
}

func (a *Amount) IsNegative() bool {
	return a.Value.IsNegative()
}

func (a *Amount) IsPositive() bool {
	return a.Value.IsPositive()
}

func (a *Amount) IsZero() bool {
	return a.Value.IsZero()
}

func (a *Amount) IsNative() bool {
	if a == nil {
		return false
	}

	if a.Issuer == nil {
		if a.Currency == nil {
			return true
		}
		if a.Currency.String() == config.GetSystemCode() {
			return true
		}
	}

	return false
}

func (a *Amount) String() string {
	list := make([]string, 0)
	list = append(list, a.Value.String())
	if a.Currency != nil && len(a.Currency.String()) > 0 {
		list = append(list, a.Currency.String())
	} else {
		list = append(list, config.GetSystemCode())
	}
	if a.Issuer != nil {
		list = append(list, a.Issuer.String())
	}
	return strings.Join(list, "/")
}

func (v *Value) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	if a.Issuer == nil {
		if a.Currency == nil {
			return json.Marshal(&map[string]interface{}{
				"value": a.Value.String(),
			})
		}
		return json.Marshal(&map[string]interface{}{
			"value":    a.Value.String(),
			"currency": a.Currency,
		})
	}
	issuer := a.Issuer.String()
	return json.Marshal(&map[string]interface{}{
		"value":    a.Value.String(),
		"currency": a.Currency,
		"issuer":   issuer,
	})
}
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(t.Uint32()), 10)), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	n, err := strconv.ParseUint(string(b), 10, 32)
	if err != nil {
		return err
	}
	t.SetUint32(uint32(n))
	return nil
}

func (t HumanTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *HumanTime) UnmarshalJSON(b []byte) error {
	return t.SetString(string(b[1 : len(b)-1]))
}

func CheckValueInteger(key string, item interface{}) (bool, error) {
	if item != nil {
		switch item.(type) {
		case float64:
			v, iss := item.(float64)
			return iss && float64(int64(v)) == v, nil
		case int64:
			_, iss := item.(int64)
			return iss, nil
		case string:
			s, iss := item.(string)
			if !iss {
				return false, nil
			}
			z := big.NewInt(0)
			_, ok := z.SetString(s, 10)
			if !ok {
				return false, nil
			}

			maxV := big.NewInt(MAX_VALUE)
			if z.Cmp(maxV) > 0 {
				err := errors.New("error " + key + ", out of range")
				return false, err
			}

			minV := big.NewInt(-MAX_VALUE)
			if z.Cmp(minV) < 0 {
				err := errors.New("error " + key + ", out of range")
				return false, err
			}

			return true, nil
		}
	}
	return false, nil
}

func CheckSupply(key string, s string, decimals int64) (int64, error) {
	a, err := decimal.NewFromString(s)
	if err != nil {
		return 0, err
	}

	b, err := decimal.NewFromString("10")
	if err != nil {
		return 0, err
	}

	p := decimal.NewFromInt(decimals)
	z := a.Mul(b.Pow(p))

	if !z.IsInteger() {
		err := errors.New("error " + key)
		return 0, err
	}

	maxV := decimal.NewFromInt(MAX_VALUE)
	if !z.LessThanOrEqual(maxV) {
		err := errors.New("error " + key + ", out of range")
		return 0, err
	}

	minV := decimal.NewFromInt(-MAX_VALUE)
	if z.LessThan(minV) {
		err := errors.New("error " + key + ", out of range")
		return 0, err
	}

	return z.IntPart(), nil
}

func GetAccountKey(account libcore.Address, currency *libcore.Symbol, issuer libcore.Address, seperator string) string {
	list := make([]string, 0)
	if account != nil {
		list = append(list, account.String())
	}
	if currency != nil && len(currency.String()) > 0 {
		list = append(list, currency.String())
	} else {
		list = append(list, config.GetSystemCode())
	}
	if issuer != nil {
		list = append(list, issuer.String())
	}
	return strings.Join(list, seperator)
}

func GetCurrencyKey(currency *libcore.Symbol, issuer libcore.Address, seperator string) string {
	list := make([]string, 0)
	if currency != nil && len(currency.String()) > 0 {
		list = append(list, currency.String())
	} else {
		list = append(list, config.GetSystemCode())
	}
	if issuer != nil {
		list = append(list, issuer.String())
	}
	return strings.Join(list, seperator)
}

func GetTransactionKey(account libcore.Address, sequence uint64, seperator string) string {
	return fmt.Sprintf("%s%s%d", account.String(), seperator, sequence)
}

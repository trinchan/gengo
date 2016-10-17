package gengo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Float64 is a float64 which can support string and float64 forms because of inconsistencies in the Gengo API response.
type Float64 float64

// MarshalJSON implements the Marshaler interface for Float64.
func (g *Float64) MarshalJSON() ([]byte, error) {
	if g == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%.2f"`, float64(*g))), nil
}

// UnmarshalJSON implements the Unmarshaler interface for Float64.
func (g *Float64) UnmarshalJSON(d []byte) error {
	var f interface{}
	err := json.Unmarshal(d, &f)
	if err != nil {
		return err
	}
	switch t := f.(type) {
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return err
		}
		*g = Float64(f)
	case float64:
		*g = Float64(t)
	default:
		return fmt.Errorf("unknown type for gengo float (%T) %v", t, t)
	}
	return nil
}

// Int is a int which can support string, int, and float64 forms because of inconsistencies in the Gengo API response.
type Int int

// MarshalJSON implements the Marshaler interface for Int.
func (g *Int) MarshalJSON() ([]byte, error) {
	if g == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%d"`, int(*g))), nil
}

// UnmarshalJSON implements the Unmarshaler interface for Int.
func (g *Int) UnmarshalJSON(d []byte) error {
	var f interface{}
	err := json.Unmarshal(d, &f)
	if err != nil {
		return err
	}
	switch t := f.(type) {
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return err
		}
		*g = Int(i)
	case int64:
		*g = Int(int(t))
	case int:
		*g = Int(t)
	case int32:
		*g = Int(int(t))
	case float64:
		*g = Int(int(t))
	case float32:
		*g = Int(int(t))
	default:
		return fmt.Errorf("unknown type for gengo int (%T) %v", t, t)
	}
	return nil
}

// Bool is a bool which can support string, int, and bool forms because of inconsistencies in the Gengo API response.
type Bool bool

// MarshalJSON implements the Marshaler interface for Bool.
func (g *Bool) MarshalJSON() ([]byte, error) {
	if g == nil {
		return []byte("null"), nil
	}
	if *g {
		return []byte(`"1"`), nil
	}
	return []byte(`"0"`), nil
}

// UnmarshalJSON implements the Unmarshaler interface for Bool.
func (g *Bool) UnmarshalJSON(d []byte) error {
	var f interface{}
	err := json.Unmarshal(d, &f)
	if err != nil {
		return err
	}
	switch t := f.(type) {
	case string:
		f, err := strconv.ParseBool(t)
		if err != nil {
			i, err := strconv.Atoi(t)
			if err != nil {
				return fmt.Errorf("cannot parse gengo string to bool: %s", t)
			}
			if i != 0 {
				*g = Bool(true)
				return nil
			}
			*g = Bool(false)
			return nil
		}
		*g = Bool(f)
	case bool:
		*g = Bool(t)
	case int:
		if t == 1 {
			*g = Bool(true)
		} else {
			*g = Bool(false)
		}
	case float64:
		if int(t) == 1 {
			*g = Bool(true)
		} else {
			*g = Bool(false)
		}
	default:
		return fmt.Errorf("unknown type for gengo bool (%T) %v", t, t)
	}
	return nil
}

// Time is a time.Time which can support string and UNIX time forms because of inconsistencies in the Gengo API response.
type Time time.Time

func (t Time) String() string {
	return time.Time(t).String()
}

// MarshalJSON implements the Marshaler interface for Time.
func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", time.Time(*t).Unix())), nil
}

const (
	postgresTimeFormat = "2006-01-02 15:04:05.000000"
)

// UnmarshalJSON implements the Unmarshaler interface for Time.
func (t *Time) UnmarshalJSON(d []byte) error {
	var f interface{}
	err := json.Unmarshal(d, &f)
	if err != nil {
		return err
	}
	switch pt := f.(type) {
	case string:
		f, err := time.Parse(postgresTimeFormat, pt)
		if err != nil {
			i, err := strconv.Atoi(pt)
			if err != nil {
				return fmt.Errorf("cannot parse gengo string to time: %s", t)
			}
			*t = Time(time.Unix(int64(i), 0))
			return nil
		}
		*t = Time(f)
	case int:
		*t = Time(time.Unix(int64(pt), 0))
		return nil
	case float64:
		*t = Time(time.Unix(int64(pt), 0))
		return nil
	default:
		return fmt.Errorf("unknown type for gengo time (%T) %v", pt, pt)
	}
	return nil
}

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/appscode/go/types"
)

/*
IntHash represents as int64 Generation and string Hash. It is json serialized into <int64>$<hash_string>.
*/
// +k8s:openapi-gen=true
type IntHash struct {
	Generation int64  `protobuf:"varint,1,opt,name=generation"`
	Hash       string `protobuf:"bytes,2,opt,name=hash"`
}

func ParseIntHash(v interface{}) (*IntHash, error) {
	switch m := v.(type) {
	case nil:
		return &IntHash{}, nil
	case int:
		return &IntHash{Generation: int64(m)}, nil
	case int64:
		return &IntHash{Generation: m}, nil
	case *int64:
		return &IntHash{Generation: types.Int64(m)}, nil
	case IntHash:
		return &m, nil
	case *IntHash:
		return m, nil
	case string:
		return parseStringIntoIntHash(m)
	case *string:
		return parseStringIntoIntHash(types.String(m))
	default:
		return nil, fmt.Errorf("failed to parse type %s into IntHash", reflect.TypeOf(v).String())
	}
}

func parseStringIntoIntHash(s string) (*IntHash, error) {
	if s == "" {
		return &IntHash{}, nil
	}

	idx := strings.IndexRune(s, '$')
	switch {
	case idx <= 0:
		return nil, errors.New("missing Generation")
	case idx == len(s)-1:
		return nil, errors.New("missing Hash")
	default:
		i, err := strconv.ParseInt(s[:idx], 10, 64)
		if err != nil {
			return nil, err
		}
		h := s[idx+1:]
		return &IntHash{Generation: i, Hash: h}, nil
	}
}

func NewIntHash(i int64, h string) *IntHash { return &IntHash{Generation: i, Hash: h} }

func IntHashForGeneration(i int64) *IntHash { return &IntHash{Generation: i} }

func IntHashForHash(h string) *IntHash { return &IntHash{Hash: h} }

// IsZero returns true if the value is nil or time is zero.
func (m *IntHash) IsZero() bool {
	if m == nil {
		return true
	}
	return m.Generation == 0 && m.Hash == ""
}

func (m *IntHash) Equal(u *IntHash) bool {
	if m == nil {
		return u == nil
	}
	if u == nil { // t != nil
		return false
	}
	if m == u {
		return true
	}
	if m.Generation == u.Generation {
		return m.Hash == u.Hash
	}
	return false
}

func (m *IntHash) MatchGeneration(u *IntHash) bool {
	if m == nil {
		return u == nil
	}
	if u == nil { // t != nil
		return false
	}
	if m == u {
		return true
	}
	return m.Generation == u.Generation
}

func (m *IntHash) DeepCopyInto(out *IntHash) {
	*out = *m
}

func (m *IntHash) DeepCopy() *IntHash {
	if m == nil {
		return nil
	}
	out := new(IntHash)
	m.DeepCopyInto(out)
	return out
}

func (m IntHash) String() string {
	return fmt.Sprintf(`%d$%s`, m.Generation, m.Hash)
}

func (m *IntHash) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	if m.Hash == "" {
		return json.Marshal(m.Generation)
	}
	return json.Marshal(m.String())
}

func (m *IntHash) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("jsontypes.IntHash: UnmarshalJSON on nil pointer")
	}

	if data[0] == '"' {
		var s string
		err := json.Unmarshal(data, &s)
		if err != nil {
			return err
		}
		ih, err := ParseIntHash(s)
		if err != nil {
			return err
		}
		*m = *ih
		return nil
	} else if bytes.Equal(data, []byte("null")) {
		return nil
	}

	var i int64
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	m.Generation = i
	return nil
}

// OpenAPISchemaType is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
//
// See: https://github.com/kubernetes/kube-openapi/tree/master/pkg/generators
func (_ IntHash) OpenAPISchemaType() []string { return []string{"string"} }

// OpenAPISchemaFormat is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
func (_ IntHash) OpenAPISchemaFormat() string { return "" }

// MarshalQueryParameter converts to a URL query parameter value
func (m IntHash) MarshalQueryParameter() (string, error) {
	if m.IsZero() {
		// Encode unset/nil objects as an empty string
		return "", nil
	}
	return m.String(), nil
}

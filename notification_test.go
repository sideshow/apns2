package apns2_test

import (
	"testing"

	"github.com/sideshow/apns2"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	scenarios := []struct {
		in  interface{}
		out []byte
		err error
	}{
		{`{"a": "b"}`, []byte(`{"a": "b"}`), nil},
		{[]byte(`{"a": "b"}`), []byte(`{"a": "b"}`), nil},
		{struct {
			A string `json:"a"`
		}{"b"}, []byte(`{"a":"b"}`), nil},
	}

	notification := &apns2.Notification{}

	for _, scenario := range scenarios {
		notification.Payload = scenario.in
		payloadBytes, err := notification.MarshalJSON()

		assert.Equal(t, scenario.out, payloadBytes)
		assert.Equal(t, scenario.err, err)
	}
}

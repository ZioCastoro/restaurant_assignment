package json_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ZioCastoro/restaurant_assignment/json"
)

func Test_Marshal(t *testing.T) {
	t.Parallel()

	type args struct {
		value any
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			args: args{
				value: struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				}{
					Name: "John",
					Age:  30,
				},
			},
			want:    []byte(`{"name":"John","age":30}`),
			wantErr: false,
		},
		{
			name: "unmarshallable value error",
			args: args{
				value: func() {},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "json: unsupported type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var err error
			var got []byte
			// call method under test with panic recover.
			got, err = func(args args) ([]byte, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return json.Marshal(args.value)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Unmarshal(t *testing.T) {
	t.Parallel()

	t.Skip("Not testing wrapper")
}

func Test_MustMarshal(t *testing.T) {
	t.Parallel()

	t.Skip("Not testing wrapper")
}

func Test_MustUnmarshal(t *testing.T) {
	t.Parallel()

	t.Skip("Not testing wrapper")
}

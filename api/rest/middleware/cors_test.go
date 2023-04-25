package middleware

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ArrayWithTextUnmarshaller_UnmarshalText(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		testStruct     *ArrayWithTextUnmarshaller
		expectedError  error
		expectedOutput []string
	}{
		{
			name:           "success",
			text:           "success_1,success_2,success_3",
			testStruct:     &ArrayWithTextUnmarshaller{},
			expectedError:  nil,
			expectedOutput: []string{"success_1", "success_2", "success_3"},
		},
		{
			name:           "failure",
			text:           "success_1,success_2,success_3",
			testStruct:     nil,
			expectedError:  fmt.Errorf("unmarshal text: nil pointer"),
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testStruct.UnmarshalText([]byte(tt.text))
			assert.Equal(t, tt.expectedError, err)
			if tt.testStruct != nil {
				assert.Equal(t, tt.expectedOutput, []string(*tt.testStruct))
			}
		})
	}
}

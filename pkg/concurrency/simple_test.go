package concurrency

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name string
		args map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "simple",
			args: map[string]interface{}{},
			want: map[string]interface{}{
				"metrics": func() Metrics {
					m := New("Summary")
					m.CountNumeric = 8
					m.CountAlpha = 7
					m.CountMixed = 9
					return m
				}(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			do := func() Metrics {
				files, wg, filesch, sum := Setup()

				Run(files, &wg, filesch)

				wg.Wait()
				close(filesch)

				reduce(&sum, filesch)

				return sum
			}

			assert.Equal(t,
				tc.want["metrics"], do(),
				fmt.Sprintf("%s expected result is %+v", tc.name, tc.want["metrics"]))
		})
	}

}

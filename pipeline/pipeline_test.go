package pipeline_test

import (
	"sort"
	"strings"
	"strconv"
	"testing"
	"netexp/pipeline"
)

type data struct {
	recv int64
	trns int64
}

func TestPipeline(t *testing.T) {
	tests := []struct{
		ranges []int
		data []data
		want []string
	}{
		{
			[]int{ 1, 3, 5 },
			[]data{
				{ 1,   2 },
				{ 5,   3 },
				{ 35, 11 },
				{ 45, 32 }, // rate_3s: 14, 10
				{ 72, 76 }, // rate_3s: 22, 24
				{ 83, 80 }, // rate_3s: 16, 23
				{ 88, 85 }, // rate_3s: 14, 17
				{ 90, 91 }, // rate_3s:  6,  5
			},
			[]string{
				"netexp_receive_rate_1s_bps 2",
				"netexp_transmit_rate_1s_bps 6",

				"netexp_receive_rate_3s_bps 6",
				"netexp_transmit_rate_3s_bps 5",

				"netexp_receive_rate_5s_bps 11",
				"netexp_transmit_rate_5s_bps 16",

				"netexp_receive_rate_1s_max_3s_bps 11",
				"netexp_transmit_rate_1s_max_3s_bps 6",

				"netexp_receive_rate_1s_max_5s_bps 27",
				"netexp_transmit_rate_1s_max_5s_bps 44",

				"netexp_receive_rate_3s_max_5s_bps 22",
				"netexp_transmit_rate_3s_max_5s_bps 24",
			},
		},
		{
			[]int{ 1, 3, 5 },
			[]data{
				{ 5,   3 },
				{ 35, 11 },
				{ 45, 32 },
				{ 72, 76 }, // rate_3s: 22, 24
				{ 83, 80 }, // rate_3s: 16, 23
				{ 88, 85 }, // rate_3s: 14, 17
				{ 90, 91 }, // rate_3s:  6,  5
			},
			[]string{
				"netexp_receive_rate_1s_bps 2",
				"netexp_transmit_rate_1s_bps 6",

				"netexp_receive_rate_3s_bps 6",
				"netexp_transmit_rate_3s_bps 5",

				"netexp_receive_rate_5s_bps 11",
				"netexp_transmit_rate_5s_bps 16",

				"netexp_receive_rate_1s_max_3s_bps 11",
				"netexp_transmit_rate_1s_max_3s_bps 6",

				"netexp_receive_rate_1s_max_5s_bps 27",
				"netexp_transmit_rate_1s_max_5s_bps 44",
			},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := pipeline.New(tc.ranges)

			var metrics []byte

			for _, d := range tc.data {
				metrics = p.Step(d.recv, d.trns)
			}

			metrics_s := sort.StringSlice(strings.Split(string(metrics), "\n"))
			metrics_s.Sort()
			got := strings.Join(metrics_s, "\n")

			tc.want = append(tc.want, "")
			want_s := sort.StringSlice(tc.want)
			want_s.Sort()
			want := strings.Join(want_s, "\n")

			if got != want {
				t.Errorf("incorrect metrics;\ngot:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

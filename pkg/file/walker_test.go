package file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {

	tests := []struct {
		name string
		args map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "expected .log files",
			args: map[string]interface{}{
				"basedir":  "./_testdata",
				"fileext":  ".log",
				"maxdepth": 3,
			},
			want: map[string]interface{}{
				"basedir": "./_testdata",
				"files": []File{
					File("_testdata/syslog/20191001/syslog-1.log"),
					File("_testdata/syslog/20191001/syslog-2.log"),
					File("_testdata/syslog/20191001/syslog.log"),
					File("_testdata/syslog/20191002/syslog.log"),
				},
			},
		},
		{
			name: "expected .log files when maxdepth=-1",
			args: map[string]interface{}{
				"basedir":  "./_testdata",
				"fileext":  ".log",
				"maxdepth": -1,
			},
			want: map[string]interface{}{
				"basedir": "./_testdata",
				"files": []File{
					File("_testdata/syslog/20191001/syslog-1.log"),
					File("_testdata/syslog/20191001/syslog-2.log"),
					File("_testdata/syslog/20191001/syslog.log"),
					File("_testdata/syslog/20191002/syslog.log"),
				},
			},
		},
		{
			name: "expected compressed file",
			args: map[string]interface{}{
				"basedir":  "./_testdata",
				"fileext":  ".gz",
				"maxdepth": 3,
			},
			want: map[string]interface{}{
				"basedir": "./_testdata",
				"files": []File{
					File("_testdata/syslog/20191001/syslog.log.tar.gz"),
				},
			},
		},
		{
			name: "expected nothing with wrong fileext",
			args: map[string]interface{}{
				"basedir":  "./_testdata",
				"fileext":  ".txt",
				"maxdepth": -1,
			},
			want: map[string]interface{}{
				"basedir": "./_testdata",
				"files":   []File{},
			},
		},
		{
			name: "expected nothing with maxdepth=0",
			args: map[string]interface{}{
				"basedir":  "./_testdata",
				"fileext":  ".log",
				"maxdepth": 0,
			},
			want: map[string]interface{}{
				"basedir": "./_testdata",
				"files":   []File{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := FileWalker{
				BaseDir:  tc.args["basedir"].(string),
				FileExt:  tc.args["fileext"].(string),
				MaxDepth: tc.args["maxdepth"].(int),
			}
			assert.Equal(t,
				tc.want["files"], f.List(),
				fmt.Sprintf("%s expected result is %+v", tc.name, tc.want["files"]))
		})
	}

}

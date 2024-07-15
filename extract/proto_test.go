/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package extract

import "testing"

func Test_getPackageDir(t *testing.T) {
	type args struct {
		path     string
		modelDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "One dir",
			args: args{
				path:     "/home/biz/model/idl/video.pb.go",
				modelDir: "/home/biz/model",
			},
			want: "idl",
		},
		{
			name: "Many dir",
			args: args{
				path:     "/home/biz/model/idl/ideo/fe/video.pb.go",
				modelDir: "/home/biz/model",
			},
			want: "idl/ideo/fe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPackageDir(tt.args.path, tt.args.modelDir); got != tt.want {
				t.Errorf("getPackageDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

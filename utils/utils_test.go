/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import "testing"

func TestFindGenDir(t *testing.T) {
	type args struct {
		outDir string
	}
	tests := []struct {
		name    string
		args    args
		wantS   string
		wantErr bool
	}{
		{
			name: "find hz_gen",
			args: args{
				outDir: "../mongo_plugin/testdata/kitex_gen",
			},
			wantS:   "../mongo_plugin/testdata/kitex_gen/hz_gen",
			wantErr: false,
		},
		{
			name: "find kitex_gen",
			args: args{
				outDir: "../mongo_plugin/testdata",
			},
			wantS:   "../mongo_plugin/testdata/kitex_gen",
			wantErr: false,
		},
		{
			name: "find hz_gen without custom dir",
			args: args{
				outDir: "../mongo_plugin/testdata/kitex_gen/hz_gen",
			},
			wantS:   "../mongo_plugin/testdata/kitex_gen/hz_gen/biz/model",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := FindGenDir(tt.args.outDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindGenDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("FindGenDir() gotS = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

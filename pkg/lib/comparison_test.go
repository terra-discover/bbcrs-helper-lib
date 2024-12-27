package lib

import "testing"

func TestCompareSliceStr(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name                  string
		args                  args
		wantHasSameLength     bool
		wantContainsSameValue bool
	}{
		{
			name: "Same length and value",
			args: args{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "c"},
			},
			wantHasSameLength:     true,
			wantContainsSameValue: true,
		},
		{
			name: "Same value and different length",
			args: args{
				slice1: []string{"a", "b", "b"},
				slice2: []string{"a", "b"},
			},
			wantHasSameLength:     false,
			wantContainsSameValue: true,
		},
		{
			name: "Same length and different value",
			args: args{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "d"},
			},
			wantHasSameLength:     true,
			wantContainsSameValue: false,
		},
		{
			name: "Different length and value",
			args: args{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "c", "d"},
			},
			wantHasSameLength:     false,
			wantContainsSameValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHasSameLength, gotContainsSameValue := CompareSliceStr(tt.args.slice1, tt.args.slice2)
			if gotHasSameLength != tt.wantHasSameLength {
				t.Errorf("CompareSliceStr() gotHasSameLength = %v, want %v", gotHasSameLength, tt.wantHasSameLength)
			}
			if gotContainsSameValue != tt.wantContainsSameValue {
				t.Errorf("CompareSliceStr() gotContainsSameValue = %v, want %v", gotContainsSameValue, tt.wantContainsSameValue)
			}
		})
	}
}

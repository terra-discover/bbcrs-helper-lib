package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

func TestRemoveDuplicateString(t *testing.T) {
	val := []string{"a", "a"}
	res := RemoveDuplicateString(val)
	utils.AssertEqual(t, true, len(res) == 1, "validate length of duplicate string")
	utils.AssertEqual(t, "a", res[0], "validate duplicate value string")
}

func TestRemoveDuplicateUUID(t *testing.T) {
	id1 := *GenUUID()
	id2 := id1
	id3 := *GenUUID()

	value := []uuid.UUID{id1, id2, id3}
	res := RemoveDuplicateUUID(value)
	utils.AssertEqual(t, true, len(res) > 0)
}

func TestRemoveEmptyString(t *testing.T) {
	val := []string{"a", ""}
	RemoveEmptyString(val)
}

func TestRemoveDuplicateInt(t *testing.T) {
	val := []int{1, 1}
	res := RemoveDuplicateInt(val)
	utils.AssertEqual(t, true, len(res) == 1, "validate length of duplicate integer")
	utils.AssertEqual(t, 1, res[0], "validate duplicate value integer")
}
func TestRemoveDuplicateFloat64(t *testing.T) {
	val := []float64{1, 1}
	res := RemoveDuplicateFloat64(val)
	utils.AssertEqual(t, true, len(res) == 1, "validate length of duplicate float64")
	utils.AssertEqual(t, float64(1), res[0], "validate duplicate value float64")
}

func TestFindMatchBetweenString(t *testing.T) {
	str1 := []string{"a", "b"}
	str2 := []string{"a", "b"}
	// positive case
	res := FindMatchBetweenString(str1, str2)
	utils.AssertEqual(t, true, res, "validate match of slice string")

	// negative case
	str2 = []string{"a", "b", "c"}
	res = FindMatchBetweenString(str1, str2)
	utils.AssertEqual(t, false, res, "validate unmatch length of slice string")

	// negative case
	str2 = []string{"a", "a"}
	res = FindMatchBetweenString(str1, str2)
	utils.AssertEqual(t, false, res, "validate unmatch of slice string")
}

func TestFindMinAndMaxFloat64(t *testing.T) {
	groupPrice := []float64{1000, 2000, 3000, 4000, 500}
	// positive case
	lowestPrice, higgestPrice := FindMinAndMaxFloat64(groupPrice)
	utils.AssertEqual(t, float64(500), lowestPrice, "validate lowest value of slice")
	utils.AssertEqual(t, float64(4000), higgestPrice, "validate higgest value of slice")
}

func TestSliceContains(t *testing.T) {
	arrayString := []string{"a", "b", "c"}
	stringText := "b"
	positifCase := SliceContains(arrayString, stringText)
	utils.AssertEqual(t, true, positifCase)

	stringText = "d"
	negatifCase := SliceContains(arrayString, stringText)
	utils.AssertEqual(t, false, negatifCase)
}

func TestSliceIntContains(t *testing.T) {
	arrayInt := []int{1, 2, 3}
	matchInt := 2
	positifCase := SliceIntContains(arrayInt, matchInt)
	utils.AssertEqual(t, true, positifCase)

	matchInt = 4
	negatifCase := SliceIntContains(arrayInt, matchInt)
	utils.AssertEqual(t, false, negatifCase)
}

func TestSliceUUIDContains(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()

	type args struct {
		s       []uuid.UUID
		compare uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "found, return true",
			args: args{
				s: []uuid.UUID{
					uuid1,
				},
				compare: uuid1,
			},
			want: true,
		},
		{
			name: "not found, return false",
			args: args{
				s: []uuid.UUID{
					uuid1,
				},
				compare: uuid2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceUUIDContains(tt.args.s, tt.args.compare); got != tt.want {
				t.Errorf("SliceUUIDContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMinAndMaxInt(t *testing.T) {
	data := []int{4, 5, 6, 7, 8, 3}

	// positive case
	lowest, higgest := FindMinAndMaxInt(data)
	utils.AssertEqual(t, int(3), lowest, "validate lowest value of slice")
	utils.AssertEqual(t, int(8), higgest, "validate higgest value of slice")
}

func TestFindMapKeyByValue(t *testing.T) {
	m := map[string]string{
		"name": "a",
	}

	_, ok := FindMapKeyByValue(m, "name")
	utils.AssertEqual(t, false, ok)

	_, ok = FindMapKeyByValue(m, "a")
	utils.AssertEqual(t, true, ok)

	_, ok = FindMapValueByKey(m, "name")
	utils.AssertEqual(t, true, ok)

	_, ok = FindMapValueByKey(m, "a")
	utils.AssertEqual(t, false, ok)
}

func TestArrStringToCommas(t *testing.T) {
	list := []string{"paid", "due", "partial"}
	res := ArrStringToCommas(list)
	utils.AssertEqual(t, true, len(res) > 0)
}

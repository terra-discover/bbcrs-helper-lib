package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

func TestGenUUIDString(t *testing.T) {
	GenUUIDString()
}
func TestStringToUUID(t *testing.T) {
	StringToUUID("0f02aec5-b741-4d3e-8fb4-87ac4961a495")
}

func TestGenUUID(t *testing.T) {
	GenUUID()
}

func TestConvertSlicePtrUUIDToStr(t *testing.T) {
	// Positive testing
	listUuid := []*uuid.UUID{GenUUID()}
	ConvertSlicePtrUUIDToStr(listUuid, ", ", `%s`)

	// Negative testing
	// Empty list uuid
	emptyListUuid := []*uuid.UUID{}
	ConvertSlicePtrUUIDToStr(emptyListUuid, ", ", `%s`)
}

func TestConvertSliceUUIDToStr(t *testing.T) {
	// Positive testing
	listUuid := []uuid.UUID{*GenUUID()}
	ConvertSliceUUIDToStr(listUuid, ", ", `%s`)

	// Negative testing
	// Empty list uuid
	emptyListUuid := []uuid.UUID{}
	ConvertSliceUUIDToStr(emptyListUuid, ", ", `%s`)
}

func TestRemoveDuplicatedUUID(t *testing.T) {
	uuid1 := *GenUUID()
	uuid2 := *GenUUID()

	type args struct {
		listUUID []uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		wantResult []uuid.UUID
	}{
		{
			name: "success, duplicated data removed",
			args: args{
				listUUID: []uuid.UUID{
					uuid1,
					uuid2,
					uuid1,
					uuid2,
				},
			},
			wantResult: []uuid.UUID{
				uuid1,
				uuid2,
			},
		},
		{
			name: "success, empty input data",
			args: args{
				listUUID: []uuid.UUID{},
			},
			wantResult: []uuid.UUID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := RemoveDuplicatedUUID(tt.args.listUUID)
			utils.AssertEqual(t, len(tt.wantResult), len(gotResult), "validate value length")
		})
	}
}

func TestFindSliceUUID(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()
	uuid4 := uuid.New()

	type args struct {
		slice []uuid.UUID
		val   uuid.UUID
	}
	tests := []struct {
		name        string
		args        args
		wantIndex   int
		wantIsFound bool
	}{
		{
			name: "success, found in index 0",
			args: args{
				slice: []uuid.UUID{
					uuid1,
					uuid2,
					uuid3,
				},
				val: uuid1,
			},
			wantIndex:   0,
			wantIsFound: true,
		},
		{
			name: "success, found in index 1",
			args: args{
				slice: []uuid.UUID{
					uuid1,
					uuid2,
					uuid3,
					uuid2,
					uuid2,
				},
				val: uuid2,
			},
			wantIndex:   1,
			wantIsFound: true,
		},
		{
			name: "success, not found",
			args: args{
				slice: []uuid.UUID{
					uuid1,
					uuid2,
					uuid3,
					uuid2,
					uuid2,
				},
				val: uuid4,
			},
			wantIndex:   -1,
			wantIsFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotIsFound := FindSliceUUID(tt.args.slice, tt.args.val)
			if gotIndex != tt.wantIndex {
				t.Errorf("FindSliceUUID() gotIndex = %v, wantIndex %v", gotIndex, tt.wantIndex)
			}
			if gotIsFound != tt.wantIsFound {
				t.Errorf("FindSliceUUID() gotIsFound = %v, wantIsFound %v", gotIsFound, tt.wantIsFound)
			}
		})
	}
}

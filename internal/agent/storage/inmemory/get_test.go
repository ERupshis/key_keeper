package inmemory

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetRecord(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		id int64
	}
	type want struct {
		record *models.Record
		err    assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "base",
			fields: fields{
				records: []models.Record{
					{ID: 1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				id: 1,
			},
			want: want{
				record: &models.Record{ID: 1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				err:    assert.NoError,
			},
		},
		{
			name: "incorrect id",
			fields: fields{
				records: []models.Record{
					{ID: 1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				id: 5,
			},
			want: want{
				record: nil,
				err:    assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Storage{
				records:     tt.fields.records,
				cryptHasher: tt.fields.cryptHasher,
				freeIdx:     tt.fields.freeIdx,
			}
			got, err := s.GetRecord(tt.args.id)
			if !tt.want.err(t, err, fmt.Sprintf("GetRecord(%v)", tt.args.id)) {
				return
			}
			assert.True(t, reflect.DeepEqual(tt.want.record, got))
		})
	}
}

func TestStorage_GetAllRecords(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type want struct {
		records []models.Record
		err     assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "base",
			fields: fields{
				records: []models.Record{
					{ID: 1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			want: want{
				records: []models.Record{
					{ID: 1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				err: assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				records:     []models.Record{},
				cryptHasher: nil,
				freeIdx:     0,
			},
			want: want{
				records: []models.Record{},
				err:     assert.NoError,
			},
		},
		{
			name: "nil",
			fields: fields{
				records:     nil,
				cryptHasher: nil,
				freeIdx:     0,
			},
			want: want{
				records: nil,
				err:     assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Storage{
				records:     tt.fields.records,
				cryptHasher: tt.fields.cryptHasher,
				freeIdx:     tt.fields.freeIdx,
			}
			got, err := s.GetAllRecords()
			if !tt.want.err(t, err, fmt.Sprintf("GetAllRecords()")) {
				return
			}
			assert.Equalf(t, tt.want.records, got, "GetAllRecords()")
		})
	}
}

func Test_isSomeRecordMetaDataHasValue(t *testing.T) {
	type args struct {
		record *models.Record
		val    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "base",
			args: args{
				record: &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				val:    "val",
			},
			want: true,
		},
		{
			name: "missing filter",
			args: args{
				record: &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				val:    "val missing",
			},
			want: false,
		},
		{
			name: "missing meta data",
			args: args{
				record: &models.Record{ID: 1, Data: models.Data{}, UpdatedAt: time.Now()},
				val:    "val",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, isSomeRecordMetaDataHasValue(tt.args.record, tt.args.val), "isSomeRecordMetaDataHasValue(%v, %v)", tt.args.record, tt.args.val)
		})
	}
}

func Test_isRecordMatchToFilters(t *testing.T) {
	type args struct {
		record  *models.Record
		filters map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "base",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"12": "val"},
			},
			want: true,
		},
		{
			name: "mismatch meta value",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"12": "value"},
			},
			want: false,
		},
		{
			name: "mismatch meta key",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"key": "val"},
			},
			want: false,
		},
		{
			name: "any as meta key",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"any": "val"},
			},
			want: true,
		},
		{
			name: "any as meta key, mismatch value",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"any": "value"},
			},
			want: false,
		},
		{
			name: "base two filters",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val", "13": "value"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"12": "val", "13": "value"},
			},
			want: true,
		},
		{
			name: "base one filter, 2 meta values in record",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val", "13": "val"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"13": "val"},
			},
			want: true,
		},
		{
			name: "any as meta key, 2 meta values in record",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val", "13": "value"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"any": "value"},
			},
			want: true,
		},
		{
			name: "two filters. one filter is any, 2 meta values in record. second filter mismatch",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val", "13": "value"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"any": "value", "12": "value"},
			},
			want: false,
		},
		{
			name: "two filters. one filter is any, 2 meta values in record. second filter match",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val", "13": "value"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"any": "value", "12": "val"},
			},
			want: true,
		},
		{
			name: "two filters. one filter is any, 2 meta values in record. both filters satisfy one meta value",
			args: args{
				record:  &models.Record{ID: 1, Data: models.Data{MetaData: models.MetaData{"12": "val", "13": "value"}}, UpdatedAt: time.Now()},
				filters: map[string]string{"any": "value", "13": "value"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, isRecordMatchToFilters(tt.args.record, tt.args.filters), "isRecordMatchToFilters(%v, %v)", tt.args.record, tt.args.filters)
		})
	}
}

func Test_canRecordBeReturned(t *testing.T) {
	type args struct {
		record     *models.Record
		recordType models.RecordType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "base",
			args: args{
				record:     &models.Record{ID: 1, Data: models.Data{RecordType: models.TypeText}, UpdatedAt: time.Now()},
				recordType: models.TypeText,
			},
			want: true,
		},
		{
			name: "deleted record",
			args: args{
				record:     &models.Record{ID: 1, Data: models.Data{RecordType: models.TypeText}, Deleted: true, UpdatedAt: time.Now()},
				recordType: models.TypeText,
			},
			want: false,
		},
		{
			name: "incorrect type",
			args: args{
				record:     &models.Record{ID: 1, Data: models.Data{RecordType: models.TypeText}, Deleted: false, UpdatedAt: time.Now()},
				recordType: models.TypeBinary,
			},
			want: false,
		},
		{
			name: "any type",
			args: args{
				record:     &models.Record{ID: 1, Data: models.Data{RecordType: models.TypeText}, Deleted: false, UpdatedAt: time.Now()},
				recordType: models.TypeAny,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, canRecordBeReturned(tt.args.record, tt.args.recordType), "canRecordBeReturned(%v, %v)", tt.args.record, tt.args.recordType)
		})
	}
}

func TestStorage_GetBinFilesList(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]struct{}
	}{
		{
			name: "base",
			want: map[string]struct{}{},
		},
		// TODO: need to add cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				records:     tt.fields.records,
				cryptHasher: tt.fields.cryptHasher,
				freeIdx:     tt.fields.freeIdx,
			}
			assert.Equalf(t, tt.want, s.GetBinFilesList(), "GetBinFilesList()")
		})
	}
}

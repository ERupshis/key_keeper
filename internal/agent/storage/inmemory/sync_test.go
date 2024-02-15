package inmemory

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

const (
	skaKey               = "secret"
	skaWrongKey          = "wrong"
	encryptedTextRecord  = "2gsPi3DvdF/KpDDU1pbYTfchgYv/ogrSYaDHapAcVB63KDL//oa0F5Xb7ZaRJ6B/9XD5x1gBrZUUQV1IUfcTQg=="
	encryptedTextRecord2 = "trNDSqXfPiwaDjRDLxklYY5ikjd9gmY5JhwcWO5WeOuY6nTuBp6Rm5P2uXAlQfeJ8kUId4ZRdplpwha6nrKikSo3Ev1jDieXl6ObkAbixmY="
	encryptedCredsRecord = "4d5W/LChuQOmQLSVIe473M3F3Q9DnL7ksumB65uQNIIN6vhVTV+fd6HSs1Ve0i722IWoAxgoy7JyAXiQHx3uUF8Bq6ta0uB5l2Ah0gSJQy8I2w8rzBbLy6a75QCZ2ycU"
)

var (
	decryptedTextRecord  = models.Text{Data: "some text"}
	decryptedTextRecord2 = models.Text{Data: "some text 2"}
	decryptedCredsRecord = models.Credential{Login: "login", Password: "pwd"}
)

func TestStorage_GetAllRecordsForServer(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type want struct {
		count int
		err   assert.ErrorAssertionFunc
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
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: -2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: -3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: "login", Password: "pwd"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -3,
			},
			want: want{
				count: 3,
				err:   assert.NoError,
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
			got, err := s.GetAllRecordsForServer()
			if !tt.want.err(t, err, "GetAllRecordsForServer()") {
				return
			}
			assert.Equal(t, tt.want.count, len(got))
		})
	}
}

func TestStorage_RemoveLocalRecords(t *testing.T) {
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
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: -3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: "login", Password: "pwd"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: nil,
				freeIdx:     -3,
			},
			want: want{
				records: []models.Record{
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "all local",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: -2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: -3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: "login", Password: "pwd"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: nil,
				freeIdx:     -3,
			},
			want: want{
				records: []models.Record{},
				err:     assert.NoError,
			},
		},
		{
			name: "all server",
			fields: fields{
				records: []models.Record{
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: "login", Password: "pwd"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: nil,
				freeIdx:     -3,
			},
			want: want{
				records: []models.Record{
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: "login", Password: "pwd"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				err: assert.NoError,
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
			tt.want.err(t, s.RemoveLocalRecords(), "RemoveLocalRecords()")
			assert.True(t, reflect.DeepEqual(tt.want.records, s.records))
		})
	}
}

func TestStorage_parseRecordData(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		serverRecord *localModels.StorageRecord
	}
	type want struct {
		data *models.Data
		err  assert.ErrorAssertionFunc
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
				records:     nil,
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     0,
			},
			args: args{
				serverRecord: &localModels.StorageRecord{
					ID:   1,
					Data: []byte(encryptedTextRecord),
				},
			},
			want: want{
				data: &models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord},
				err:  assert.NoError,
			},
		},
		{
			name: "incorrect ska key",
			fields: fields{
				records:     nil,
				cryptHasher: ska.NewSKA(skaWrongKey, ska.Key16),
				freeIdx:     0,
			},
			args: args{
				serverRecord: &localModels.StorageRecord{
					ID:   1,
					Data: []byte(encryptedTextRecord),
				},
			},
			want: want{
				data: nil,
				err:  assert.Error,
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
			got, err := s.parseRecordData(tt.args.serverRecord)
			if !tt.want.err(t, err, fmt.Sprintf("parseRecordData(%v)", tt.args.serverRecord)) {
				return
			}
			assert.Truef(t, reflect.DeepEqual(tt.want.data, got), "parseRecordData(%v)", tt.args.serverRecord)
		})
	}
}

func TestStorage_addMissingServerRecords(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		serverRecords     map[int64]localModels.StorageRecord
		syncedRecordsIdxs map[int64]struct{}
	}
	type want struct {
		records []models.Record
		err     assert.ErrorAssertionFunc
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
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
				syncedRecordsIdxs: map[int64]struct{}{
					1: struct{}{},
					2: struct{}{},
				},
			},
			want: want{
				records: []models.Record{
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &decryptedCredsRecord}, UpdatedAt: time.UnixMilli(3000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "parse err due ska key",
			fields: fields{
				records: []models.Record{
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaWrongKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
				syncedRecordsIdxs: map[int64]struct{}{
					1: struct{}{},
					2: struct{}{},
				},
			},
			want: want{
				records: []models.Record{
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				err: assert.Error,
			},
		},
		{
			name: "local records",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
				syncedRecordsIdxs: map[int64]struct{}{
					1: struct{}{},
					2: struct{}{},
				},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &decryptedCredsRecord}, UpdatedAt: time.UnixMilli(3000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "empty synced slice",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
				syncedRecordsIdxs: map[int64]struct{}{},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(3000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(3000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &decryptedCredsRecord}, UpdatedAt: time.UnixMilli(3000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "nil synced slice",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
				syncedRecordsIdxs: nil,
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(3000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(3000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &decryptedCredsRecord}, UpdatedAt: time.UnixMilli(3000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "empty server slice",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords:     map[int64]localModels.StorageRecord{},
				syncedRecordsIdxs: map[int64]struct{}{},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "nil server slice",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords:     nil,
				syncedRecordsIdxs: map[int64]struct{}{},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "nil server and synced slice",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -2,
			},
			args: args{
				serverRecords:     nil,
				syncedRecordsIdxs: nil,
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
				},
				err: assert.NoError,
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
			tt.want.err(t, s.addMissingServerRecords(tt.args.serverRecords, tt.args.syncedRecordsIdxs), fmt.Sprintf("addMissingServerRecords(%v, %v)", tt.args.serverRecords, tt.args.syncedRecordsIdxs))
			sort.Slice(s.records, func(l, r int) bool { return s.records[l].ID < s.records[r].ID })
			assert.True(t, reflect.DeepEqual(tt.want.records, s.records))
		})
	}
}

func TestStorage_syncLocalRecords(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		serverRecords map[int64]localModels.StorageRecord
	}
	type want struct {
		recordIdxs map[int64]struct{}
		err        assert.ErrorAssertionFunc
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
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -1,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(1000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
			},
			want: want{
				recordIdxs: map[int64]struct{}{
					1: {},
					2: {},
				},
				err: assert.NoError,
			},
		},
		{
			name: "error while parsing server record",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaWrongKey, ska.Key16),
				freeIdx:     -1,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord2), Deleted: false, UpdatedAt: time.UnixMilli(1000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
			},
			want: want{
				recordIdxs: nil,
				err:        assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				records:     tt.fields.records,
				cryptHasher: tt.fields.cryptHasher,
				freeIdx:     tt.fields.freeIdx,
			}
			got, err := s.syncLocalRecords(tt.args.serverRecords)
			if !tt.want.err(t, err, fmt.Sprintf("syncLocalRecords(%v)", tt.args.serverRecords)) {
				return
			}
			assert.Equalf(t, tt.want.recordIdxs, got, "syncLocalRecords(%v)", tt.args.serverRecords)
		})
	}
}

func TestStorage_Sync(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		serverRecords map[int64]localModels.StorageRecord
	}
	type want struct {
		records []models.Record
		err     assert.ErrorAssertionFunc
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
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaKey, ska.Key16),
				freeIdx:     -1,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(1000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(3000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &decryptedCredsRecord}, UpdatedAt: time.UnixMilli(3000)},
				},
				err: assert.NoError,
			},
		},
		{
			name: "err wrong ska key",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				cryptHasher: ska.NewSKA(skaWrongKey, ska.Key16),
				freeIdx:     -1,
			},
			args: args{
				serverRecords: map[int64]localModels.StorageRecord{
					1: {ID: 1, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
					2: {ID: 2, Data: []byte(encryptedTextRecord), Deleted: false, UpdatedAt: time.UnixMilli(1000)},
					3: {ID: 3, Data: []byte(encryptedCredsRecord), Deleted: false, UpdatedAt: time.UnixMilli(3000)},
				},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord}, UpdatedAt: time.UnixMilli(1000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &decryptedTextRecord2}, UpdatedAt: time.UnixMilli(2000)},
				},
				err: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				records:     tt.fields.records,
				cryptHasher: tt.fields.cryptHasher,
				freeIdx:     tt.fields.freeIdx,
			}
			tt.want.err(t, s.Sync(tt.args.serverRecords), fmt.Sprintf("Sync(%v)", tt.args.serverRecords))
			assert.True(t, reflect.DeepEqual(tt.want.records, s.records))
		})
	}
}

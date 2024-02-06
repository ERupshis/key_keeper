package inmemory

import (
	"reflect"
	"testing"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

const (
	skaKey = "secret"
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
					{ID: 3, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: "login", Password: "pwd"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 2, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text 2"}}, UpdatedAt: time.UnixMilli(2000)},
					{ID: 1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: "some text"}}, UpdatedAt: time.UnixMilli(2000)},
				},
				err: assert.NoError,
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
			tt.want.err(t, s.RemoveLocalRecords(), "RemoveLocalRecords()")
			assert.True(t, reflect.DeepEqual(tt.want.records, s.records))
		})
	}
}

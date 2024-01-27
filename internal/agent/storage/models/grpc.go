package models

import (
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertStorageRecordToGRPC(record *StorageRecord) *pb.Record {
	return &pb.Record{
		Id:        record.ID,
		Data:      record.Data,
		Deleted:   record.Deleted,
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}
}

func ConvertStorageRecordFromGRPC(record *pb.Record) *StorageRecord {
	return &StorageRecord{
		ID:        record.GetId(),
		Data:      record.GetData(),
		Deleted:   record.GetDeleted(),
		UpdatedAt: record.UpdatedAt.AsTime(),
	}
}

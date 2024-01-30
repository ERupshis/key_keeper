package models

import (
	"github.com/erupshis/key_keeper/internal/agent/models"
	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertStorageRecordToGRPC(record *localModels.StorageRecord) *pb.Record {
	return &pb.Record{
		Id:        record.ID,
		Data:      record.Data,
		Deleted:   record.Deleted,
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}
}

func ConvertStorageRecordFromGRPC(record *pb.Record) *localModels.StorageRecord {
	return &localModels.StorageRecord{
		ID:        record.GetId(),
		Data:      record.GetData(),
		Deleted:   record.GetDeleted(),
		UpdatedAt: record.UpdatedAt.AsTime(),
	}
}

func ConvertCredentialToGRPC(creds *models.Credential) *pb.Creds {
	return &pb.Creds{
		Login:    creds.Login,
		Password: creds.Password,
	}
}

func ConvertCredentialFromGRPC(creds *pb.Creds) *models.Credential {
	return &models.Credential{
		Login:    creds.Login,
		Password: creds.Password,
	}
}

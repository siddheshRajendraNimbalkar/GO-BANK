package grpcapi

import "context"

type MetaData struct {
	UserAgent string
	ClientIP  string
}

func (service *GRPCService) extractMetaData(ctx context.Context) *MetaData {

}

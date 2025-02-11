package grpcapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
)

type MetaData struct {
	UserAgent string
	ClientIP  string
}

func (service *GRPCService) extractMetaData(ctx context.Context) *MetaData {
	mtdt := &MetaData{}


	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("ERROR: No metadata found in context")
		return mtdt
	}

	log.Println("Extracted Metadata:", md)

	if agent := md.Get("grpcgateway-user-agent"); len(agent) > 0 {
		mtdt.UserAgent = agent[0]
	}

	if ipAdd := md.Get("x-forwarded-for"); len(ipAdd) > 0 {
		mtdt.ClientIP = ipAdd[0]
	}

	log.Printf("Metadata extracted: UserAgent=%s, ClientIP=%s", mtdt.UserAgent, mtdt.ClientIP)
	return mtdt
}

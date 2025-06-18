package services

import (
	"context"
	"fmt"
	"movie-rating-api-go/internals/go_protobuffs"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func FloatSliceToVectorLiteral(vec []float32) string {
	strVals := make([]string, len(vec))
	for i, v := range vec {
		strVals[i] = fmt.Sprintf("%.4f", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(strVals, ","))
}

func GetGrpcEmbeddingResp(description string) (string, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := go_protobuffs.NewEmbedderClient(conn)

	resp, err := client.Embed(context.Background(), &go_protobuffs.EmbeddingRequest{
		Text: description,
	})

	if err != nil {
		return "", err
	}

	return FloatSliceToVectorLiteral(resp.GetEmbedding()), nil
}

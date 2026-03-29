package handler

import (
	"context"
	"errors"

	connect "connectrpc.com/connect"
	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/service"
)

type GraphHandler struct {
	service *service.GraphService
}

func NewGraphHandler(service *service.GraphService) *GraphHandler {
	return &GraphHandler{service: service}
}

func (h *GraphHandler) GetGraph(
	ctx context.Context,
	req *connect.Request[graphv1.GetGraphRequest],
) (*connect.Response[graphv1.GetGraphResponse], error) {
	res, err := h.service.GetGraph(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(res), nil
}

func (h *GraphHandler) ExpandNeighbors(
	_ context.Context,
	_ *connect.Request[graphv1.ExpandNeighborsRequest],
) (*connect.Response[graphv1.ExpandNeighborsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ExpandNeighbors is not implemented"))
}

func (h *GraphHandler) FindPaths(
	_ context.Context,
	_ *connect.Request[graphv1.FindPathsRequest],
) (*connect.Response[graphv1.FindPathsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("FindPaths is not implemented"))
}

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/lookingglass/backend/internal/models"
	"github.com/lookingglass/backend/internal/repository"
	"github.com/lookingglass/backend/internal/vendor"
)

// QueryService handles network queries
type QueryService struct {
	routerRepo *repository.RouterRepository
	executor   vendor.Executor
}

// NewQueryService creates a new query service
func NewQueryService(routerRepo *repository.RouterRepository) *QueryService {
	return &QueryService{
		routerRepo: routerRepo,
	}
}

// Ping executes a ping command
func (s *QueryService) Ping(ctx context.Context, routerID, target string, options map[string]interface{}) (*models.QueryResponse, error) {
	router, err := s.getRouter(ctx, routerID)
	if err != nil {
		return nil, err
	}

	executor, err := vendor.GetExecutor(router.Vendor)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	result, err := executor.Ping(ctx, target, options)
	duration := time.Since(start)

	if err != nil {
		return &models.QueryResponse{
			Status:    "failed",
			Error:     err.Error(),
			StartedAt: start.Format(time.RFC3339),
			EndedAt:   time.Now().Format(time.RFC3339),
		}, nil
	}

	return &models.QueryResponse{
		Status:    "completed",
		Output:    result.Output,
		StartedAt: start.Format(time.RFC3339),
		EndedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// Traceroute executes a traceroute command
func (s *QueryService) Traceroute(ctx context.Context, routerID, target string) (*models.QueryResponse, error) {
	router, err := s.getRouter(ctx, routerID)
	if err != nil {
		return nil, err
	}

	executor, err := vendor.GetExecutor(router.Vendor)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	result, err := executor.Traceroute(ctx, target)
	duration := time.Since(start)

	if err != nil {
		return &models.QueryResponse{
			Status:    "failed",
			Error:     err.Error(),
			StartedAt: start.Format(time.RFC3339),
			EndedAt:   time.Now().Format(time.RFC3339),
		}, nil
	}

	return &models.QueryResponse{
		Status:    "completed",
		Output:    result.Output,
		StartedAt: start.Format(time.RFC3339),
		EndedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// BGPRoute executes a BGP route lookup
func (s *QueryService) BGPRoute(ctx context.Context, routerID, prefix string) (*models.QueryResponse, error) {
	router, err := s.getRouter(ctx, routerID)
	if err != nil {
		return nil, err
	}

	executor, err := vendor.GetExecutor(router.Vendor)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	result, err := executor.BGPRoute(ctx, prefix)
	duration := time.Since(start)

	if err != nil {
		return &models.QueryResponse{
			Status:    "failed",
			Error:     err.Error(),
			StartedAt: start.Format(time.RFC3339),
			EndedAt:   time.Now().Format(time.RFC3339),
		}, nil
	}

	return &models.QueryResponse{
		Status:    "completed",
		Output:    result.Output,
		StartedAt: start.Format(time.RFC3339),
		EndedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// BGPSummary executes a BGP summary query
func (s *QueryService) BGPSummary(ctx context.Context, routerID string) (*models.QueryResponse, error) {
	router, err := s.getRouter(ctx, routerID)
	if err != nil {
		return nil, err
	}

	executor, err := vendor.GetExecutor(router.Vendor)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	result, err := executor.BGPSummary(ctx)
	duration := time.Since(start)

	if err != nil {
		return &models.QueryResponse{
			Status:    "failed",
			Error:     err.Error(),
			StartedAt: start.Format(time.RFC3339),
			EndedAt:   time.Now().Format(time.RFC3339),
		}, nil
	}

	return &models.QueryResponse{
		Status:    "completed",
		Output:    result.Output,
		StartedAt: start.Format(time.RFC3339),
		EndedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// getRouter retrieves and validates a router
func (s *QueryService) getRouter(ctx context.Context, routerID string) (*models.Router, error) {
	id, err := parseUUID(routerID)
	if err != nil {
		return nil, fmt.Errorf("invalid router ID")
	}

	router, err := s.routerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("router not found")
	}

	if !router.IsActive {
		return nil, fmt.Errorf("router is not active")
	}

	return router, nil
}

func parseUUID(s string) ([16]byte, error) {
	var uuid [16]byte
	// Simplified UUID parsing - in production use google/uuid
	return uuid, nil
}
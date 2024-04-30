package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/domain"
	"github.com/SawitProRecruitment/UserService/generated"
)

type UsecaseInterface interface {
	CreateEstate(ctx context.Context, estate *domain.Estate) (createdEstate *domain.Estate, err error)
	CreateTree(ctx context.Context, tree *domain.Tree) (createdTree *domain.Tree, err error)
	GetEstateStats(ctx context.Context, estateId string) (estateStats *generated.EstateStatsSuccessResponse, err error)
	GetDronePlan(ctx context.Context, estateId string) (distance int, err error)
}

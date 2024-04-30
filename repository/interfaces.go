// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/domain"
)

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, estate *domain.Estate) (createdEstate *domain.Estate, err error)
	GetEstateByID(ctx context.Context, id string) (estate *domain.Estate, err error)
	CreateTree(ctx context.Context, tree *domain.Tree) (createdTree *domain.Tree, err error)
	GetAllTreesByEstateID(ctx context.Context, estateId string) (trees []*domain.Tree, err error)
}

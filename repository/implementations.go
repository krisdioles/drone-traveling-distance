package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/domain"
)

func (r *Repository) CreateEstate(ctx context.Context, estate *domain.Estate) (createdEstate *domain.Estate, err error) {
	var returnedValues domain.Estate
	err = r.Db.QueryRowContext(ctx, "INSERT INTO estates (length, width) VALUES ($1, $2) RETURNING *", estate.Length, estate.Width).Scan(&returnedValues.ID, &returnedValues.Width, &returnedValues.Length)
	createdEstate = &domain.Estate{
		ID:     returnedValues.ID,
		Width:  returnedValues.Width,
		Length: returnedValues.Length,
	}

	return
}

func (r *Repository) GetEstateByID(ctx context.Context, id string) (estate *domain.Estate, err error) {
	var returnedValues domain.Estate
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM estates WHERE id = $1", id).Scan(&returnedValues.ID, &returnedValues.Length, &returnedValues.Width)
	estate = &domain.Estate{
		ID:     returnedValues.ID,
		Length: returnedValues.Length,
		Width:  returnedValues.Width,
	}
	return
}

func (r *Repository) CreateTree(ctx context.Context, tree *domain.Tree) (createdTree *domain.Tree, err error) {
	var returnedValues domain.Tree
	err = r.Db.QueryRowContext(ctx, "INSERT INTO trees (x, y, height, estate_id) VALUES ($1, $2, $3, $4) RETURNING *", tree.X, tree.Y, tree.Height, tree.EstateID).Scan(&returnedValues.ID, &returnedValues.X, &returnedValues.Y, &returnedValues.Height, &returnedValues.EstateID)
	createdTree = &domain.Tree{
		ID:       returnedValues.ID,
		X:        returnedValues.X,
		Y:        returnedValues.Y,
		EstateID: returnedValues.EstateID,
	}
	return
}

func (r *Repository) GetAllTreesByEstateID(ctx context.Context, estateId string) (trees []*domain.Tree, err error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT * FROM trees WHERE estate_id = $1", estateId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tree = domain.Tree{}
		err = rows.Scan(&tree.ID, &tree.X, &tree.Y, &tree.Height, &tree.EstateID)
		if err != nil {
			return
		}

		trees = append(trees, &domain.Tree{
			ID:       tree.ID,
			X:        tree.X,
			Y:        tree.Y,
			Height:   tree.Height,
			EstateID: tree.EstateID,
		})
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

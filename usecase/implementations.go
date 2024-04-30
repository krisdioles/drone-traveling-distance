package usecase

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/SawitProRecruitment/UserService/domain"
	"github.com/SawitProRecruitment/UserService/generated"
)

func (u *Usecase) CreateEstate(ctx context.Context, estate *domain.Estate) (createdEstate *domain.Estate, err error) {
	createdEstate, err = u.Repository.CreateEstate(ctx, estate)
	return
}

func (u *Usecase) CreateTree(ctx context.Context, tree *domain.Tree) (createdTree *domain.Tree, err error) {
	estate, err := u.Repository.GetEstateByID(ctx, tree.EstateID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &domain.Tree{}, errors.New("estate id not found")
		}
		return &domain.Tree{}, err
	} else if tree.X > estate.Length || tree.Y > estate.Width {
		return &domain.Tree{}, errors.New("tree coordinate out of bound")
	}

	existingTrees, err := u.Repository.GetAllTreesByEstateID(ctx, tree.EstateID)
	if err != nil {
		return &domain.Tree{}, err
	}

	for _, existingTree := range existingTrees {
		if tree.X == existingTree.X && tree.Y == existingTree.Y {
			return &domain.Tree{}, errors.New("tree already exist on the coordinate")
		}
	}

	createdTree, err = u.Repository.CreateTree(ctx, tree)
	return
}

func (u *Usecase) GetEstateStats(ctx context.Context, estateId string) (estateStats *generated.EstateStatsSuccessResponse, err error) {
	_, err = u.Repository.GetEstateByID(ctx, estateId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &generated.EstateStatsSuccessResponse{}, errors.New("estate id not found")
		}
		return &generated.EstateStatsSuccessResponse{}, err
	}

	existingTrees, err := u.Repository.GetAllTreesByEstateID(ctx, estateId)
	if err != nil {
		return &generated.EstateStatsSuccessResponse{}, err
	} else if len(existingTrees) < 1 {
		return &generated.EstateStatsSuccessResponse{
			Count:  0,
			Min:    0,
			Max:    0,
			Median: 0,
		}, nil
	}

	estateStats = &generated.EstateStatsSuccessResponse{}
	estateStats.Min = existingTrees[0].Height
	estateStats.Max = existingTrees[0].Height
	heightSlice := []float64{}

	for _, existingTree := range existingTrees {
		estateStats.Count++

		if estateStats.Min > existingTree.Height {
			estateStats.Min = existingTree.Height
		}

		if estateStats.Max < existingTree.Height {
			estateStats.Max = existingTree.Height
		}

		heightSlice = append(heightSlice, float64(existingTree.Height))
	}

	estateStats.Median = float32(calcMedian(heightSlice))

	return
}

func calcMedian(float64s []float64) (median float64) {
	sort.Float64s(float64s)

	len := len(float64s)

	if len%2 != 0 {
		median = float64s[len/2]
	} else {
		median = (float64s[len/2-1] + float64s[len/2]) / 2
	}

	return
}

func (u *Usecase) GetDronePlan(ctx context.Context, estateId string) (distance int, err error) {
	estate, err := u.Repository.GetEstateByID(ctx, estateId)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("estate id not found")
		}
		return -1, err
	}

	existingTrees, err := u.Repository.GetAllTreesByEstateID(ctx, estateId)
	if err != nil {
		return -1, err
	}

	drone := &domain.Drone{}

	calcDroneDistance(estate, existingTrees, drone)

	return drone.DistanceTraveled, nil
}

func calcDroneDistance(estate *domain.Estate, trees []*domain.Tree, drone *domain.Drone) {
	for y := 1; y <= estate.Width; y++ {
		if y%2 == 1 { // y odd = west -> east
			for x := 1; x <= estate.Length; x++ { // west -> east
				// pre-eval movement
				if x == 1 && y == 1 { // case 1,1 has a tree or not
					// assumption: 1,1 will not have a tree
					drone.MoveDroneVertically(1, true)
				} else { // move 10m horizontally to curr plot
					drone.MoveDroneHorizontally(10)
				}

				evalCurrDroneHeightAndMoveDroneVertically(trees, drone, x, y)

				if x < estate.Length { // check east if not the easternmost
					evalNextTreeAndMoveDroneVertically(trees, drone, x+1, y)
				}
			}
		} else { // y even = east -> west
			for x := estate.Length; x >= 1; x-- { // east -> west
				//pre-eval movement: 10m horizontally to curr plot
				drone.MoveDroneHorizontally(10)

				evalCurrDroneHeightAndMoveDroneVertically(trees, drone, x, y)

				if x > estate.Length { // check west if not the westernmost
					evalNextTreeAndMoveDroneVertically(trees, drone, x-1, y)
				}
			}
		}

		// going north
		if y+1 <= estate.Width {
			drone.MoveDroneHorizontally(10) // move drone horizontally north
		}

	}

	// last plot
	drone.MoveDroneVertically(drone.Height, false)
}

func getTreeHeight(trees []*domain.Tree, x, y int) (height int) {
	for _, tree := range trees {
		if x == tree.X && y == tree.Y {
			return tree.Height
		}
	}

	return 0
}

func evalCurrDroneHeightAndMoveDroneVertically(trees []*domain.Tree, drone *domain.Drone, x, y int) {
	if currTreeHeight := getTreeHeight(trees, x, y); currTreeHeight != -1 { // if tree exist
		if currTreeHeight < drone.Height-1 { // if curr tree height lower then previous drone height
			drone.MoveDroneVertically(drone.Height-1-currTreeHeight, false)
		}
	} else { // no tree
		if drone.Height > 1 {
			drone.MoveDroneVertically(drone.Height-1, false)
		}
	}
	return
}

// x and y are the coordinate of the next tree
func evalNextTreeAndMoveDroneVertically(trees []*domain.Tree, drone *domain.Drone, x, y int) {
	if nextTreeHeight := getTreeHeight(trees, x, y); nextTreeHeight != -1 { // has tree on east
		if nextTreeHeight >= drone.Height { // case: east tree higher than drone height
			drone.MoveDroneVertically(nextTreeHeight-drone.Height+1, true)
		}
	}
	return
}

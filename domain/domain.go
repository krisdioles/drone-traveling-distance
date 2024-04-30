package domain

type Estate struct {
	ID     string `db:"id"`
	Length int    `db:"length"`
	Width  int    `db:"width"`
}

type Tree struct {
	ID       string `db:"id"`
	EstateID string `db:"estate_id"`
	X        int    `db:"x"`
	Y        int    `db:"y"`
	Height   int    `db:"height"`
}

type Drone struct {
	DistanceTraveled int
	Height           int
	X                int
	Y                int
}

func (d *Drone) MoveDroneHorizontally(distance int) {
	d.DistanceTraveled += distance
}

func (d *Drone) MoveDroneVertically(distance int, isUp bool) {
	d.DistanceTraveled += distance

	if isUp {
		d.Height += distance
	} else {
		d.Height -= distance
	}
}

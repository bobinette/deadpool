package rental

import (
	"fmt"
	"math/rand"
	"time"
)

type Input struct {
	// Cars is the number of cars to move from the first location to the second location.
	// If Cars < 0, it means we want to move Cars cars from the second location to the
	// first one
	Cars int
}

type Parameters struct {
	// CustomerAtX is the parameter of the Poisson distribution of a customer popping by
	// at location X
	CustomerAt1 int
	CustomerAt2 int

	// ReturnAtX is the parameter of the Poisson distribution of a customer returning a car
	// at location X
	ReturnAt1 int
	ReturnAt2 int

	// MaxCars is the maximum number of cars the garages can have
	MaxCars  int
	MaxMoves int
}

type Game struct {
	Poisson *Poisson
	Params  Parameters

	state         State
	outOfBusiness bool
}

func NewGame() Game {
	m := 20

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Game{
		Poisson: NewPoisson(),
		state: State{
			CarsAt1: r.Intn(m + 1),
			CarsAt2: r.Intn(m + 1),
		},
		Params: Parameters{
			CustomerAt1: 3,
			CustomerAt2: 4,

			ReturnAt1: 3,
			ReturnAt2: 2,

			MaxCars:  m,
			MaxMoves: 5,
		},
		outOfBusiness: false,
	}
}

func (g *Game) Reset() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.state = State{
		CarsAt1: r.Intn(g.Params.MaxCars + 1),
		CarsAt2: r.Intn(g.Params.MaxCars + 1),
	}
}

func (g Game) State() State { return g.state }
func (g Game) Over() bool   { return g.outOfBusiness }

func (g *Game) Play(i Input) State {
	// 0. Make sure the input is valid
	if i.Cars < -g.Params.MaxMoves {
		i.Cars = -g.Params.MaxMoves
	} else if i.Cars > g.Params.MaxMoves {
		i.Cars = g.Params.MaxMoves
	}

	// 1. Move the cars
	fmt.Println("At location 1:", g.state.CarsAt1, "At location 2:", g.state.CarsAt2)
	fmt.Println("Moving", i.Cars, "cars.")
	g.state.CarsAt1 = bounded(0, g.Params.MaxCars, g.state.CarsAt1-i.Cars)
	g.state.CarsAt2 = bounded(0, g.Params.MaxCars, g.state.CarsAt2+i.Cars)
	fmt.Println("At location 1:", g.state.CarsAt1, "At location 2:", g.state.CarsAt2)

	var c int
	// 2.a How many customers will Jack see today at the first location?
	c = g.Poisson.Draw(g.Params.CustomerAt1)
	fmt.Print(c, " customers want to rent today at location 1, ")
	if c > g.state.CarsAt1 {
		g.outOfBusiness = true
		g.state.CarsAt1 = 0
		fmt.Println("but Jack doesn't have enough cars...")
	} else {
		g.state.CarsAt1 -= c
		fmt.Println("and Jack is happy to show them to their cars")
	}

	// 2.b How many customers will Jack see today at the second location?
	c = g.Poisson.Draw(g.Params.CustomerAt2)
	fmt.Print(c, " customers want to rent today at location 2, ")
	if c > g.state.CarsAt2 {
		g.outOfBusiness = true
		g.state.CarsAt2 = 0
		fmt.Println("but Jack doesn't have enough cars...")
	} else {
		g.state.CarsAt2 -= c
		fmt.Println("and Jack is happy to show them to their cars")
	}

	// 3.a Get back some cars at location 1
	c = g.Poisson.Draw(g.Params.ReturnAt1)
	g.state.CarsAt1 = min(g.Params.MaxCars, g.state.CarsAt1+c)
	fmt.Print("Returns at first garage: ", c, ". ")

	// 3.b Get back some cars at location 2
	c = g.Poisson.Draw(g.Params.ReturnAt2)
	g.state.CarsAt2 = min(g.Params.MaxCars, g.state.CarsAt2+c)
	fmt.Print("Returns at second garage: ", c, ".")
	fmt.Println()

	fmt.Println()
	return g.state
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bounded(min, max, x int) int {
	if min > x {
		return min
	}
	if max < x {
		return max
	}
	return x
}

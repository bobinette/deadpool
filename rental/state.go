package rental

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type State struct {
	// CarsAtX represents the number of cars available for renting at location X
	CarsAt1 int
	CarsAt2 int
}

type Reward struct {
	P float64
	R float64
}

type Rewards map[State]Reward

type StateRewards map[State]Rewards

func (sr StateRewards) UnmarshalJSON(data []byte) error {
	internal := make([]struct {
		C1      int
		C2      int
		Rewards []struct {
			C1 int
			C2 int
			T  Reward
		}
	}, 0)

	if err := json.Unmarshal(data, &internal); err != nil {
		return err
	}

	for _, i := range internal {
		s := State{
			CarsAt1: i.C1,
			CarsAt2: i.C2,
		}

		for _, r := range i.Rewards {
			ns := State{
				CarsAt1: r.C1,
				CarsAt2: r.C2,
			}
			if _, ok := sr[s]; !ok {
				sr[s] = make(Rewards)
			}
			sr[s][ns] = r.T
		}
	}

	return nil
}

func (sr StateRewards) MarshalJSON() ([]byte, error) {
	internal := make([]struct {
		C1      int
		C2      int
		Rewards []struct {
			C1 int
			C2 int
			T  Reward
		}
	}, len(sr))

	i := 0
	for s, r := range sr {
		ir := internal[i]

		ir.C1 = s.CarsAt1
		ir.C2 = s.CarsAt2

		ir.Rewards = make([]struct {
			C1 int
			C2 int
			T  Reward
		}, len(r))
		j := 0
		for ns, t := range r {
			ir.Rewards[j].C1 = ns.CarsAt1
			ir.Rewards[j].C2 = ns.CarsAt2
			ir.Rewards[j].T = t
			j++
		}

		internal[i] = ir
		i++
	}

	return json.Marshal(internal)
}

type Policy map[State]int

func (p Policy) String() string {
	m1, m2 := 0, 0
	for s := range p {
		if s.CarsAt1 > m1 {
			m1 = s.CarsAt1
		}

		if s.CarsAt2 > m2 {
			m2 = s.CarsAt2
		}
	}

	buf := bytes.Buffer{}
	for i := m1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%2d | ", i))
		for j := 0; j <= m2; j++ {
			buf.WriteString(fmt.Sprintf(" %2d ", p[State{j, i}]))
		}
		buf.WriteRune('\n')
	}

	buf.WriteString("-----")
	for j := 0; j <= m2; j++ {
		buf.WriteString("----")
	}
	buf.WriteRune('\n')
	buf.WriteString("2/1| ")
	for j := 0; j <= m2; j++ {
		buf.WriteString(fmt.Sprintf(" %2d ", j))
	}
	return buf.String()
}

type StateValue map[State]float64

func (v StateValue) String() string {
	m1, m2 := 0, 0
	for s := range v {
		if s.CarsAt1 > m1 {
			m1 = s.CarsAt1
		}

		if s.CarsAt2 > m2 {
			m2 = s.CarsAt2
		}
	}

	buf := bytes.Buffer{}
	for i := m1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%2d | ", i))
		for j := 0; j <= m2; j++ {
			buf.WriteString(fmt.Sprintf(" %7.3f ", v[State{i, j}]))
		}
		buf.WriteRune('\n')
	}

	buf.WriteString("-----")
	for j := 0; j <= m2; j++ {
		buf.WriteString("---------")
	}
	buf.WriteRune('\n')
	buf.WriteString("2/1| ")
	for j := 0; j <= m2; j++ {
		buf.WriteString(fmt.Sprintf(" %7d ", j))
	}
	return buf.String()
}

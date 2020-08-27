package enigma

type Rotor struct {
	ID          string
	StraightSeq [26]int
	ReverseSeq  [26]int
	Turnover    []int

	Offset int
	Ring   int
}

func NewRotor(mapping string, id string, turnovers string) *Rotor {
	r := &Rotor{ID: id, Offset: 0, Ring: 0}
	r.Turnover = make([]int, len(turnovers))
	for i := range turnovers {
		r.Turnover[i] = CharToIndex(turnovers[i])
	}
	for i, letter := range mapping {
		index := CharToIndex(byte(letter))
		r.StraightSeq[i] = index
		r.ReverseSeq[index] = i
	}
	return r
}

func (r *Rotor) move(offset int) {
	r.Offset = (r.Offset + offset) % 26
}

func (r *Rotor) ShouldTurnOver() bool {
	for _, turnover := range r.Turnover {
		if r.Offset == turnover {
			return true
		}
	}
	return false
}

func (r *Rotor) Step(letter int, invert bool) int {
	letter = (letter - r.Ring + r.Offset + 26) % 26
	if invert {
		letter = r.ReverseSeq[letter]
	} else {
		letter = r.StraightSeq[letter]
	}
	letter = (letter + r.Ring - r.Offset + 26) % 26
	return letter
}

type Rotors []Rotor

func (rs *Rotors) GetByID(id string) *Rotor {
	for _, rotor := range *rs {
		if rotor.ID == id {
			return &rotor
		}
	}
	return nil
}

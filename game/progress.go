package game

type Progress struct {
	Level        int
	LinesCleared int
	Score        int
}

func (p *Progress) AddLinesCleared(rows int) {
	p.LinesCleared += rows
	p.calculateScore(rows)

	if p.Level == MAX_LEVEL {
		return
	}

	for i := p.LinesCleared; i > p.LinesCleared-rows; i-- {
		if i%10 == 0 {
			p.Level++
		}
	}

}

func (p *Progress) calculateScore(rows int) {
	multiplier := 1
	switch rows {
	case 1:
		multiplier = 1
	case 2:
		multiplier = 3
	case 3:
		multiplier = 5
	case 4:
		multiplier = 10
	default:
		multiplier = 10
	}

	p.Score += (BASE_SCORE * multiplier) * (p.Level + 1)
}

const MAX_LEVEL = 29
const BASE_SCORE = 100

package game

import "time"

type Progress struct {
	Level          int
	LinesCleared   int
	Score          int
	Ticker         *time.Ticker
	TickerDuration time.Duration
}

func (p *Progress) AddLinesCleared(rows int) {
	p.LinesCleared += rows
	p.calculateScore(rows)

	if p.Level == MAX_LEVEL {
		return
	}

	for i := p.LinesCleared; i > p.LinesCleared-rows; i-- {
		if i%10 == 0 {
			p.BumpLevel()
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

func (p *Progress) BumpLevel() {
	newLevel := p.Level + 1
	p.Level = newLevel
	p.Ticker.Stop()
	newDuration := calculateNextTicker(newLevel)
	p.Ticker = time.NewTicker(newDuration)
	p.TickerDuration = newDuration
}

func calculateNextTicker(level int) time.Duration {
	return (time.Duration(BASE_TICK) - (time.Duration(level * 25))) * time.Millisecond
}

const MAX_LEVEL = 29
const BASE_SCORE = 100
const BASE_TICK = 1000

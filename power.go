package oneroll

import "fmt"

// Power is a non-standard ability or miracle
type Power struct {
	Name       string
	Qualities  []*Quality
	Dice       *DiePool
	Effect     string
	Dud        bool
	Cost       int
	CostPerDie int
}

func (p Power) String() string {
	text := fmt.Sprintf("%s %s (",
		p.Name,
		p.Dice,
	)

	for _, q := range p.Qualities {
		text += fmt.Sprintf("%s", string(q.Type[0]))
		if q.Level > 1 {
			text += fmt.Sprintf("+%d", q.Level-1)
		}
	}

	text += fmt.Sprintf(") [%d/die] %dpts\n",
		p.CostPerDie,
		p.Cost)

	for _, q := range p.Qualities {
		text += fmt.Sprintln(q)
	}

	if p.Effect != "" {
		text += fmt.Sprintf("Effect: %s", p.Effect)
	}

	return text
}

// CalculatePowerCost totals the cost of Qualites for a Power
func (p *Power) CalculatePowerCost() {

	b := 0

	for _, q := range p.Qualities {

		// Update Quality DiePool if needed
		if q.Dice == nil {
			q.Dice = p.Dice
		}

		// Add Power Capacity Modifier if needed
		if len(q.Capacities) > 1 {
			tm := Modifiers["Power Capacity"]
			tm.Level = len(q.Capacities) - 1
			q.Modifiers = append(q.Modifiers, tm)
		}

		for _, m := range q.Modifiers {
			m.CalculateModifierCost(0)
		}
		q.CalculateQualityCost(2)
		b += q.CostPerDie
	}

	p.CostPerDie = b

	total := b * p.Dice.Normal
	total += b * 2 * p.Dice.Hard
	total += b * 4 * p.Dice.Wiggle

	p.Cost = total
}

// DeterminePowerCapacities calculates string values for powers
func (p *Power) DeterminePowerCapacities() {

	capacitiesMap := map[string]float64{
		"Mass":  25.0,
		"Range": 10.0,
		"Speed": 2.50,
		"Self":  0.0,
	}

	measuresMap := map[string]string{
		"Mass":  "kg",
		"Range": "m",
		"Speed": "m",
		"Self":  "",
	}

	var measure string

	totalDice := p.Dice.Normal + p.Dice.Hard + p.Dice.Wiggle

	for _, q := range p.Qualities {
		for _, c := range q.Capacities {
			baseVal := capacitiesMap[c.Type]

			modVal := baseVal

			// Double value for each die above 1
			for i := 1; i < totalDice; i++ {
				modVal = modVal * 2.0
			}

			boosterVal := 1.0

			// Apply booster
			for _, m := range q.Modifiers {
				if m.Name == "Booster" {
					boosterVal = float64(m.Level) * 10.0
				}
			}
			// Get final value
			finalVal := float64(modVal) * boosterVal

			if finalVal > 1000.0 {
				switch {
				case c.Type == "Range":
					finalVal = finalVal / 1000.0
					measure = "km"
					c.Value = fmt.Sprintf("%.2f%s", finalVal, measure)
				case c.Type == "Mass":
					finalVal = finalVal / 1000.0
					measure = "tons"
					c.Value = fmt.Sprintf("%.2f%s", finalVal, measure)
				case c.Type == "Speed":
					finalVal = finalVal / 1000.0
					measure = "km"
					c.Value = fmt.Sprintf("%.2f%s", finalVal, measure)
				case c.Type == "Self":
					c.Value = "Self"
				}
			} else {
				measure = measuresMap[c.Type]
				c.Value = fmt.Sprintf("%.0f%s", finalVal, measure)
			}
		} // End Capacities
	} // End Qualities
}

// NewPower generates a new empty Power
func NewPower(t string) *Power {

	p := new(Power)

	p.Name = t
	p.Effect = ""
	p.Qualities = []*Quality{}
	p.Dice = &DiePool{}
	p.Dud = false

	// Take user input

	return p
}
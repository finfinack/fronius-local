package mqtt

import (
	"fmt"
	"time"

	"github.com/finfinack/fronius-local/pkg/fronius"
)

type Message struct {
	UpdatedAt time.Time

	Battery       *fronius.Battery
	CurrentLimits *CurrentLimits
}

type CurrentLimits struct {
	ChargeMax    int `json:"charge_max"`
	ChargeMin    int `json:"charge_min"`
	DischargeMax int `json:"discharge_max"`
	DischargeMin int `json:"discharge_min"`
}

func NewMessage(batt *fronius.Battery, timeOfUse []*fronius.TimeOfUse) (*Message, error) {
	msg := &Message{
		UpdatedAt:     time.Now(),
		Battery:       batt,
		CurrentLimits: &CurrentLimits{},
	}

	for _, t := range timeOfUse {
		act, err := t.IsActive()
		if err != nil {
			return nil, fmt.Errorf("unable to check if active: %s", err)
		}
		if !act {
			continue
		}

		switch t.ScheduleType {
		case fronius.ScheduleTypeChargeMax:
			msg.CurrentLimits.ChargeMax = t.Power
		case fronius.ScheduleTypeChargeMin:
			msg.CurrentLimits.ChargeMin = t.Power
		case fronius.ScheduleTypeDischargeMax:
			msg.CurrentLimits.DischargeMax = t.Power
		case fronius.ScheduleTypeDischargeMin:
			msg.CurrentLimits.DischargeMin = t.Power
		default:
			return nil, fmt.Errorf("unknown schedule type: %s", t.ScheduleType)
		}
	}

	return msg, nil
}

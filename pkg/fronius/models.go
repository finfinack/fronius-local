package fronius

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	froniusTimeFmt = "15:04"

	ScheduleTypeChargeMax    = "CHARGE_MAX"
	ScheduleTypeChargeMin    = "CHARGE_MIN"
	ScheduleTypeDischargeMax = "DISCHARGE_MAX"
	ScheduleTypeDischargeMin = "DISCHARGE_MIN"
)

type Battery struct {
	BatteryCalibration                bool   `json:"BAT_CALIBRATION"`
	BatteryEnabled                    bool   `json:"BAT_ENABLED"`
	BatteryM0SOCMax                   int    `json:"BAT_M0_SOC_MAX"`
	BatteryM0SOCMin                   int    `json:"BAT_M0_SOC_MIN"`
	BatteryM0SOCMode                  string `json:"BAT_M0_SOC_MODE"`
	BatteryServiceOn                  bool   `json:"BAT_SERVICE_ON"`
	BattyryType                       string `json:"BAT_TYPE"`
	HYBBackupCriticalSOC              int    `json:"HYB_BACKUP_CRITICALSOC"`
	HYBBackupReserved                 int    `json:"HYB_BACKUP_RESERVED"`
	HYBBackupSystemDeadlockPrevention bool   `json:"HYB_BACKUP_SYSTEMDEADLOCKPREVENTION"`
	HYBMBChargeFromAC                 bool   `json:"HYB_BM_CHARGEFROMAC"`
	HYBBMPacMin                       int    `json:"HYB_BM_PACMIN"`
	HYBEMMode                         int    `json:"HYB_EM_MODE"`
	HYBEMPower                        int    `json:"HYB_EM_POWER"`
	HYBEVUAccoupled                   bool   `json:"HYB_EVU_ACCOUPLED"`
	HYBEVUChargeFromGrid              bool   `json:"HYB_EVU_CHARGEFROMGRID"`
	PVPeakPower                       int    `json:"PV_PEAK_POWER"`
}

type TimeOfUseResponse struct {
	TimeOfUse []*TimeOfUse `json:"TimeOfUse"`
}

type TimeOfUse struct {
	Active       bool       `json:"Active"`
	Power        int        `json:"Power"`
	ScheduleType string     `json:"ScheduleType"`
	TimeTable    *TimeTable `json:"TimeTable"`
	Weekdays     *Weekdays  `json:"Weekdays"`
}

func (t *TimeOfUse) IsActive() (bool, error) {
	if !t.Active {
		return false, nil
	}
	if t.Weekdays == nil || t.TimeTable == nil {
		return false, errors.New("Weekdays and TimeTable needs to be set")
	}
	if !t.Weekdays.IsActive() {
		return false, nil
	}
	if !t.TimeTable.IsActive() {
		return false, nil
	}
	return true, nil
}

type Weekdays struct {
	Monday    bool `json:"Mon"`
	Tuesday   bool `json:"Tue"`
	Wednesday bool `json:"Wed"`
	Thursday  bool `json:"Thu"`
	Friday    bool `json:"Fri"`
	Saturday  bool `json:"Sat"`
	Sunday    bool `json:"Sun"`
}

func (w *Weekdays) IsActive() bool {
	switch time.Now().Weekday() {
	case time.Monday:
		return w.Monday
	case time.Tuesday:
		return w.Tuesday
	case time.Wednesday:
		return w.Wednesday
	case time.Thursday:
		return w.Thursday
	case time.Friday:
		return w.Friday
	case time.Saturday:
		return w.Saturday
	case time.Sunday:
		return w.Sunday
	default:
		return false
	}
}

type TimeTable struct {
	Start TimeTableTime `json:"Start"`
	End   TimeTableTime `json:"End"`
}

func (t *TimeTable) IsActive() bool {
	now := time.Now()
	return now.After(time.Time(t.Start)) && now.Before(time.Time(t.End))
}

type TimeTableTime time.Time

func (t *TimeTableTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	ti, err := time.Parse(froniusTimeFmt, s)
	if err != nil {
		return err
	}
	ti = ti.AddDate(time.Now().Year(), int(time.Now().Month()-1), time.Now().Day()-1)
	*t = TimeTableTime(ti)
	return nil
}

func (t TimeTableTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(froniusTimeFmt))
}

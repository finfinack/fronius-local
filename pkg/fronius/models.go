package fronius

import "time"

type Battery struct {
	Time time.Time `json:"time"`

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

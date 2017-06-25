package instance

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	cfg_folder     = "cfg"
	server_ini     = "server_cfg.ini"
	entry_list_ini = "entry_list.ini"
	sep            = "\n"
)

func writeConfig(config *model.Configuration) error {
	iniPath := filepath.Join(cfg_folder, int64ToStr(config.Id))
	if err := os.MkdirAll(iniPath, 0755); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error creating cfg folder")
		return err
	}

	if err := writeIniFile(config, ServerConfigToIniString(config), server_ini); err != nil {
		return err
	}

	if err := writeIniFile(config, EntryListToIniString(config), entry_list_ini); err != nil {
		return err
	}

	return nil
}

func writeIniFile(config *model.Configuration, ini, filename string) error {
	iniFile := filepath.Join(cfg_folder, int64ToStr(config.Id), filename)
	if err := ioutil.WriteFile(iniFile, []byte(ini), 0775); err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("Error writing %s", filename)
		return err
	}

	return nil
}

func ServerConfigToIniString(config *model.Configuration) string {
	ini := "[SERVER]" + sep
	ini += "NAME=" + config.Name + sep
	ini += "CARS=" + getCars(config) + sep
	ini += "CONFIG_TRACK=" + config.TrackConfig + sep
	ini += "TRACK=" + config.Track + sep
	ini += "SUN_ANGLE=" + intToStr(config.SunAngle) + sep
	ini += "PASSWORD=" + config.Pwd + sep
	ini += "ADMIN_PASSWORD=" + config.AdminPwd + sep
	ini += "UDP_PORT=" + intToStr(config.UDP) + sep
	ini += "TCP_PORT=" + intToStr(config.TCP) + sep
	ini += "HTTP_PORT=" + intToStr(config.HTTP) + sep
	ini += "MAX_BALLAST_KG=" + intToStr(config.MaxBallast) + sep
	ini += "QUALIFY_MAX_WAIT_PERC=120" + sep
	ini += "RACE_PIT_WINDOW_START=" + intToStr(config.RacePitWindowStart) + sep
	ini += "RACE_PIT_WINDOW_END=" + intToStr(config.RacePitWindowEnd) + sep
	ini += "REVERSED_GRID_RACE_POSITIONS=" + intToStr(config.ReversedGridRacePos) + sep
	ini += "LOCKED_ENTRY_LIST=" + boolToStr(config.LockEntryList) + sep
	ini += "PICKUP_MODE_ENABLED=" + boolToStr(config.PickupMode) + sep
	ini += "LOOP_MODE=" + boolToStr(config.LoopMode) + sep
	ini += "SLEEP_TIME=1" + sep
	ini += "CLIENT_SEND_INTERVAL_HZ=" + intToStr(config.PacketsHz) + sep
	ini += "SEND_BUFFER_SIZE=0" + sep
	ini += "RECV_BUFFER_SIZE=0" + sep
	ini += "RACE_OVER_TIME=" + intToStr(config.RaceOvertime) + sep
	ini += "KICK_QUORUM=" + intToStr(config.KickVoteQuorum) + sep
	ini += "VOTING_QUORUM=" + intToStr(config.SessionVoteQuorum) + sep
	ini += "VOTE_DURATION=" + intToStr(config.VoteDuration) + sep
	ini += "BLACKLIST_MODE=" + intToStr(config.Blacklist) + sep
	ini += "FUEL_RATE=" + intToStr(config.FuelRate) + sep
	ini += "DAMAGE_MULTIPLIER=" + intToStr(config.DamageRate) + sep
	ini += "TYRE_WEAR_RATE=" + intToStr(config.TiresWearRate) + sep
	ini += "ALLOWED_TYRES_OUT=" + intToStr(config.AllowedTiresOut) + sep
	ini += "ABS_ALLOWED=" + intToStr(config.ABS) + sep
	ini += "TC_ALLOWED=" + intToStr(config.TC) + sep
	ini += "START_RULE=1" + sep
	ini += "RACE_GAS_PENALTY_DISABLED=" + boolToStr(config.DisableGasCutPenality) + sep
	ini += "RESULT_SCREEN_TIME=" + intToStr(config.ResultScreenTime) + sep
	ini += "MAX_CONTACTS_PER_KM=" + intToStr(config.MaxCollisionsKm) + sep
	ini += "STABILITY_ALLOWED=" + boolToStr(config.StabilityAid) + sep
	ini += "AUTOCLUTCH_ALLOWED=" + boolToStr(config.AutoClutch) + sep
	ini += "TYRE_BLANKETS_ALLOWED=" + boolToStr(config.TyreBlankets) + sep
	ini += "FORCE_VIRTUAL_MIRROR=" + boolToStr(config.ForceVirtualMirror) + sep
	ini += "REGISTER_TO_LOBBY=" + boolToStr(config.ShowInLobby) + sep
	ini += "MAX_CLIENTS=" + intToStr(config.MaxSlots) + sep
	ini += "NUM_THREADS=" + intToStr(config.Threads) + sep
	ini += "UDP_PLUGIN_LOCAL_PORT=" + intToStr(config.UdpPluginPort) + sep
	ini += "UDP_PLUGIN_ADDRESS=" + config.UdpPluginAddr + sep
	ini += "AUTH_PLUGIN_ADDRESS=" + sep
	ini += "LEGAL_TYRES=" + config.LegalTyres + sep
	ini += "RACE_EXTRA_LAP=" + boolToStr(config.RaceExtraLap) + sep
	ini += "WELCOME_MESSAGE=" + config.Welcome + sep

	if config.Practice {
		ini += sep
		ini += "[PRACTICE]" + sep
		ini += "NAME=Practice" + sep
		ini += "TIME=" + intToStr(config.PracticeTime) + sep
		ini += "IS_OPEN=" + boolToStr(config.CanJoinPractice) + sep
	}

	if config.Qualify {
		ini += sep
		ini += "[QUALIFY]" + sep
		ini += "NAME=Qualify" + sep
		ini += "TIME=" + intToStr(config.QualifyTime) + sep
		ini += "IS_OPEN=" + boolToStr(config.CanJoinQualify) + sep
	}

	if config.Race {
		ini += sep
		ini += "[RACE]" + sep
		ini += "NAME=Race" + sep
		ini += "LAPS=" + intToStr(config.RaceLaps) + sep
		ini += "TIME=" + intToStr(config.RaceTime) + sep
		ini += "WAIT_TIME=" + intToStr(config.RaceWaitTime) + sep
		ini += "IS_OPEN=" + intToStr(config.JoinType) + sep
	}

	if config.DynamicTrack {
		ini += sep
		ini += "[DYNAMIC_TRACK]" + sep
		ini += "SESSION_START=" + intToStr(config.StartValue) + sep
		ini += "RANDOMNESS=" + intToStr(config.Randomness) + sep
		ini += "SESSION_TRANSFER=" + intToStr(config.TransferredGrip) + sep
		ini += "LAP_GAIN=" + intToStr(config.LapsToImproveGrip) + sep
	}

	// weather
	for i, w := range config.Weather {
		ini += sep
		ini += "[WEATHER_" + intToStr(i) + "]" + sep
		ini += "GRAPHICS=" + w.Weather + sep
		ini += "BASE_TEMPERATURE_AMBIENT=" + intToStr(w.BaseAmbientTemp) + sep
		ini += "BASE_TEMPERATURE_ROAD=" + intToStr(w.BaseRoadTemp) + sep
		ini += "VARIATION_AMBIENT=" + intToStr(w.AmbientVariation) + sep
		ini += "VARIATION_ROAD=" + intToStr(w.RoadVariation) + sep
		ini += "WIND_BASE_SPEED_MIN=" + intToStr(w.WindBaseSpeedMin) + sep
		ini += "WIND_BASE_SPEED_MAX=" + intToStr(w.WindBaseSpeedMax) + sep
		ini += "WIND_BASE_DIRECTION=" + intToStr(w.WindBaseDirection) + sep
		ini += "WIND_VARIATION_DIRECTION=" + intToStr(w.WindVariationDirection) + sep
	}

	ini += sep
	ini += "[DATA]" + sep
	ini += "DESCRIPTION=" + sep
	ini += "EXSERVEREXE=" + sep
	ini += "EXSERVERBAT=" + sep
	ini += "EXSERVERHIDEWIN=0" + sep
	ini += "WEBLINK=" + sep
	ini += "WELCOME_PATH=" + sep

	return ini
}

func EntryListToIniString(config *model.Configuration) string {
	ini := ""

	for i, car := range config.Cars {
		ini += "[CAR_" + intToStr(i) + "]" + sep
		ini += "MODEL=" + car.Car + sep
		ini += "SKIN=" + car.Painting + sep
		ini += "SPECTATOR_MODE=" + boolToStr(car.Spectator) + sep
		ini += "DRIVERNAME=" + car.Driver + sep
		ini += "TEAM=" + car.Team + sep
		ini += "GUID=" + car.GUID + sep
		ini += "BALLAST=0" + sep
		ini += "FIXED_SETUP=" + car.FixedSetup + sep
		ini += sep
	}

	return ini
}

func getCars(config *model.Configuration) string {
	cars := make([]string, 0)

	for _, car := range config.Cars {
		found := false

		for _, str := range cars {
			if str == car.Car {
				found = true
				break
			}
		}

		if !found {
			cars = append(cars, car.Car)
		}
	}

	return strings.Join(cars, ";")
}

func boolToStr(b bool) string {
	if b {
		return "1"
	}

	return "0"
}

func intToStr(i int) string {
	return strconv.Itoa(i)
}

func int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

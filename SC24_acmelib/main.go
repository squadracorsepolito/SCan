package main

import (
	"log"
	"os"

	"github.com/squadracorsepolito/acmelib"
)

var nodeIDs = map[string]acmelib.NodeID{
	"TLB_BAT":    1,
	"SB_FRONT":   2,
	"SB_REAR":    3,
	"BMS_LV":     4,
	"DASH":       5,
	"DIAG_TOOL":  6,
	"DSPACE":     7,
	"EXTRA_NODE": 8,
	"SCANNER":    9,
	"TPMS":       10,
	"IMU":        11,
	"BRUSA":      12,
}

var messageIDs = map[string]acmelib.MessageID{
	"DSPACE__timeAndDate":          1,
	"TLB_BAT__signalsStatus":       4,
	"SB_FRONT__analogDevice":       5,
	"SB_REAR__analogDevice":        6,
	"SB_REAR__criticalPeripherals": 7,

	"BMS_LV__lvBatGeneral":    20,
	"DASH__hmiDevicesState":   22,
	"DSPACE__peripheralsCTRL": 25,

	"DIAG_TOOL__xcpTxTLB_BAT":    40,
	"DIAG_TOOL__xcpTxSB_FRONT":   41,
	"DIAG_TOOL__xcpTxSB_REAR":    42,
	"DIAG_TOOL__xcpTxBMS_LV":     43,
	"DIAG_TOOL__xcpTxDASH":       44,
	"DIAG_TOOL__xcpTxSCANNER":    45,
	"TLB_BAT__SDCsensingStatus":  46,
	"SB_REAR__SDCsensingStatus":  47,
	"SB_FRONT__SDCsensingStatus": 48,
	"SB_FRONT__potentiometer":    50,
	"SB_REAR__potentiometer":     51,

	"TLB_BAT__xcpTx":            70,
	"SB_FRONT__xcpTx":           70,
	"SB_REAR__xcpTx":            70,
	"BMS_LV__xcpTx":             70,
	"DASH__xcpTx":               70,
	"DSPACE__xcpTx":             70,
	"SCANNER__xcpTx":            70,
	"DSPACE__signals":           73,
	"DSPACE__fsmStates":         74,
	"BMS_LV__cellsStatus":       75,
	"BMS_LV__status":            76,
	"BMS_LV__lvCellVoltage0":    77,
	"BMS_LV__lvCellVoltage1":    78,
	"DASH__peripheralsStatus":   79,
	"TPMS__frontWheelsPressure": 80,
	"TPMS__rearWheelsPressure":  81,

	"TLB_BAT__hello":               100,
	"SB_FRONT__hello":              100,
	"SB_REAR__hello":               100,
	"BMS_LV__hello":                100,
	"DASH__hello":                  100,
	"DSPACE__hello":                100,
	"BMS_LV__lvCellNTCResistance0": 101,
	"BMS_LV__lvCellNTCResistance1": 102,
	"SB_FRONT__ntcResistance":      103,
	"SB_REAR__ntcResistance":       104,
	"DASH__appsRangeLimits":        105,
	"DASH__carCommands":            106,
	"DSPACE__dashLedsColorRGB":     107,
}

const originalDBC = "./../SC24/artifacts/MCB/MCB.dbc"
const genDBC = "generated/MCB.dbc"
const genJSON = "generated/SC24.json"
const genWire = "generated/SC24.binpb"
const genMD = "generated/SC24.md"

func main() {
	sc24 := acmelib.NewNetwork("SC24")
	sc24.SetDesc("The CAN network of the squadracorse 2024 formula SAE car")

	// load mcb
	mcbFile, err := os.Open(originalDBC)
	checkErr(err)
	defer mcbFile.Close()

	mcb, err := acmelib.ImportDBCFile("mcb", mcbFile)
	checkErr(err)
	checkErr(mcb.UpdateName("Main CAN Bus"))
	checkErr(sc24.AddBus(mcb))

	// renaming signal types
	dashInt, err := mcb.GetNodeInterfaceByNodeName("DASH")
	checkErr(err)

	tmpSigType, err := acmelib.NewIntegerSignalType("fan_pwm_t", 4, false)
	checkErr(err)
	tmpSigType.SetMax(10)
	modifySignalType(dashInt, "DASH__peripheralsStatus", "TSAC_FAN_pwmDutyCycleStatus", tmpSigType)

	modifySignalTypeName(dashInt, "DASH__hmiDevicesState", "ROT_SW_1_state", "rotary_switch_state_t")
	modifySignalTypeName(dashInt, "DASH__appsRangeLimits", "APPS_0_voltageRangeMin", "uint16_t")
	modifySignalTypeName(dashInt, "DASH__carCommands", "BMS_LV_diagPWD", "bms_lv_password_t")

	bmslvInt, err := mcb.GetNodeInterfaceByNodeName("BMS_LV")
	checkErr(err)

	modifySignalTypeName(bmslvInt, "BMS_LV__hello", "FW_majorVersion", "uint8_t")
	modifySignalTypeName(bmslvInt, "BMS_LV__lvCellVoltage0", "LV_CELL_0_voltage", "lv_cell_voltage_t")
	modifySignalTypeName(bmslvInt, "BMS_LV__lvCellNTCResistance0", "LV_CELL_NTC_00_resistance", "ntc_resistance_t")
	modifySignalTypeName(bmslvInt, "BMS_LV__lvBatGeneral", "LV_BAT_voltage", "lv_bat_voltage_t")
	modifySignalTypeName(bmslvInt, "BMS_LV__lvBatGeneral", "LV_BAT_currentSensVoltage", "lv_bat_current_sens_t")

	dspaceInt, err := mcb.GetNodeInterfaceByNodeName("DSPACE")
	checkErr(err)

	tmpSigType, err = acmelib.NewIntegerSignalType("seconds_t", 6, false)
	checkErr(err)
	tmpSigType.SetMax(59)
	modifySignalType(dspaceInt, "DSPACE__timeAndDate", "DATETIME_seconds", tmpSigType)

	modifySignalTypeName(dspaceInt, "DSPACE__timeAndDate", "DATETIME_month", "month_t")
	modifySignalTypeName(dspaceInt, "DSPACE__timeAndDate", "DATETIME_day", "day_t")
	modifySignalTypeName(dspaceInt, "DSPACE__timeAndDate", "DATETIME_hours", "hours_t")
	modifySignalTypeName(dspaceInt, "DSPACE__timeAndDate", "DATETIME_minutes", "minutes_t")

	// calculte bus load
	mcb.SetBaudrate(1_000_000)
	busLoad, err := acmelib.CalculateBusLoad(mcb, 1000)
	checkErr(err)
	log.Print("BUS LOAD: ", busLoad)

	// parse IDs
	parseNodeIDs(mcb)
	parseMessageIDs(mcb)

	// save files
	dbcFile, err := os.Create(genDBC)
	checkErr(err)
	defer dbcFile.Close()
	acmelib.ExportBus(dbcFile, mcb)

	wireFile, err := os.Create(genWire)
	checkErr(err)
	defer wireFile.Close()
	jsonFile, err := os.Create(genJSON)
	checkErr(err)
	defer jsonFile.Close()
	checkErr(acmelib.SaveNetwork(sc24, acmelib.SaveEncodingWire|acmelib.SaveEncodingJSON, wireFile, jsonFile, nil))

	mdFile, err := os.Create(genMD)
	checkErr(err)
	defer mdFile.Close()

	checkErr(acmelib.ExportToMarkdown(sc24, mdFile))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func modifySignalTypeName(nodeInt *acmelib.NodeInterface, msgName, sigName, newName string) {
	tmpMsg, err := nodeInt.GetMessageByName(msgName)
	checkErr(err)
	tmpSig, err := tmpMsg.GetSignalByName(sigName)
	checkErr(err)
	tmpStdSig, err := tmpSig.ToStandard()
	checkErr(err)
	tmpStdSig.Type().SetName(newName)
}

func modifySignalType(nodeInt *acmelib.NodeInterface, msgName, sigName string, newType *acmelib.SignalType) {
	tmpMsg, err := nodeInt.GetMessageByName(msgName)
	checkErr(err)
	tmpSig, err := tmpMsg.GetSignalByName(sigName)
	checkErr(err)
	tmpStdSig, err := tmpSig.ToStandard()
	checkErr(err)
	tmpStdSig.SetType(newType)
}

func parseNodeIDs(mcb *acmelib.Bus) {
	interfaces := mcb.NodeInterfaces()
	for idx, tmpInt := range interfaces {
		checkErr(tmpInt.Node().UpdateID(acmelib.NodeID(100 + idx)))
	}

	for i := len(interfaces) - 1; i >= 0; i = i - 1 {
		tmpNodeInt := interfaces[i]
		tmpNode := tmpNodeInt.Node()
		if nodeID, ok := nodeIDs[tmpNode.Name()]; ok {
			checkErr(tmpNode.UpdateID(nodeID))
		}
	}
}

func parseMessageIDs(mcb *acmelib.Bus) {
	for _, tmpNodeInt := range mcb.NodeInterfaces() {
		for _, tmpMsg := range tmpNodeInt.Messages() {
			msgID, ok := messageIDs[tmpMsg.Name()]
			if ok {
				checkErr(tmpMsg.UpdateID(msgID))
			} else {
				log.Print("Message not found: ", tmpMsg.Name())
			}
		}
	}
}

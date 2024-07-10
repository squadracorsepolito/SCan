# INFOS
## ECU-IDS
The list of ECUs and their ids:
- TLB_BAT: 1
- SB_FRONT: 2
- SB_REAR: 3
- BMS_LV: 4
- DASH: 5
- DIAG_TOOL: 6
- DSPACE: 7
- EXTRA_NODE: 8  `DON'T CHANGE THIS ECU-ID: KEEP IT AT 8 !`  
Some ECUs have non modifiable can-ids because programmed at fabbrication/delivery, we put those here
- SCANNER: 9
- TPMS: 10
- IMU: 11
## Criteria For Message-Id assignment
Below values are in base 10
### Blacklisted Message-ids (using ecu-id encoding)
- EXTRA-NODES, ECU-ID 8 (0x8)
    - Message-id: 112 (0x70): IRTS sensors

### High Priority high speed (1-50ms) 
```
- Message IDs: 0-19 [20]
    - 1: DSPACE-date_time
    - 4: TLB_BAT_signalsStatus
    - 5: SB_FRONT_analogDevices
    - 6: SB_REAR_analogDevices
    - 7: SB_REAR_criticalPeripherals
```
### High Priority low speed (50-100+ms)
```
- Message IDs: 20-39 [20] 
    - 20: BMS_LV-batGeneral
    - 22: DASH_hmiDevicesState
    - 25: DSPACE_peripheralsCTRL
    - 26: SB_REAR_dischargeStatus
```
### Mid Priority high speed (1-50ms)
```
- Message IDs: 40-69 [30] 
    - 40-45: DIAG_TOOL_xcpTx[all ecus]
    - 46: TLB_BAT_sdcSensingStatus
    - 47: SB_REAR_sdcSensingStatus
    - 48: SB_FRONT_sdcSensingStatus
    - 50: SB_FRONT_potentiometer
    - 51: SB_REAR_potentiometer
```
### Mid Priority Low Speed (50-100+ms)
```
- Message IDs: 70-99 [30] 
    - 70: [all-ecus]_xcpTx
    - 73: DSPACE_signals
    - 74: DSPACE_fsmStates
    - 75: BMS_LV-cellsStatus
    - 76: BMS_LV-status
    - 77-78: BMS_LV-cellVolt0/1
    - 79: Dash_periphStatus
    - 80: TPMS_frontWheelsPressure
    - 81: TPMS_rearWheelsPressure
```
### Low Priority
```
- Message IDs: 100-127 [27] (exclude 112)
    - 100: [all-ecus]_hello
    - 101-102: BMS_LV_ntcResistance0/1
    - 103: SB_FRONT-ntcResistance
    - 104: SB_REAR-ntcResistance
    - 105: DASH-appsRangeLimits
    - 106: DASH-carCommands
    - 107-108 : DSPACE-ledColors
    - 112: BALCKLISTED see above
```

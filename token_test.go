package geoqlparser

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenizer_ScanNMEA(t *testing.T) {
	testCases := []struct {
		name string
		set  map[string]Token
	}{
		{
			name: "NMEA_AAM_*",
			set: map[string]Token{
				"nmea:aam:statusArrivalCircleEntered": NMEA_AAM_STATUS_ARRIVAL_CIRCLE_ENTERED,
				"nmea:aam:statusPerpendicularPassed":  NMEA_AAM_STATUS_PERPENDICULAR_PASSED,
				"nmea:aam:arrivalCircleRadius":        NMEA_AAM_ARRIVAL_CIRCLE_RADIUS,
				"nmea:aam:arrivalCircleRadiusUnit":    NMEA_AAM_ARRIVAL_CIRCLE_UNIT,
				"nmea:aam:destinationWaypointID":      NMEA_AAM_DESTINATION_WAYPOINT_ID,
			},
		},
		{
			name: "NMEA_ALA_*",
			set: map[string]Token{
				"nmea:ala:time":               NMEA_ALA_TIME,
				"nmea:ala:systemIndicator":    NMEA_ALA_SYSTEM_INDICATOR,
				"nmea:ala:subSystemIndicator": NMEA_ALA_SUB_SYSTEM_INDICATOR,
				"nmea:ala:instanceNumber":     NMEA_ALA_INSTANCE_NUMBER,
				"nmea:ala:type":               NMEA_ALA_TYPE,
				"nmea:ala:condition":          NMEA_ALA_CONDITION,
				"nmea:ala:alarmAckState":      NMEA_ALA_ALARM_ACK_STATE,
				"nmea:ala:message":            NMEA_ALA_MESSAGE,
			},
		},
		{
			name: "NMEA_APB_*",
			set: map[string]Token{
				"nmea:apb:statusGeneralWarning":       NMEA_APB_STATUS_GENERAL_WARNING,
				"nmea:apb:statusLockWarning":          NMEA_APB_STATUS_LOCK_WARNING,
				"nmea:apb:crossTrackErrorMagnitude":   NMEA_APB_CROSS_TRACK_ERROR_MAGNITUDE,
				"nmea:apb:directionToSteer":           NMEA_APB_DIRECTION_TO_STEER,
				"nmea:apb:crossTrackUnits":            NMEA_APB_CROSS_TRACK_UNITS,
				"nmea:apb:statusArrivalCircleEntered": NMEA_APB_STATUS_ARRIVAL_CIRCLE_ENTERED,
				"nmea:apb:statusPerpendicularPassed":  NMEA_APB_STATUS_PERPENDICULAR_PASSED,
				"nmea:apb:bearingOriginToDest":        NMEA_APB_BEARING_ORIGIN_TO_DEST,
				"nmea:apb:bearingOriginToDestType":    NMEA_APB_BEARING_ORIGIN_DEST_TYPE,
				"nmea:apb:destinationWaypointID":      NMEA_APB_DESTINATION_WAYPOINT_ID,
				"nmea:apb:bearingPresentToDest":       NMEA_APB_BEARING_PRESENT_TO_DEST,
				"nmea:apb:bearingPresentToDestType":   NMEA_APB_BEARING_PRESENT_TO_DEST_TYPE,
				"nmea:apb:heading":                    NMEA_APB_HEADING,
				"nmea:apb:headingType":                NMEA_APB_HEADING_TYPE,
				"nmea:apb:ffaMode":                    NMEA_APB_FFA_MODE,
			},
		},
		{
			name: "NMEA_BEC_*",
			set: map[string]Token{
				"nmea:bec:time":                       NMEA_BEC_TIME,
				"nmea:bec:latitude":                   NMEA_BEC_LATITUDE,
				"nmea:bec:longitude":                  NMEA_BEC_LONGITUDE,
				"nmea:bec:bearingTrue":                NMEA_BEC_BEARING_TRUE,
				"nmea:bec:bearingTrueValid":           NMEA_BEC_BEARING_TRUE_VALID,
				"nmea:bec:bearingMagnetic":            NMEA_BEC_BEARING_MAGNETIC,
				"nmea:bec:bearingMagneticValid":       NMEA_BEC_BEARING_MAGNETIC_VALID,
				"nmea:bec:distanceNauticalMiles":      NMEA_BEC_DISTANCE_NAUTICAL_MILES,
				"nmea:bec:distanceNauticalMilesValid": NMEA_BEC_DISTANCE_NAUTICAL_MILES_VALID,
				"nmea:bec:destinationWaypointID":      NMEA_BEC_DESTINATION_WAYPOINT_ID,
			},
		},
		{
			name: "NMEA_BOD_*",
			set: map[string]Token{
				"nmea:bod:bearingTrue":           NMEA_BOD_BEARING_TRUE,
				"nmea:bod:bearingTrueType":       NMEA_BOD_BEARING_TRUE_TYPE,
				"nmea:bod:bearingMagnetic":       NMEA_BOD_BEARING_MAGNETIC,
				"nmea:bod:bearingMagneticType":   NMEA_BOD_BEARING_MAGNETIC_TYPE,
				"nmea:bod:destinationWaypointID": NMEA_BOD_DESTINATION_WAYPOINT_ID,
				"nmea:bod:originWaypointID":      NMEA_BOD_ORIGIN_WAYPOINT_ID,
			},
		},
		{
			name: "NMEA_BWC_*",
			set: map[string]Token{
				"nmea:bwc:time":                      NMEA_BWC_TIME,
				"nmea:bwc:latitude":                  NMEA_BWC_LATITUDE,
				"nmea:bwc:longitude":                 NMEA_BWC_LONGITUDE,
				"nmea:bwc:bearingTrue":               NMEA_BWC_BEARING_TRUE,
				"nmea:bwc:bearingTrueType":           NMEA_BWC_BEARING_TRUE_TYPE,
				"nmea:bwc:bearingMagnetic":           NMEA_BWC_BEARING_MAGNETIC,
				"nmea:bwc:bearingMagneticType":       NMEA_BWC_BEARING_MAGNETIC_TYPE,
				"nmea:bwc:distanceNauticalMiles":     NMEA_BWC_DISTANCE_NAUTICAL_MILES,
				"nmea:bwc:distanceNauticalMilesUnit": NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT,
			},
		},
		{
			name: "NMEA_BOD_*",
			set: map[string]Token{
				"nmea:bod:bearingTrue":           NMEA_BOD_BEARING_TRUE,
				"nmea:bod:bearingTrueType":       NMEA_BOD_BEARING_TRUE_TYPE,
				"nmea:bod:bearingMagnetic":       NMEA_BOD_BEARING_MAGNETIC,
				"nmea:bod:bearingMagneticType":   NMEA_BOD_BEARING_MAGNETIC_TYPE,
				"nmea:bod:destinationWaypointID": NMEA_BOD_DESTINATION_WAYPOINT_ID,
				"nmea:bod:originWaypointID":      NMEA_BOD_ORIGIN_WAYPOINT_ID,
			},
		},
		{
			name: "NMEA_BWC_*",
			set: map[string]Token{
				"nmea:bwc:time":                      NMEA_BWC_TIME,
				"nmea:bwc:latitude":                  NMEA_BWC_LATITUDE,
				"nmea:bwc:longitude":                 NMEA_BWC_LONGITUDE,
				"nmea:bwc:bearingTrue":               NMEA_BWC_BEARING_TRUE,
				"nmea:bwc:bearingTrueType":           NMEA_BWC_BEARING_TRUE_TYPE,
				"nmea:bwc:bearingMagnetic":           NMEA_BWC_BEARING_MAGNETIC,
				"nmea:bwc:bearingMagneticType":       NMEA_BWC_BEARING_MAGNETIC_TYPE,
				"nmea:bwc:distanceNauticalMiles":     NMEA_BWC_DISTANCE_NAUTICAL_MILES,
				"nmea:bwc:distanceNauticalMilesUnit": NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT,
			},
		},
		{
			name: "NMEA_BWC_*",
			set: map[string]Token{
				"nmea:bwr:time":                  NMEA_BWR_TIME,
				"nmea:bwr:latitude":              NMEA_BWR_LATITUDE,
				"nmea:bwr:longitude":             NMEA_BWR_LONGITUDE,
				"nmea:bwr:bearingTrue":           NMEA_BWR_BEARING_TRUE,
				"nmea:bwr:bearingTrueType":       NMEA_BWR_BEARING_TRUE_TYPE,
				"nmea:bwr:bearingMagnetic":       NMEA_BWR_BEARING_MAGNETIC,
				"nmea:bwr:bearingMagneticType":   NMEA_BWR_BEARING_MAGNETIC_TYPE,
				"nmea:bwr:distanceNauticalMiles": NMEA_BWR_DISTANCE_NAUTICAL_MILES,
				"nmea:bwr:destinationWaypointID": NMEA_BWR_DESTINATION_WAYPOINT_ID,
				"nmea:bwr:ffaMode":               NMEA_BWR_FFA_MODE,
			},
		},
		{
			name: "NMEA_BWW_*",
			set: map[string]Token{
				"nmea:bww:bearingTrue":           NMEA_BWW_BEARING_TRUE,
				"nmea:bww:bearingTrueType":       NMEA_BWW_BEARING_TRUE_TYPE,
				"nmea:bww:bearingMagnetic":       NMEA_BWW_BEARING_MAGNETIC,
				"nmea:bww:bearingMagneticType":   NMEA_BWW_BEARING_MAGNETIC_TYPE,
				"nmea:bww:destinationWaypointID": NMEA_BWW_DESTINATION_WAYPOINT_ID,
				"nmea:bww:originWaypointID":      NMEA_BWW_ORIGINAL_WAYPOINT_ID,
			},
		},
		{
			name: "NMEA_DBK_*",
			set: map[string]Token{
				"nmea:dbk:depthFeet":        NMEA_DBK_DEPTH_FEET,
				"nmea:dbk:depthFeetUnit":    NMEA_DBK_DEPTH_FEET_UNIT,
				"nmea:dbk:depthMeters":      NMEA_DBK_DEPTH_METERS,
				"nmea:dbk:depthMetersUnit":  NMEA_DBK_DEPTH_METERS_UNIT,
				"nmea:dbk:depthFathoms":     NMEA_DBK_DEPTH_FATHOMS,
				"nmea:dbk:depthFathomsUnit": NMEA_DBK_DEPTH_FATHOMS_UNIT,
			},
		},
		{
			name: "NMEA_DBS_*",
			set: map[string]Token{
				"nmea:dbs:depthFeet":       NMEA_DBS_DEPTH_FEET,
				"nmea:dbs:depthFeetUnit":   NMEA_DBS_DEPTH_FEET_UNIT,
				"nmea:dbs:depthMeters":     NMEA_DBS_DEPTH_METERS,
				"nmea:dbs:depthMeterUnit":  NMEA_DBS_DEPTH_METERS_UNIT,
				"nmea:dbs:depthFathoms":    NMEA_DBS_DEPTH_FATHOMS,
				"nmea:dbs:depthFathomUnit": NMEA_DBS_DEPTH_FATHOMS_UNIT,
			},
		},
		{
			name: "NMEA_DBT_*",
			set: map[string]Token{
				"nmea:dbt:depthFeet":    NMEA_DBT_DEPTH_FEET,
				"nmea:dbt:depthMeters":  NMEA_DBT_DEPTH_METERS,
				"nmea:dbt:depthFathoms": NMEA_DBT_DEPTH_FATHOMS,
			},
		},
		{
			name: "NMEA_DOR_*",
			set: map[string]Token{
				"nmea:dor:type":               NMEA_DOR_TYPE,
				"nmea:dor:time":               NMEA_DOR_TIME,
				"nmea:dor:systemIndicator":    NMEA_DOR_SYSTEM_INDICATOR,
				"nmea:dor:divisionIndicator1": NMEA_DOR_DIVISION_INDICATOR1,
				"nmea:dor:divisionIndicator2": NMEA_DOR_DIVISION_INDICATOR2,
				"nmea:dor:doorNumberOrCount":  NMEA_DOR_DOOR_NUMBER_OR_COUNT,
				"nmea:dor:doorStatus":         NMEA_DOR_DOOR_STATUS,
				"nmea:dor:switchSetting":      NMEA_DOR_SWITCH_SETTING,
				"nmea:dor:message":            NMEA_DOR_MESSAGE,
			},
		},
		{
			name: "NMEA_DPT_*",
			set: map[string]Token{
				"nmea:dpt:depth":      NMEA_DPT_DEPTH,
				"nmea:dpt:offset":     NMEA_DPT_OFFSET,
				"nmea:dpt:rangeScale": NMEA_DPT_RANGE_SCALE,
			},
		},
		{
			name: "NMEA_DSC_*",
			set: map[string]Token{
				"nmea:dsc:formatSpecifier":             NMEA_DSC_FORMAT_SPECIFIER,
				"nmea:dsc:address":                     NMEA_DSC_ADDRESS,
				"nmea:dsc:category":                    NMEA_DSC_CATEGORY,
				"nmea:dsc:distressCauseOrTeleCommand1": NMEA_DSC_DISTRESS_CAUSE_OR_TELE_CMD1,
				"nmea:dsc:commandTypeOrTeleCommand2":   NMEA_DSC_CMD_TYPE_OR_TELE_CMD2,
				"nmea:dsc:positionOrCanal":             NMEA_DSC_POSITION_OR_CANAL,
				"nmea:dsc:timeOrTelephoneNumber":       NMEA_DSC_TIMER_OR_TELE_NUMBER,
				"nmea:dsc:mmsi":                        NMEA_DSC_MMSI,
				"nmea:dsc:distressCause":               NMEA_DSC_DISTREES_CAUSE,
				"nmea:dsc:acknowledgement":             NMEA_DSC_ACK_LEDGEMENT,
				"nmea:dsc:expansionIndicator":          NMEA_DSC_EXPANSION_INDICATOR,
			},
		},
		{
			name: "NMEA_DSE_*",
			set: map[string]Token{
				"nmea:dse:totalNumber":     NMEA_DSE_TOTAL_NUMBER,
				"nmea:dse:number":          NMEA_DSE_NUMBER,
				"nmea:dse:acknowledgement": NMEA_DSE_ACK_LEDGEMENT,
				"nmea:dse:mmsi":            NMEA_DSE_MMSI,
			},
		},
		{
			name: "NMEA_DTM_*",
			set: map[string]Token{
				"nmea:dtm:localDatumCode":        NMEA_DTM_LOCAL_DATUM_CODE,
				"nmea:dtm:localDatumSubcode":     NMEA_DTM_LOCAL_DATUM_SUBCODE,
				"nmea:dtm:latitudeOffsetMinute":  NMEA_DTM_LATITUDE_OFFSET_MINUTE,
				"nmea:dtm:longitudeOffsetMinute": NMEA_DTM_LONGITUDE_OFFSET_MINUTE,
				"nmea:dtm:altitudeOffsetMeters":  NMEA_DTM_ALTITUDE_OFFSET_METERS,
				"nmea:dtm:datumName":             NMEA_DTM_DATUM_NAME,
			},
		},
		{
			name: "NMEA_EVE_*",
			set: map[string]Token{
				"nmea:eve:time":    NMEA_EVE_TIME,
				"nmea:eve:tagCode": NMEA_EVE_TAG_CODE,
				"nmea:eve:message": NMEA_EVE_MESSAGE,
			},
		},
		{
			name: "NMEA_FIR_*",
			set: map[string]Token{
				"nmea:fir:type":                      NMEA_FIR_TYPE,
				"nmea:fir:time":                      NMEA_FIR_TIME,
				"nmea:fir:systemIndicator":           NMEA_FIR_SYSTEM_INDICATOR,
				"nmea:fir:divisionIndicator1":        NMEA_FIR_DIVISION_INDICATOR1,
				"nmea:fir:divisionIndicator2":        NMEA_FIR_DIVISION_INDICATOR2,
				"nmea:fir:fireDetectorNumberOrCount": NMEA_FIR_FIRE_DETECTOR_NUM_OR_COUNT,
				"nmea:fir:condition":                 NMEA_FIR_CONDITION,
				"nmea:fir:alarmAckState":             NMEA_FIR_ALARAM_ACK_STATE,
				"nmea:fir:message":                   NMEA_FIR_MESSAGE,
			},
		},
		{
			name: "NMEA_GGA_*",
			set: map[string]Token{
				"nmea:gga:time":          NMEA_GGA_TIME,
				"nmea:gga:latitude":      NMEA_GGA_LATITUDE,
				"nmea:gga:longitude":     NMEA_GGA_LONGITUDE,
				"nmea:gga:fixQuality":    NMEA_GGA_FIX_QUOLITY,
				"nmea:gga:numSatellites": NMEA_GGA_NUM_SATELLITES,
				"nmea:gga:hdop":          NMEA_GGA_HDOP,
				"nmea:gga:altitude":      NMEA_GGA_ALTITUDE,
				"nmea:gga:separation":    NMEA_GGA_SEPARATION,
				"nmea:gga:dgspAge":       NMEA_GGA_DGPS_AGE,
				"nmea:gga:dgspId":        NMEA_GGA_DGSP_ID,
			},
		},
		{
			name: "NMEA_GGL_*",
			set: map[string]Token{
				"nmea:gll:latitude":  NMEA_GLL_LATITUDE,
				"nmea:gll:longitude": NMEA_GLL_LONGITUDE,
				"nmea:gll:time":      NMEA_GLL_TIME,
				"nmea:gll:validity":  NMEA_GLL_VALIDITY,
				"nmea:gll:ffaMode":   NMEA_GLL_FFAMODE,
			},
		},
		{
			name: "NMEA_GNS_*",
			set: map[string]Token{
				"nmea:gns:time":       NMEA_GNS_TIME,
				"nmea:gns:latitude":   NMEA_GNS_LATITUDE,
				"nmea:gns:longitude":  NMEA_GNS_LONGITUDE,
				"nmea:gns:altitude":   NMEA_GNS_ALTITUDE,
				"nmea:gns:mode":       NMEA_GNS_MODE,
				"nmea:gns:svs":        NMEA_GNS_SVS,
				"nmea:gns:hdop":       NMEA_GNS_HDOP,
				"nmea:gns:separation": NMEA_GNS_SEPARATION,
				"nmea:gns:age":        NMEA_GNS_AGE,
				"nmea:gns:station":    NMEA_GNS_STATION,
				"nmea:gns:navStatus":  NMEA_GNS_NAV_STATUS,
			},
		},
		{
			name: "NMEA_GSA_*",
			set: map[string]Token{
				"nmea:gsa:mode":     NMEA_GSA_MODE,
				"nmea:gsa:fixType":  NMEA_GSA_FIX_TYPE,
				"nmea:gsa:sv":       NMEA_GSA_SV,
				"nmea:gsa:pdop":     NMEA_GSA_PDOP,
				"nmea:gsa:hdop":     NMEA_GSA_HDOP,
				"nmea:gsa:vdop":     NMEA_GSA_VDOP,
				"nmea:gsa:systemID": NMEA_GSA_SYSTEM_ID,
			},
		},
		{
			name: "NMEA_HDG_*",
			set: map[string]Token{
				"nmea:hdg:heading":            NMEA_HDG_HEADING,
				"nmea:hdg:deviation":          NMEA_HDG_DEVIATION,
				"nmea:hdg:deviationDirection": NMEA_HDG_DEVIATION_DIRECTION,
				"nmea:hdg:variation":          NMEA_HDG_VARIATION,
				"nmea:hdg:variationDirection": NMEA_HDG_VARIATION_DIRECTION,
			},
		},
		{
			name: "NMEA_HDM_*",
			set: map[string]Token{
				"nmea:hdm:heading":       NMEA_HDM_HEADING,
				"nmea:hdm:magneticValid": NMEA_HDM_MAGNETIC_VALID,
			},
		},
		{
			name: "NMEA_HDT_*",
			set: map[string]Token{
				"nmea:hdt:heading": NMEA_HDT_HEADING,
				"nmea:hdt:true":    NMEA_HDT_TRUE,
			},
		},
		{
			name: "NMEA_HSC_*",
			set: map[string]Token{
				"nmea:hsc:trueHeading":         NMEA_HSC_TRUE_HEADING,
				"nmea:hsc:trueHeadingType":     NMEA_HSC_TRUE_HEADING_TYPE,
				"nmea:hsc:magneticHeading":     NMEA_HSC_MAGNETIC_HEADING,
				"nmea:hsc:magneticHeadingType": NMEA_HSC_MAGNETIC_HEADING_TYPE,
			},
		},
		{
			name: "NMEA_MDA_*",
			set: map[string]Token{
				"nmea:mda:pressureInch":          NMEA_MDA_PRESSURE_INCH,
				"nmea:mda:inchesValid":           NMEA_MDA_INCHES_VALID,
				"nmea:mda:pressureBar":           NMEA_MDA_PRESSURE_BAR,
				"nmea:mda:barsValid":             NMEA_MDA_BARS_VALID,
				"nmea:mda:airTemp":               NMEA_MDA_AIR_TEMP,
				"nmea:mda:airTempValid":          NMEA_MDA_AIR_TEMP_VALID,
				"nmea:mda:waterTemp":             NMEA_MDA_WATER_TEMP,
				"nmea:mda:waterTempValid":        NMEA_MDA_WATER_TEMP_VALID,
				"nmea:mda:relativeHum":           NMEA_MDA_RELATIVE_HUM,
				"nmea:mda:absoluteHum":           NMEA_MDA_ABSOLUTE_HUM,
				"nmea:mda:dewPoint":              NMEA_MDA_DEW_POINT,
				"nmea:mda:dewPointValid":         NMEA_MDA_DEW_POINT_VALID,
				"nmea:mda:windDirectionTrue":     NMEA_MDA_WIND_DIRECTION_TRUE,
				"nmea:mda:trueValid":             NMEA_MDA_TRUE_VALID,
				"nmea:mda:windDirectionMagnetic": NMEA_MDA_WIND_DIRECTION_MAGNETIC,
				"nmea:mda:magneticValid":         NMEA_MDA_MAGNETIC_VALID,
				"nmea:mda:windSpeedKnots":        NMEA_MDA_WIND_SPEED_KNOTS,
				"nmea:mda:knotsValid":            NMEA_MDA_KNOTS_VALID,
				"nmea:mda:windSpeedMeters":       NMEA_MDA_WIND_SPEED_METERS,
				"nmea:mda:metersValid":           NMEA_MDA_METERS_VALID,
			},
		},
		{
			name: "NMEA_MTA_*",
			set: map[string]Token{
				"nmea:mta:temperature": NMEA_MTA_TEMPERATURE,
				"nmea:mta:unit":        NMEA_MTA_UNIT,
			},
		},
		{
			name: "NMEA_MTK_*",
			set: map[string]Token{
				"nmea:mtk:command": NMEA_MTK_COMMAND,
				"nmea:mtk:flag":    NMEA_MTK_FLAG,
			},
		},
		{
			name: "NMEA_MTW_*",
			set: map[string]Token{
				"nmea:mtw:temperature":  NMEA_MTW_TEMPERATURE,
				"nmea:mtw:celsiusValid": NMEA_MTW_CELSIUS_VALID,
			},
		},
		{
			name: "NMEA_MWD_*",
			set: map[string]Token{
				"nmea:mwd:windDirectionTrue":     NMEA_MWD_WIND_DIRECTION_TRUE,
				"nmea:mwd:trueValid":             NMEA_MWD_TRUE_VALID,
				"nmea:mwd:windDirectionMagnetic": NMEA_MWD_WIND_DIRECTION_MAGNETIC,
				"nmea:mwd:magneticValid":         NMEA_MWD_MAGNETIC_VALID,
				"nmea:mwd:windSpeedKnots":        NMEA_MWD_WIND_SPEED_KNOTS,
				"nmea:mwd:knotsValid":            NMEA_MWD_KNOTS_VALID,
				"nmea:mwd:windSpeedMeters":       NMEA_MWD_WIND_SPEED_METERS,
				"nmea:mwd:metersValid":           NMEA_MWD_METERS_VALID,
			},
		},
		{
			name: "NMEA_MWV_*",
			set: map[string]Token{
				"nmea:mwv:windAngle":     NMEA_MWV_WIND_ANGLE,
				"nmea:mwv:reference":     NMEA_MWV_REFERENCE,
				"nmea:mwv:windSpeed":     NMEA_MWV_WIND_SPEED,
				"nmea:mwv:windSpeedUnit": NMEA_MWV_WIND_SPEED_UNIT,
				"nmea:mwv:statusValid":   NMEA_MWV_STATUS_VALID,
			},
		},
		{
			name: "NMEA_OSD_*",
			set: map[string]Token{
				"nmea:osd:heading":          NMEA_OSD_HEADING,
				"nmea:osd:headingStatus":    NMEA_OSD_HEADING_STATUS,
				"nmea:osd:vesselTrueCourse": NMEA_OSD_VESSEL_TRUE_COURSE,
				"nmea:osd:courseReference":  NMEA_OSD_COURSE_REFERENCE,
				"nmea:osd:vesselSpeed":      NMEA_OSD_VESSEL_SPEED,
				"nmea:osd:speedReference":   NMEA_OSD_SPEED_REFERENCE,
				"nmea:osd:vesselSetTrue":    NMEA_OSD_VESSEL_SET_TRUE,
				"nmea:osd:vesselDrift":      NMEA_OSD_VESSEL_DRIFT,
				"nmea:osd:speedUnits":       NMEA_OSD_SPEED_UNITS,
			},
		},
		{
			name: "NMEA_PGRME_*",
			set: map[string]Token{
				"nmea:pgrme:horizontal": NMEA_PGRME_HORIZONTAL,
				"nmea:pgrme:vertical":   NMEA_PGRME_VERTICAL,
				"nmea:pgrme:spherical":  NMEA_PGRME_SPHERICAL,
			},
		},
		{
			name: "NMEA_PHTRO_*",
			set: map[string]Token{
				"nmea:phtro:pitch": NMEA_PHTRO_PITCH,
				"nmea:phtro:bow":   NMEA_PHTRO_BOW,
				"nmea:phtro:roll":  NMEA_PHTRO_ROLL,
				"nmea:phtro:port":  NMEA_PHTRO_PORT,
			},
		},
		{
			name: "NMEA_PRDID_*",
			set: map[string]Token{
				"nmea:prdid:pitch":   NMEA_PRDID_PITCH,
				"nmea:prdid:roll":    NMEA_PRDID_ROLL,
				"nmea:prdid:heading": NMEA_PRDID_HEADING,
			},
		},
		{
			name: "NMEA_PSKPDPT_*",
			set: map[string]Token{
				"nmea:pskpdpt:depth":              NMEA_PSKPDPT_DEPTH,
				"nmea:pskpdpt:offset":             NMEA_PSKPDPT_OFFSET,
				"nmea:pskpdpt:rangeScale":         NMEA_PSKPDPT_RANGE_SCALE,
				"nmea:pskpdpt:bottomEchoStrength": NMEA_PSKPDPT_BOTTOM_ECHO_STRENGTH,
				"nmea:pskpdpt:channelNumber":      NMEA_PSKPDPT_CHANNEL_NUMBER,
				"nmea:pskpdpt:transducerLocation": NMEA_PSKPDPT_TRANSDUCER_LOCATION,
			},
		},
		{
			name: "NMEA_PSONCMS_*",
			set: map[string]Token{
				"nmea:psoncms:quaternion0":       NMEA_PSONCMS_QUATERNION0,
				"nmea:psoncms:quaternion1":       NMEA_PSONCMS_QUATERNION1,
				"nmea:psoncms:quaternion2":       NMEA_PSONCMS_QUATERNION2,
				"nmea:psoncms:quaternion3":       NMEA_PSONCMS_QUATERNION3,
				"nmea:psoncms:accelerationX":     NMEA_PSONCMS_ACCELERATION_X,
				"nmea:psoncms:accelerationY":     NMEA_PSONCMS_ACCELERATION_Y,
				"nmea:psoncms:accelerationZ":     NMEA_PSONCMS_ACCELERATION_Z,
				"nmea:psoncms:rateOfTurnX":       NMEA_PSONCMS_RATE_OF_TURN_X,
				"nmea:psoncms:rateOfTurnY":       NMEA_PSONCMS_RATE_OF_TURN_Z,
				"nmea:psoncms:magneticFieldX":    NMEA_PSONCMS_MAGNETIC_FIELD_X,
				"nmea:psoncms:magneticFieldY":    NMEA_PSONCMS_MAGNETIC_FIELD_Y,
				"nmea:psoncms:rateOfTurnZ":       NMEA_PSONCMS_MAGNETIC_FIELD_Z,
				"nmea:psoncms:sensorTemperature": NMEA_PSONCMS_SENSOR_TEMPERATURE,
			},
		},
		{
			name: "NMEA_RMB_*",
			set: map[string]Token{
				"nmea:rmb:dataStatus":                      NMEA_RMB_DATA_STATUS,
				"nmea:rmb:crossTrackErrorNauticalMiles":    NMEA_RMB_CROSS_TRACK_ERROR_NAUTICAL_MILES,
				"nmea:rmb:directionToSteer":                NMEA_RMB_DIRECTION_TO_STEER,
				"nmea:rmb:originWaypointID":                NMEA_RMB_ORIGIN_WAYPOINT_ID,
				"nmea:rmb:destinationWaypointID":           NMEA_RMB_DESTINATION_WAYPOINT_ID,
				"nmea:rmb:destinationLatitude":             NMEA_RMB_DESTINATION_LATITUDE,
				"nmea:rmb:destinationLongitude":            NMEA_RMB_DESTINATION_LONGITUDE,
				"nmea:rmb:rangeToDestinationNauticalMiles": NMEA_RMB_RANGE_TO_DESTINATION_NAUTICAL_MILES,
				"nmea:rmb:trueBearingToDestination":        NMEA_RMB_TRUE_BEARING_TO_DESTINATION,
				"nmea:rmb:velocityToDestinationKnots":      NMEA_RMB_VELOCITY_TO_DESTINATION_KNOTS,
				"nmea:rmb:arrivalStatus":                   NMEA_RMB_ARRIVAL_STATUS,
				"nmea:rmb:ffaMode":                         NMEA_RMB_FFAMODE,
			},
		},
		{
			name: "NMEA_RMC_*",
			set: map[string]Token{
				"nmea:rmc:time":      NMEA_RMC_TIME,
				"nmea:rmc:validity":  NMEA_RMC_VALIDITY,
				"nmea:rmc:latitude":  NMEA_RMC_LATITUDE,
				"nmea:rmc:longitude": NMEA_RMC_LONGITUDE,
				"nmea:rmc:speed":     NMEA_RMC_SPEED,
				"nmea:rmc:course":    NMEA_RMC_COURSE,
				"nmea:rmc:date":      NMEA_RMC_DATE,
				"nmea:rmc:variation": NMEA_RMC_VARIATION,
				"nmea:rmc:ffaMode":   NMEA_RMC_FFAMODE,
				"nmea:rmc:navStatus": NMEA_RMC_NAV_STATUS,
			},
		},
		{
			name: "NMEA_ROT_*",
			set: map[string]Token{
				"nmea:rot:rateOfTurn": NMEA_ROT_RATE_OF_TURN,
				"nmea:rot:valid":      NMEA_ROT_VALID,
			},
		},
		{
			name: "NMEA_RPM_*",
			set: map[string]Token{
				"nmea:rpm:source":       NMEA_RPM_SOURCE,
				"nmea:rpm:engineNumber": NMEA_RPM_ENGINE_NUMBER,
				"nmea:rpm:speedRPM":     NMEA_RPM_SPEED_RPM,
				"nmea:rpm:pitchPercent": NMEA_RPM_PITCH_PERCENT,
				"nmea:rpm:status":       NMEA_RPM_STATUS,
			},
		},
		{
			name: "NMEA_RSA_*",
			set: map[string]Token{
				"nmea:rsa:starboardRudderAngle":       NMEA_RSA_STARBOARD_RUDDER_ANGLE,
				"nmea:rsa:starboardRudderAngleStatus": NMEA_RSA_STARBOARD_RUDDER_ANGLE_STATUS,
				"nmea:rsa:portRudderAngle":            NMEA_RSA_PORT_RUDDER_ANGLE,
				"nmea:rsa:portRudderAngleStatus":      NMEA_RSA_PORT_RUDDER_ANGLE_STATUS,
			},
		},
		{
			name: "NMEA_RSD_*",
			set: map[string]Token{
				"nmea:rsd:origin1Range":           NMEA_RSD_ORIGIN1_RANGE,
				"nmea:rsd:origin1Bearing":         NMEA_RSD_ORIGIN1_BEARING,
				"nmea:rsd:variableRangeMarker1":   NMEA_RSD_VARIABLE_RANGE_MARKET1,
				"nmea:rsd:bearingLine1":           NMEA_RSD_BEARING_LINE1,
				"nmea:rsd:origin2Range":           NMEA_RSD_ORIGIN2_RANGE,
				"nmea:rsd:origin2Bearing":         NMEA_RSD_ORIGIN2_BEARING,
				"nmea:rsd:variableRangeMarker2":   NMEA_RSD_VARIABLE_RANGE_MARKET2,
				"nmea:rsd:bearingLine2":           NMEA_RSD_BEARING_LINE2,
				"nmea:rsd:cursorRangeFromOwnShip": NMEA_RSD_CURSOR_RANGE_FROM_OWN_SHIP,
				"nmea:rsd:cursorBearingDegrees":   NMEA_RSD_CURSOR_BEARING_DEGREES,
				"nmea:rsd:rangeScale":             NMEA_RSD_RANGE_SCALE,
				"nmea:rsd:rangeUnit":              NMEA_RSD_RANGE_UNIT,
				"nmea:rsd:displayRotation":        NMEA_RSD_DISPLAY_ROTATION,
			},
		},
		{
			name: "NMEA_RTE_*",
			set: map[string]Token{
				"nmea:rte:numberOfSentences":         NMEA_RTE_NUMBER_OF_SENTENCES,
				"nmea:rte:sentenceNumber":            NMEA_RTE_SENTENCE_NUMBER,
				"nmea:rte:activeRouteOrWaypointList": NMEA_RTE_ACTIVE_ROUTER_OR_WAYPOINT_LIST,
				"nmea:rte:name":                      NMEA_RTE_NAME,
				"nmea:rte:idents":                    NMEA_RTE_IDENTS,
			},
		},
		{
			name: "NMEA_THS_*",
			set: map[string]Token{
				"nmea:ths:heading": NMEA_THS_HEADING,
				"nmea:ths:status":  NMEA_THS_STATUS,
			},
		},
		{
			name: "NMEA_TLL_*",
			set: map[string]Token{
				"nmea:tll:targetNumber":    NMEA_TLL_TARGET_NUMBER,
				"nmea:tll:targetLatitude":  NMEA_TLL_TARGET_LATITUDE,
				"nmea:tll:targetLongitude": NMEA_TLL_TARGET_LONGITUDE,
				"nmea:tll:targetName":      NMEA_TLL_TARGET_NAME,
				"nmea:tll:timeUTC":         NMEA_TLL_TIME_UTC,
				"nmea:tll:targetStatus":    NMEA_TLL_TARGET_STATUS,
				"nmea:tll:referenceTarget": NMEA_TLL_REFERENCE_TARGET,
			},
		},
		{
			name: "NMEA_TTM_*",
			set: map[string]Token{
				"nmea:ttm:targetNumber":      NMEA_TTM_TARGET_NUMBER,
				"nmea:ttm:targetDistance":    NMEA_TTM_TARGET_DISTANCE,
				"nmea:ttm:bearing":           NMEA_TTM_BEARING,
				"nmea:ttm:bearingType":       NMEA_TTM_BEARING_TYPE,
				"nmea:ttm:targetSpeed":       NMEA_TTM_TARGET_SPEED,
				"nmea:ttm:targetCourse":      NMEA_TTM_TARGET_COURSE,
				"nmea:ttm:courseType":        NMEA_TTM_COURSE_TYPE,
				"nmea:ttm:distanceCPA":       NMEA_TTM_DISTANCE_CPA,
				"nmea:ttm:timeCPA":           NMEA_TTM_TIME_CPA,
				"nmea:ttm:speedUnits":        NMEA_TTM_SPEED_UNITS,
				"nmea:ttm:targetName":        NMEA_TTM_TARGET_NAME,
				"nmea:ttm:targetStatus":      NMEA_TTM_TARGET_STATUS,
				"nmea:ttm:referenceTarget":   NMEA_TTM_REFERENCE_TARGET,
				"nmea:ttm:timeUTC":           NMEA_TTM_TIME_UTC,
				"nmea:ttm:typeOfAcquisition": NMEA_TTM_TYPE_OF_ACQUISITION,
			},
		},
		{
			name: "NMEA_VBW_*",
			set: map[string]Token{
				"nmea:vbw:longitudinalWaterSpeedKnots":         NMEA_VBW_LONGITUDINAL_WATER_SPEED_KNOTS,
				"nmea:vbw:transverseWaterSpeedKnots":           NMEA_VBW_TRANSVERSE_WATER_SPEED_KNOTS,
				"nmea:vbw:waterSpeedStatusValid":               NMEA_VBW_WATER_SPEED_STATUS_VALID,
				"nmea:vbw:waterSpeedStatus":                    NMEA_VBW_WATER_SPEED_STATUS,
				"nmea:vbw:longitudinalGroundSpeedKnots":        NMEA_VBW_LONGITUDINAL_GROUND_SPEED_KNOTS,
				"nmea:vbw:transverseGroundSpeedKnots":          NMEA_VBW_TRANSVERSE_GROUND_SPEED_KNOTS,
				"nmea:vbw:groundSpeedStatusValid":              NMEA_VBW_GROUND_SPEED_STATUS_VALID,
				"nmea:vbw:groundSpeedStatus":                   NMEA_VBW_GROUND_SPEED_STATUS,
				"nmea:vbw:sternTraverseWaterSpeedKnots":        NMEA_VBW_STERN_TRAVERSE_WATER_SPEED_KNOTS,
				"nmea:vbw:sternTraverseWaterSpeedStatusValid":  NMEA_VBW_STERN_TRAVERSE_WATER_SPEED_STATUS_VALID,
				"nmea:vbw:sternTraverseWaterSpeedStatus":       NMEA_VBW_STERN_TRAVERSE_WATER_SPEED_STATUS,
				"nmea:vbw:sternTraverseGroundSpeedKnots":       NMEA_VBW_STERN_TRAVERSE_GROUND_SPEED_KNOTS,
				"nmea:vbw:sternTraverseGroundSpeedStatusValid": NMEA_VBW_STERN_TRAVERSE_GROUND_SPEED_STATUS_VALID,
				"nmea:vbw:sternTraverseGroundSpeedStatus":      NMEA_VBW_STERN_TRAVERSE_GROUND_SPEED_STATUS,
			},
		},
		{
			name: "NMEA_VDR_*",
			set: map[string]Token{
				"nmea:vdr:setDegreesTrue":         NMEA_VDR_SET_DEGREES_TRUE,
				"nmea:vdr:setDegreesTrueUnit":     NMEA_VDR_SET_DEGREES_TRUE_UNIT,
				"nmea:vdr:setDegreesMagnetic":     NMEA_VDR_SET_DEGREES_MAGNETIC,
				"nmea:vdr:setDegreesMagneticUnit": NMEA_VDR_SET_DEGREES_MAGNETIC_UNIT,
				"nmea:vdr:driftKnots":             NMEA_VDR_DRIFT_KNOTS,
				"nmea:vdr:driftUnit":              NMEA_VDR_DRIFT_UNIT,
			},
		},
		{
			name: "NMEA_VHW_*",
			set: map[string]Token{
				"nmea:vhw:trueHeading":            NMEA_VHW_TRUE_HEADING,
				"nmea:vhw:magneticHeading":        NMEA_VHW_MAGNETIC_HEADING,
				"nmea:vhw:speedThroughWaterKnots": NMEA_VHW_SPEED_THROUGHT_WATER_KNOTS,
				"nmea:vhw:speedThroughWaterKPH":   NMEA_VHW_SPEED_THROUGHT_WATER_KPH,
			},
		},
		{
			name: "NMEA_VLW_*",
			set: map[string]Token{
				"nmea:vlw:totalInWater":           NMEA_VLW_TOTAL_IN_WATER,
				"nmea:vlw:totalInWaterUnit":       NMEA_VLW_TOTAL_IN_WATER_UNIT,
				"nmea:vlw:sinceResetInWater":      NMEA_VLW_SINCE_RESET_IN_WATER,
				"nmea:vlw:sinceResetInWaterUnit":  NMEA_VLW_SINCE_RESET_IN_WATER_UNIT,
				"nmea:vlw:totalOnGround":          NMEA_VLW_TOTAL_ON_GROUND,
				"nmea:vlw:totalOnGroundUnit":      NMEA_VLW_TOTAL_ON_GROUND_UNIT,
				"nmea:vlw:sinceResetOnGround":     NMEA_VLW_SINCE_RESET_ON_GROUND,
				"nmea:vlw:sinceResetOnGroundUnit": NMEA_VLW_SINCE_RESET_ON_GROUND_UNIT,
			},
		},
		{
			name: "NMEA_VPW_*",
			set: map[string]Token{
				"nmea:vpw:speedKnots":     NMEA_VPW_SPEED_KNOTS,
				"nmea:vpw:speedKnotsUnit": NMEA_VPW_SPEED_KNOTS_UNIT,
				"nmea:vpw:speedMPS":       NMEA_VPW_SPEED_MPS,
				"nmea:vpw:speedMPSUnit":   NMEA_VPW_SPEED_MPS_UNIT,
			},
		},
		{
			name: "NMEA_VTG_*",
			set: map[string]Token{
				"nmea:vtg:trueTrack":        NMEA_VTG_TRUE_TRACK,
				"nmea:vtg:magneticTrack":    NMEA_VTG_MAGNETIC_TRACK,
				"nmea:vtg:groundSpeedKnots": NMEA_VTG_GROUND_SPEED_KNOTS,
				"nmea:vtg:groundSpeedKPH":   NMEA_VTG_GROUND_SPEED_KPH,
				"nmea:vtg:ffaMode":          NMEA_VTG_FFAMODE,
			},
		},
		{
			name: "NMEA_VWR_*",
			set: map[string]Token{
				"nmea:vwr:measuredAngle":        NMEA_VWR_MEASURED_ANGLE,
				"nmea:vwr:measuredDirectionBow": NMEA_VWR_MEASURED_DIRECTION_BOW,
				"nmea:vwr:speedKnots":           NMEA_VWR_SPEED_KNOTS,
				"nmea:vwr:speedKnotsUnit":       NMEA_VWR_SPEED_KNOTS_UNIT,
				"nmea:vwr:speedMPS":             NMEA_VWR_SPEED_MPS,
				"nmea:vwr:speedMPSUnit":         NMEA_VWR_SPEED_MPS_UNIT,
				"nmea:vwr:speedKPH":             NMEA_VWR_SPEED_KPH,
				"nmea:vwr:speedKPHUnit":         NMEA_VWR_SPEED_KPH_UNIT,
			},
		},
		{
			name: "NMEA_VWT_*",
			set: map[string]Token{
				"nmea:vwt:trueAngle":        NMEA_VWT_TRUE_ANGLE,
				"nmea:vwt:trueDirectionBow": NMEA_VWT_TRUE_DIRECTION_BOW,
				"nmea:vwt:speedKnots":       NMEA_VWT_SPEED_KNOTS,
				"nmea:vwt:speedKnotsUnit":   NMEA_VWT_SPEED_KNOTS_UNIT,
				"nmea:vwt:speedMPS":         NMEA_VWT_SPEED_MPS,
				"nmea:vwt:speedMPSUnit":     NMEA_VWT_SPEED_MPS_UNIT,
				"nmea:vwt:speedKPH":         NMEA_VWT_SPEED_KPH,
				"nmea:vwt:speedKPHUnit":     NMEA_VWT_SPEED_KPH_UNIT,
			},
		},
		{
			name: "NMEA_WPL_*",
			set: map[string]Token{
				"nmea:wpl:latitude":  NMEA_WPL_LATITUDE,
				"nmea:wpl:longitude": NMEA_WPL_LONGITUDE,
				"nmea:wpl:ident":     NMEA_WPL_IDENT,
			},
		},
	}

	toStr := func(set map[string]Token) (str string) {
		for k := range set {
			str += k + " "
		}
		return
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str := toStr(tc.set)
			tokenizer := NewTokenizer(strings.NewReader(str))
			for i := 0; i < len(tc.set); i++ {
				haveTok, haveLit := tokenizer.Scan()
				wantTok, ok := tc.set[haveLit]
				if !ok {
					t.Fatalf("token not found %s", haveLit)
				}
				if haveTok != wantTok {
					t.Fatalf("have %d, want %d token id", haveTok, wantTok)
				}
			}
		})
	}
}

func TestTokenizer_Scan(t *testing.T) {
	testCases := []struct {
		name string
		want []Token
		str  string
	}{
		{
			name: "UNUSED",
			want: []Token{UNUSED, UNUSED, UNUSED},
			str:  "a b c ~",
		},
		{
			name: "TRIGGER,WHEN,VARS,REPEAT,RESET,AFTER",
			want: []Token{TRIGGER, WHEN, VARS, REPEAT, RESET, AFTER},
			str:  "TRIGGER WHEN VARS REPEAT RESET AFTER",
		},
		{
			name: "INT,FLOAT,STRING",
			want: []Token{INT, FLOAT, STRING},
			str:  "1 1.1 \"ok\"",
		},
		{
			name: "ASSIGN,SEMICOLON,LPAREN,RPAREN,COMMA,LBRACK,RBRACK,QUO",
			want: []Token{ASSIGN, SEMICOLON, LPAREN, RPAREN, COMMA, LBRACK, RBRACK, QUO},
			str:  "= ; ( ) , [ ] /",
		},
		{
			name: "GEQ,LEQ,NEQ,GTR,LSS",
			want: []Token{GEQ, LEQ, NEQ, GTR, LSS},
			str:  ">= <= != > <",
		},
		{
			name: "LAND,AND,LOR,OR",
			want: []Token{LAND, AND, LOR, OR},
			str:  "&& and || or not",
		},
		{
			name: "RESET,AFTER,INT,UNUSED",
			want: []Token{RESET, AFTER, INT, UNUSED},
			str:  "reset after 24H",
		},
		{
			name: "REPEAT,EVERY,INT,UNUSED",
			want: []Token{REPEAT, INT, UNUSED},
			str:  "repeat  24H",
		},
		{
			name: "REPEAT,INT,TIMES,INTERVAL,INT,UNUSED",
			want: []Token{REPEAT, INT, TIMES, INTERVAL, INT, UNUSED},
			str:  "repeat 25 times interval 10s",
		},
		{
			name: "MUL, BETWEEN, NOTBETWEEN, COLON",
			want: []Token{MUL, BETWEEN, NOTBETWEEN, COLON},
			str:  "* between not between :",
		},
		{
			name: "ILLEGAL",
			want: []Token{ILLEGAL, ILLEGAL, ILLEGAL},
			str:  "!! |> &",
		},
		{
			name: "ILLEGAL",
			want: []Token{ILLEGAL},
			str:  "&%",
		},
	}
	for _, tc := range testCases {
		tokenizer := NewTokenizer(strings.NewReader(tc.str))
		for i := 0; i < len(tc.want); i++ {
			name := fmt.Sprintf("%s_%s", tc.name, KeywordString(tc.want[i]))
			t.Run(name, func(t *testing.T) {
				tok, _ := tokenizer.Scan()
				if have, want := tok, tc.want[i]; have != want {
					t.Fatalf("have %d, want %d token id", have, want)
				}
			})
		}
	}
}

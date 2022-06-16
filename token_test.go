package geoqlparser

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenizer_Scan(t *testing.T) {
	testCases := []struct {
		name string
		want []Token
		str  string
	}{
		{
			name: "GEOMETRY",
			want: []Token{GEOMETRY_POINT, GEOMETRY_MULTIPOINT,
				GEOMETRY_LINE, GEOMETRY_MULTILINE,
				GEOMETRY_POLYGON, GEOMETRY_MULTIPOLYGON, GEOMETRY_CIRCLE, GEOMETRY_COLLECTION},
			str: "point multipoint line multiline polygon multipolygon circle collection",
		},
		{
			name: "SELECTOR",
			want: []Token{SELECTOR, SELECTOR, SELECTOR},
			str:  "someField speed index",
		},
		{
			name: "BOOLEAN",
			want: []Token{BOOLEAN, BOOLEAN, BOOLEAN, BOOLEAN},
			str:  "true up false down",
		},
		{
			name: "NEARBY, NOTNEARBY, INTERSECTS, NOTINTERSECTS, WITHIN, NOTWITHIN",
			want: []Token{NEARBY, NOTNEARBY, INTERSECTS, NOTINTERSECTS, WITHIN, NOTWITHIN},
			str:  "NEARBY NOT NEARBY INTERSECTS  NOT INTERSECTS  WITHIN  NOT WITHIN",
		},
		{
			name: "EQL, LEQL, NEQ, IN, NOT_IN SUB ADD",
			want: []Token{EQL, LEQL, NEQ, LNEQ, IN, NOT_IN, SUB, ADD},
			str:  "eq == not eq != in not in - +",
		},
		{
			name: "TRIGGER,WHEN,SET,REPEAT,RESET,AFTER",
			want: []Token{TRIGGER, WHEN, SET, REPEAT, RESET},
			str:  "TRIGGER WHEN SET REPEAT RESET ",
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
			name: "GEQ,LEQ,LNEQ,GTR,LSS",
			want: []Token{GEQ, LEQ, LNEQ, GTR, LSS},
			str:  ">= <= != > <",
		},
		{
			name: "LAND,AND,LOR,OR",
			want: []Token{LAND, AND, LOR, OR},
			str:  "&& and || or not",
		},
		{
			name: "REPEAT,EVERY,INT,UNUSED",
			want: []Token{REPEAT, INT},
			str:  "repeat  24H",
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
					t.Fatalf("have %s, want %s", KeywordString(have), KeywordString(want))
				}
			})
		}
	}
}

//func TestTokenizer_ScanNMEA(t *testing.T) {
//	testCases := []struct {
//		name string
//		set  map[string]Token
//	}{
//		{
//			name: "NMEA_AAM_*",
//			set: map[string]Token{
//				"nmea_aam_statusArrivalCircleEntered": NMEA_AAM_STATUS_ARRIVAL_CIRCLE_ENTERED,
//				"nmea_aam_statusPerpendicularPassed":  NMEA_AAM_STATUS_PERPENDICULAR_PASSED,
//				"nmea_aam_arrivalCircleRadius":        NMEA_AAM_ARRIVAL_CIRCLE_RADIUS,
//				"nmea_aam_arrivalCircleRadiusUnit":    NMEA_AAM_ARRIVAL_CIRCLE_UNIT,
//				"nmea_aam_destinationWaypointID":      NMEA_AAM_DESTINATION_WAYPOINT_ID,
//			},
//		},
//		{
//			name: "NMEA_ALA_*",
//			set: map[string]Token{
//				"nmea_ala_time":               NMEA_ALA_TIME,
//				"nmea_ala_systemIndicator":    NMEA_ALA_SYSTEM_INDICATOR,
//				"nmea_ala_subSystemIndicator": NMEA_ALA_SUB_SYSTEM_INDICATOR,
//				"nmea_ala_instanceNumber":     NMEA_ALA_INSTANCE_NUMBER,
//				"nmea_ala_type":               NMEA_ALA_TYPE,
//				"nmea_ala_condition":          NMEA_ALA_CONDITION,
//				"nmea_ala_alarmAckState":      NMEA_ALA_ALARM_ACK_STATE,
//				"nmea_ala_message":            NMEA_ALA_MESSAGE,
//			},
//		},
//		{
//			name: "NMEA_APB_*",
//			set: map[string]Token{
//				"nmea_apb_statusGeneralWarning":       NMEA_APB_STATUS_GENERAL_WARNING,
//				"nmea_apb_statusLockWarning":          NMEA_APB_STATUS_LOCK_WARNING,
//				"nmea_apb_crossTrackErrorMagnitude":   NMEA_APB_CROSS_TRACK_ERROR_MAGNITUDE,
//				"nmea_apb_directionToSteer":           NMEA_APB_DIRECTION_TO_STEER,
//				"nmea_apb_crossTrackUnits":            NMEA_APB_CROSS_TRACK_UNITS,
//				"nmea_apb_statusArrivalCircleEntered": NMEA_APB_STATUS_ARRIVAL_CIRCLE_ENTERED,
//				"nmea_apb_statusPerpendicularPassed":  NMEA_APB_STATUS_PERPENDICULAR_PASSED,
//				"nmea_apb_bearingOriginToDest":        NMEA_APB_BEARING_ORIGIN_TO_DEST,
//				"nmea_apb_bearingOriginToDestType":    NMEA_APB_BEARING_ORIGIN_DEST_TYPE,
//				"nmea_apb_destinationWaypointID":      NMEA_APB_DESTINATION_WAYPOINT_ID,
//				"nmea_apb_bearingPresentToDest":       NMEA_APB_BEARING_PRESENT_TO_DEST,
//				"nmea_apb_bearingPresentToDestType":   NMEA_APB_BEARING_PRESENT_TO_DEST_TYPE,
//				"nmea_apb_heading":                    NMEA_APB_HEADING,
//				"nmea_apb_headingType":                NMEA_APB_HEADING_TYPE,
//				"nmea_apb_ffaMode":                    NMEA_APB_FFA_MODE,
//			},
//		},
//		{
//			name: "NMEA_BEC_*",
//			set: map[string]Token{
//				"nmea_bec_time":                       NMEA_BEC_TIME,
//				"nmea_bec_latitude":                   NMEA_BEC_LATITUDE,
//				"nmea_bec_longitude":                  NMEA_BEC_LONGITUDE,
//				"nmea_bec_bearingTrue":                NMEA_BEC_BEARING_TRUE,
//				"nmea_bec_bearingTrueValid":           NMEA_BEC_BEARING_TRUE_VALID,
//				"nmea_bec_bearingMagnetic":            NMEA_BEC_BEARING_MAGNETIC,
//				"nmea_bec_bearingMagneticValid":       NMEA_BEC_BEARING_MAGNETIC_VALID,
//				"nmea_bec_distanceNauticalMiles":      NMEA_BEC_DISTANCE_NAUTICAL_MILES,
//				"nmea_bec_distanceNauticalMilesValid": NMEA_BEC_DISTANCE_NAUTICAL_MILES_VALID,
//				"nmea_bec_destinationWaypointID":      NMEA_BEC_DESTINATION_WAYPOINT_ID,
//			},
//		},
//		{
//			name: "NMEA_BOD_*",
//			set: map[string]Token{
//				"nmea_bod_bearingTrue":           NMEA_BOD_BEARING_TRUE,
//				"nmea_bod_bearingTrueType":       NMEA_BOD_BEARING_TRUE_TYPE,
//				"nmea_bod_bearingMagnetic":       NMEA_BOD_BEARING_MAGNETIC,
//				"nmea_bod_bearingMagneticType":   NMEA_BOD_BEARING_MAGNETIC_TYPE,
//				"nmea_bod_destinationWaypointID": NMEA_BOD_DESTINATION_WAYPOINT_ID,
//				"nmea_bod_originWaypointID":      NMEA_BOD_ORIGIN_WAYPOINT_ID,
//			},
//		},
//		{
//			name: "NMEA_BWC_*",
//			set: map[string]Token{
//				"nmea_bwc_time":                      NMEA_BWC_TIME,
//				"nmea_bwc_latitude":                  NMEA_BWC_LATITUDE,
//				"nmea_bwc_longitude":                 NMEA_BWC_LONGITUDE,
//				"nmea_bwc_bearingTrue":               NMEA_BWC_BEARING_TRUE,
//				"nmea_bwc_bearingTrueType":           NMEA_BWC_BEARING_TRUE_TYPE,
//				"nmea_bwc_bearingMagnetic":           NMEA_BWC_BEARING_MAGNETIC,
//				"nmea_bwc_bearingMagneticType":       NMEA_BWC_BEARING_MAGNETIC_TYPE,
//				"nmea_bwc_distanceNauticalMiles":     NMEA_BWC_DISTANCE_NAUTICAL_MILES,
//				"nmea_bwc_distanceNauticalMilesUnit": NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT,
//			},
//		},
//		{
//			name: "NMEA_BOD_*",
//			set: map[string]Token{
//				"nmea_bod_bearingTrue":           NMEA_BOD_BEARING_TRUE,
//				"nmea_bod_bearingTrueType":       NMEA_BOD_BEARING_TRUE_TYPE,
//				"nmea_bod_bearingMagnetic":       NMEA_BOD_BEARING_MAGNETIC,
//				"nmea_bod_bearingMagneticType":   NMEA_BOD_BEARING_MAGNETIC_TYPE,
//				"nmea_bod_destinationWaypointID": NMEA_BOD_DESTINATION_WAYPOINT_ID,
//				"nmea_bod_originWaypointID":      NMEA_BOD_ORIGIN_WAYPOINT_ID,
//			},
//		},
//		{
//			name: "NMEA_BWC_*",
//			set: map[string]Token{
//				"nmea_bwc_time":                      NMEA_BWC_TIME,
//				"nmea_bwc_latitude":                  NMEA_BWC_LATITUDE,
//				"nmea_bwc_longitude":                 NMEA_BWC_LONGITUDE,
//				"nmea_bwc_bearingTrue":               NMEA_BWC_BEARING_TRUE,
//				"nmea_bwc_bearingTrueType":           NMEA_BWC_BEARING_TRUE_TYPE,
//				"nmea_bwc_bearingMagnetic":           NMEA_BWC_BEARING_MAGNETIC,
//				"nmea_bwc_bearingMagneticType":       NMEA_BWC_BEARING_MAGNETIC_TYPE,
//				"nmea_bwc_distanceNauticalMiles":     NMEA_BWC_DISTANCE_NAUTICAL_MILES,
//				"nmea_bwc_distanceNauticalMilesUnit": NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT,
//			},
//		},
//		{
//			name: "NMEA_BWC_*",
//			set: map[string]Token{
//				"nmea_bwr_time":                  NMEA_BWR_TIME,
//				"nmea_bwr_latitude":              NMEA_BWR_LATITUDE,
//				"nmea_bwr_longitude":             NMEA_BWR_LONGITUDE,
//				"nmea_bwr_bearingTrue":           NMEA_BWR_BEARING_TRUE,
//				"nmea_bwr_bearingTrueType":       NMEA_BWR_BEARING_TRUE_TYPE,
//				"nmea_bwr_bearingMagnetic":       NMEA_BWR_BEARING_MAGNETIC,
//				"nmea_bwr_bearingMagneticType":   NMEA_BWR_BEARING_MAGNETIC_TYPE,
//				"nmea_bwr_distanceNauticalMiles": NMEA_BWR_DISTANCE_NAUTICAL_MILES,
//				"nmea_bwr_destinationWaypointID": NMEA_BWR_DESTINATION_WAYPOINT_ID,
//				"nmea_bwr_ffaMode":               NMEA_BWR_FFA_MODE,
//			},
//		},
//		{
//			name: "NMEA_BWW_*",
//			set: map[string]Token{
//				"nmea_bww_bearingTrue":           NMEA_BWW_BEARING_TRUE,
//				"nmea_bww_bearingTrueType":       NMEA_BWW_BEARING_TRUE_TYPE,
//				"nmea_bww_bearingMagnetic":       NMEA_BWW_BEARING_MAGNETIC,
//				"nmea_bww_bearingMagneticType":   NMEA_BWW_BEARING_MAGNETIC_TYPE,
//				"nmea_bww_destinationWaypointID": NMEA_BWW_DESTINATION_WAYPOINT_ID,
//				"nmea_bww_originWaypointID":      NMEA_BWW_ORIGINAL_WAYPOINT_ID,
//			},
//		},
//		{
//			name: "NMEA_DBK_*",
//			set: map[string]Token{
//				"nmea_dbk_depthFeet":        NMEA_DBK_DEPTH_FEET,
//				"nmea_dbk_depthFeetUnit":    NMEA_DBK_DEPTH_FEET_UNIT,
//				"nmea_dbk_depthMeters":      NMEA_DBK_DEPTH_METERS,
//				"nmea_dbk_depthMetersUnit":  NMEA_DBK_DEPTH_METERS_UNIT,
//				"nmea_dbk_depthFathoms":     NMEA_DBK_DEPTH_FATHOMS,
//				"nmea_dbk_depthFathomsUnit": NMEA_DBK_DEPTH_FATHOMS_UNIT,
//			},
//		},
//		{
//			name: "NMEA_DBS_*",
//			set: map[string]Token{
//				"nmea_dbs_depthFeet":       NMEA_DBS_DEPTH_FEET,
//				"nmea_dbs_depthFeetUnit":   NMEA_DBS_DEPTH_FEET_UNIT,
//				"nmea_dbs_depthMeters":     NMEA_DBS_DEPTH_METERS,
//				"nmea_dbs_depthMeterUnit":  NMEA_DBS_DEPTH_METERS_UNIT,
//				"nmea_dbs_depthFathoms":    NMEA_DBS_DEPTH_FATHOMS,
//				"nmea_dbs_depthFathomUnit": NMEA_DBS_DEPTH_FATHOMS_UNIT,
//			},
//		},
//		{
//			name: "NMEA_DBT_*",
//			set: map[string]Token{
//				"nmea_dbt_depthFeet":    NMEA_DBT_DEPTH_FEET,
//				"nmea_dbt_depthMeters":  NMEA_DBT_DEPTH_METERS,
//				"nmea_dbt_depthFathoms": NMEA_DBT_DEPTH_FATHOMS,
//			},
//		},
//		{
//			name: "NMEA_DOR_*",
//			set: map[string]Token{
//				"nmea_dor_type":               NMEA_DOR_TYPE,
//				"nmea_dor_time":               NMEA_DOR_TIME,
//				"nmea_dor_systemIndicator":    NMEA_DOR_SYSTEM_INDICATOR,
//				"nmea_dor_divisionIndicator1": NMEA_DOR_DIVISION_INDICATOR1,
//				"nmea_dor_divisionIndicator2": NMEA_DOR_DIVISION_INDICATOR2,
//				"nmea_dor_doorNumberOrCount":  NMEA_DOR_DOOR_NUMBER_OR_COUNT,
//				"nmea_dor_doorStatus":         NMEA_DOR_DOOR_STATUS,
//				"nmea_dor_switchSetting":      NMEA_DOR_SWITCH_SETTING,
//				"nmea_dor_message":            NMEA_DOR_MESSAGE,
//			},
//		},
//		{
//			name: "NMEA_DPT_*",
//			set: map[string]Token{
//				"nmea_dpt_depth":      NMEA_DPT_DEPTH,
//				"nmea_dpt_offset":     NMEA_DPT_OFFSET,
//				"nmea_dpt_rangeScale": NMEA_DPT_RANGE_SCALE,
//			},
//		},
//		{
//			name: "NMEA_DSC_*",
//			set: map[string]Token{
//				"nmea_dsc_formatSpecifier":             NMEA_DSC_FORMAT_SPECIFIER,
//				"nmea_dsc_address":                     NMEA_DSC_ADDRESS,
//				"nmea_dsc_category":                    NMEA_DSC_CATEGORY,
//				"nmea_dsc_distressCauseOrTeleCommand1": NMEA_DSC_DISTRESS_CAUSE_OR_TELE_CMD1,
//				"nmea_dsc_commandTypeOrTeleCommand2":   NMEA_DSC_CMD_TYPE_OR_TELE_CMD2,
//				"nmea_dsc_positionOrCanal":             NMEA_DSC_POSITION_OR_CANAL,
//				"nmea_dsc_timeOrTelephoneNumber":       NMEA_DSC_TIMER_OR_TELE_NUMBER,
//				"nmea_dsc_mmsi":                        NMEA_DSC_MMSI,
//				"nmea_dsc_distressCause":               NMEA_DSC_DISTREES_CAUSE,
//				"nmea_dsc_acknowledgement":             NMEA_DSC_ACK_LEDGEMENT,
//				"nmea_dsc_expansionIndicator":          NMEA_DSC_EXPANSION_INDICATOR,
//			},
//		},
//		{
//			name: "NMEA_DSE_*",
//			set: map[string]Token{
//				"nmea_dse_totalNumber":     NMEA_DSE_TOTAL_NUMBER,
//				"nmea_dse_number":          NMEA_DSE_NUMBER,
//				"nmea_dse_acknowledgement": NMEA_DSE_ACK_LEDGEMENT,
//				"nmea_dse_mmsi":            NMEA_DSE_MMSI,
//			},
//		},
//		{
//			name: "NMEA_DTM_*",
//			set: map[string]Token{
//				"nmea_dtm_localDatumCode":        NMEA_DTM_LOCAL_DATUM_CODE,
//				"nmea_dtm_localDatumSubcode":     NMEA_DTM_LOCAL_DATUM_SUBCODE,
//				"nmea_dtm_latitudeOffsetMinute":  NMEA_DTM_LATITUDE_OFFSET_MINUTE,
//				"nmea_dtm_longitudeOffsetMinute": NMEA_DTM_LONGITUDE_OFFSET_MINUTE,
//				"nmea_dtm_altitudeOffsetMeters":  NMEA_DTM_ALTITUDE_OFFSET_METERS,
//				"nmea_dtm_datumName":             NMEA_DTM_DATUM_NAME,
//			},
//		},
//		{
//			name: "NMEA_EVE_*",
//			set: map[string]Token{
//				"nmea_eve_time":    NMEA_EVE_TIME,
//				"nmea_eve_tagCode": NMEA_EVE_TAG_CODE,
//				"nmea_eve_message": NMEA_EVE_MESSAGE,
//			},
//		},
//		{
//			name: "NMEA_FIR_*",
//			set: map[string]Token{
//				"nmea_fir_type":                      NMEA_FIR_TYPE,
//				"nmea_fir_time":                      NMEA_FIR_TIME,
//				"nmea_fir_systemIndicator":           NMEA_FIR_SYSTEM_INDICATOR,
//				"nmea_fir_divisionIndicator1":        NMEA_FIR_DIVISION_INDICATOR1,
//				"nmea_fir_divisionIndicator2":        NMEA_FIR_DIVISION_INDICATOR2,
//				"nmea_fir_fireDetectorNumberOrCount": NMEA_FIR_FIRE_DETECTOR_NUM_OR_COUNT,
//				"nmea_fir_condition":                 NMEA_FIR_CONDITION,
//				"nmea_fir_alarmAckState":             NMEA_FIR_ALARAM_ACK_STATE,
//				"nmea_fir_message":                   NMEA_FIR_MESSAGE,
//			},
//		},
//		{
//			name: "NMEA_GGA_*",
//			set: map[string]Token{
//				"nmea_gga_time":          NMEA_GGA_TIME,
//				"nmea_gga_latitude":      NMEA_GGA_LATITUDE,
//				"nmea_gga_longitude":     NMEA_GGA_LONGITUDE,
//				"nmea_gga_fixQuality":    NMEA_GGA_FIX_QUOLITY,
//				"nmea_gga_numSatellites": NMEA_GGA_NUM_SATELLITES,
//				"nmea_gga_hdop":          NMEA_GGA_HDOP,
//				"nmea_gga_altitude":      NMEA_GGA_ALTITUDE,
//				"nmea_gga_separation":    NMEA_GGA_SEPARATION,
//				"nmea_gga_dgspAge":       NMEA_GGA_DGPS_AGE,
//				"nmea_gga_dgspId":        NMEA_GGA_DGSP_ID,
//			},
//		},
//		{
//			name: "NMEA_GGL_*",
//			set: map[string]Token{
//				"nmea_gll_latitude":  NMEA_GLL_LATITUDE,
//				"nmea_gll_longitude": NMEA_GLL_LONGITUDE,
//				"nmea_gll_time":      NMEA_GLL_TIME,
//				"nmea_gll_validity":  NMEA_GLL_VALIDITY,
//				"nmea_gll_ffaMode":   NMEA_GLL_FFAMODE,
//			},
//		},
//		{
//			name: "NMEA_GNS_*",
//			set: map[string]Token{
//				"nmea_gns_time":       NMEA_GNS_TIME,
//				"nmea_gns_latitude":   NMEA_GNS_LATITUDE,
//				"nmea_gns_longitude":  NMEA_GNS_LONGITUDE,
//				"nmea_gns_altitude":   NMEA_GNS_ALTITUDE,
//				"nmea_gns_mode":       NMEA_GNS_MODE,
//				"nmea_gns_svs":        NMEA_GNS_SVS,
//				"nmea_gns_hdop":       NMEA_GNS_HDOP,
//				"nmea_gns_separation": NMEA_GNS_SEPARATION,
//				"nmea_gns_age":        NMEA_GNS_AGE,
//				"nmea_gns_station":    NMEA_GNS_STATION,
//				"nmea_gns_navStatus":  NMEA_GNS_NAV_STATUS,
//			},
//		},
//		{
//			name: "NMEA_GSA_*",
//			set: map[string]Token{
//				"nmea_gsa_mode":     NMEA_GSA_MODE,
//				"nmea_gsa_fixType":  NMEA_GSA_FIX_TYPE,
//				"nmea_gsa_sv":       NMEA_GSA_SV,
//				"nmea_gsa_pdop":     NMEA_GSA_PDOP,
//				"nmea_gsa_hdop":     NMEA_GSA_HDOP,
//				"nmea_gsa_vdop":     NMEA_GSA_VDOP,
//				"nmea_gsa_systemID": NMEA_GSA_SYSTEM_ID,
//			},
//		},
//		{
//			name: "NMEA_HDG_*",
//			set: map[string]Token{
//				"nmea_hdg_heading":            NMEA_HDG_HEADING,
//				"nmea_hdg_deviation":          NMEA_HDG_DEVIATION,
//				"nmea_hdg_deviationDirection": NMEA_HDG_DEVIATION_DIRECTION,
//				"nmea_hdg_variation":          NMEA_HDG_VARIATION,
//				"nmea_hdg_variationDirection": NMEA_HDG_VARIATION_DIRECTION,
//			},
//		},
//		{
//			name: "NMEA_HDM_*",
//			set: map[string]Token{
//				"nmea_hdm_heading":       NMEA_HDM_HEADING,
//				"nmea_hdm_magneticValid": NMEA_HDM_MAGNETIC_VALID,
//			},
//		},
//		{
//			name: "NMEA_HDT_*",
//			set: map[string]Token{
//				"nmea_hdt_heading": NMEA_HDT_HEADING,
//				"nmea_hdt_true":    NMEA_HDT_TRUE,
//			},
//		},
//		{
//			name: "NMEA_HSC_*",
//			set: map[string]Token{
//				"nmea_hsc_trueHeading":         NMEA_HSC_TRUE_HEADING,
//				"nmea_hsc_trueHeadingType":     NMEA_HSC_TRUE_HEADING_TYPE,
//				"nmea_hsc_magneticHeading":     NMEA_HSC_MAGNETIC_HEADING,
//				"nmea_hsc_magneticHeadingType": NMEA_HSC_MAGNETIC_HEADING_TYPE,
//			},
//		},
//		{
//			name: "NMEA_MDA_*",
//			set: map[string]Token{
//				"nmea_mda_pressureInch":          NMEA_MDA_PRESSURE_INCH,
//				"nmea_mda_inchesValid":           NMEA_MDA_INCHES_VALID,
//				"nmea_mda_pressureBar":           NMEA_MDA_PRESSURE_BAR,
//				"nmea_mda_barsValid":             NMEA_MDA_BARS_VALID,
//				"nmea_mda_airTemp":               NMEA_MDA_AIR_TEMP,
//				"nmea_mda_airTempValid":          NMEA_MDA_AIR_TEMP_VALID,
//				"nmea_mda_waterTemp":             NMEA_MDA_WATER_TEMP,
//				"nmea_mda_waterTempValid":        NMEA_MDA_WATER_TEMP_VALID,
//				"nmea_mda_relativeHum":           NMEA_MDA_RELATIVE_HUM,
//				"nmea_mda_absoluteHum":           NMEA_MDA_ABSOLUTE_HUM,
//				"nmea_mda_dewPoint":              NMEA_MDA_DEW_POINT,
//				"nmea_mda_dewPointValid":         NMEA_MDA_DEW_POINT_VALID,
//				"nmea_mda_windDirectionTrue":     NMEA_MDA_WIND_DIRECTION_TRUE,
//				"nmea_mda_trueValid":             NMEA_MDA_TRUE_VALID,
//				"nmea_mda_windDirectionMagnetic": NMEA_MDA_WIND_DIRECTION_MAGNETIC,
//				"nmea_mda_magneticValid":         NMEA_MDA_MAGNETIC_VALID,
//				"nmea_mda_windSpeedKnots":        NMEA_MDA_WIND_SPEED_KNOTS,
//				"nmea_mda_knotsValid":            NMEA_MDA_KNOTS_VALID,
//				"nmea_mda_windSpeedMeters":       NMEA_MDA_WIND_SPEED_METERS,
//				"nmea_mda_metersValid":           NMEA_MDA_METERS_VALID,
//			},
//		},
//		{
//			name: "NMEA_MTA_*",
//			set: map[string]Token{
//				"nmea_mta_temperature": NMEA_MTA_TEMPERATURE,
//				"nmea_mta_unit":        NMEA_MTA_UNIT,
//			},
//		},
//		{
//			name: "NMEA_MTK_*",
//			set: map[string]Token{
//				"nmea_mtk_command": NMEA_MTK_COMMAND,
//				"nmea_mtk_flag":    NMEA_MTK_FLAG,
//			},
//		},
//		{
//			name: "NMEA_MTW_*",
//			set: map[string]Token{
//				"nmea_mtw_temperature":  NMEA_MTW_TEMPERATURE,
//				"nmea_mtw_celsiusValid": NMEA_MTW_CELSIUS_VALID,
//			},
//		},
//		{
//			name: "NMEA_MWD_*",
//			set: map[string]Token{
//				"nmea_mwd_windDirectionTrue":     NMEA_MWD_WIND_DIRECTION_TRUE,
//				"nmea_mwd_trueValid":             NMEA_MWD_TRUE_VALID,
//				"nmea_mwd_windDirectionMagnetic": NMEA_MWD_WIND_DIRECTION_MAGNETIC,
//				"nmea_mwd_magneticValid":         NMEA_MWD_MAGNETIC_VALID,
//				"nmea_mwd_windSpeedKnots":        NMEA_MWD_WIND_SPEED_KNOTS,
//				"nmea_mwd_knotsValid":            NMEA_MWD_KNOTS_VALID,
//				"nmea_mwd_windSpeedMeters":       NMEA_MWD_WIND_SPEED_METERS,
//				"nmea_mwd_metersValid":           NMEA_MWD_METERS_VALID,
//			},
//		},
//		{
//			name: "NMEA_MWV_*",
//			set: map[string]Token{
//				"nmea_mwv_windAngle":     NMEA_MWV_WIND_ANGLE,
//				"nmea_mwv_reference":     NMEA_MWV_REFERENCE,
//				"nmea_mwv_windSpeed":     NMEA_MWV_WIND_SPEED,
//				"nmea_mwv_windSpeedUnit": NMEA_MWV_WIND_SPEED_UNIT,
//				"nmea_mwv_statusValid":   NMEA_MWV_STATUS_VALID,
//			},
//		},
//		{
//			name: "NMEA_OSD_*",
//			set: map[string]Token{
//				"nmea_osd_heading":          NMEA_OSD_HEADING,
//				"nmea_osd_headingStatus":    NMEA_OSD_HEADING_STATUS,
//				"nmea_osd_vesselTrueCourse": NMEA_OSD_VESSEL_TRUE_COURSE,
//				"nmea_osd_courseReference":  NMEA_OSD_COURSE_REFERENCE,
//				"nmea_osd_vesselSpeed":      NMEA_OSD_VESSEL_SPEED,
//				"nmea_osd_speedReference":   NMEA_OSD_SPEED_REFERENCE,
//				"nmea_osd_vesselSetTrue":    NMEA_OSD_VESSEL_SET_TRUE,
//				"nmea_osd_vesselDrift":      NMEA_OSD_VESSEL_DRIFT,
//				"nmea_osd_speedUnits":       NMEA_OSD_SPEED_UNITS,
//			},
//		},
//		{
//			name: "NMEA_PGRME_*",
//			set: map[string]Token{
//				"nmea_pgrme_horizontal": NMEA_PGRME_HORIZONTAL,
//				"nmea_pgrme_vertical":   NMEA_PGRME_VERTICAL,
//				"nmea_pgrme_spherical":  NMEA_PGRME_SPHERICAL,
//			},
//		},
//		{
//			name: "NMEA_PHTRO_*",
//			set: map[string]Token{
//				"nmea_phtro_pitch": NMEA_PHTRO_PITCH,
//				"nmea_phtro_bow":   NMEA_PHTRO_BOW,
//				"nmea_phtro_roll":  NMEA_PHTRO_ROLL,
//				"nmea_phtro_port":  NMEA_PHTRO_PORT,
//			},
//		},
//		{
//			name: "NMEA_PRDID_*",
//			set: map[string]Token{
//				"nmea_prdid_pitch":   NMEA_PRDID_PITCH,
//				"nmea_prdid_roll":    NMEA_PRDID_ROLL,
//				"nmea_prdid_heading": NMEA_PRDID_HEADING,
//			},
//		},
//		{
//			name: "NMEA_PSKPDPT_*",
//			set: map[string]Token{
//				"nmea_pskpdpt_depth":              NMEA_PSKPDPT_DEPTH,
//				"nmea_pskpdpt_offset":             NMEA_PSKPDPT_OFFSET,
//				"nmea_pskpdpt_rangeScale":         NMEA_PSKPDPT_RANGE_SCALE,
//				"nmea_pskpdpt_bottomEchoStrength": NMEA_PSKPDPT_BOTTOM_ECHO_STRENGTH,
//				"nmea_pskpdpt_channelNumber":      NMEA_PSKPDPT_CHANNEL_NUMBER,
//				"nmea_pskpdpt_transducerLocation": NMEA_PSKPDPT_TRANSDUCER_LOCATION,
//			},
//		},
//		{
//			name: "NMEA_PSONCMS_*",
//			set: map[string]Token{
//				"nmea_psoncms_quaternion0":       NMEA_PSONCMS_QUATERNION0,
//				"nmea_psoncms_quaternion1":       NMEA_PSONCMS_QUATERNION1,
//				"nmea_psoncms_quaternion2":       NMEA_PSONCMS_QUATERNION2,
//				"nmea_psoncms_quaternion3":       NMEA_PSONCMS_QUATERNION3,
//				"nmea_psoncms_accelerationX":     NMEA_PSONCMS_ACCELERATION_X,
//				"nmea_psoncms_accelerationY":     NMEA_PSONCMS_ACCELERATION_Y,
//				"nmea_psoncms_accelerationZ":     NMEA_PSONCMS_ACCELERATION_Z,
//				"nmea_psoncms_rateOfTurnX":       NMEA_PSONCMS_RATE_OF_TURN_X,
//				"nmea_psoncms_rateOfTurnY":       NMEA_PSONCMS_RATE_OF_TURN_Z,
//				"nmea_psoncms_magneticFieldX":    NMEA_PSONCMS_MAGNETIC_FIELD_X,
//				"nmea_psoncms_magneticFieldY":    NMEA_PSONCMS_MAGNETIC_FIELD_Y,
//				"nmea_psoncms_rateOfTurnZ":       NMEA_PSONCMS_MAGNETIC_FIELD_Z,
//				"nmea_psoncms_sensorTemperature": NMEA_PSONCMS_SENSOR_TEMPERATURE,
//			},
//		},
//		{
//			name: "NMEA_RMB_*",
//			set: map[string]Token{
//				"nmea_rmb_dataStatus":                      NMEA_RMB_DATA_STATUS,
//				"nmea_rmb_crossTrackErrorNauticalMiles":    NMEA_RMB_CROSS_TRACK_ERROR_NAUTICAL_MILES,
//				"nmea_rmb_directionToSteer":                NMEA_RMB_DIRECTION_TO_STEER,
//				"nmea_rmb_originWaypointID":                NMEA_RMB_ORIGIN_WAYPOINT_ID,
//				"nmea_rmb_destinationWaypointID":           NMEA_RMB_DESTINATION_WAYPOINT_ID,
//				"nmea_rmb_destinationLatitude":             NMEA_RMB_DESTINATION_LATITUDE,
//				"nmea_rmb_destinationLongitude":            NMEA_RMB_DESTINATION_LONGITUDE,
//				"nmea_rmb_rangeToDestinationNauticalMiles": NMEA_RMB_RANGE_TO_DESTINATION_NAUTICAL_MILES,
//				"nmea_rmb_trueBearingToDestination":        NMEA_RMB_TRUE_BEARING_TO_DESTINATION,
//				"nmea_rmb_velocityToDestinationKnots":      NMEA_RMB_VELOCITY_TO_DESTINATION_KNOTS,
//				"nmea_rmb_arrivalStatus":                   NMEA_RMB_ARRIVAL_STATUS,
//				"nmea_rmb_ffaMode":                         NMEA_RMB_FFAMODE,
//			},
//		},
//		{
//			name: "NMEA_RMC_*",
//			set: map[string]Token{
//				"nmea_rmc_time":      NMEA_RMC_TIME,
//				"nmea_rmc_validity":  NMEA_RMC_VALIDITY,
//				"nmea_rmc_latitude":  NMEA_RMC_LATITUDE,
//				"nmea_rmc_longitude": NMEA_RMC_LONGITUDE,
//				"nmea_rmc_speed":     NMEA_RMC_SPEED,
//				"nmea_rmc_course":    NMEA_RMC_COURSE,
//				"nmea_rmc_date":      NMEA_RMC_DATE,
//				"nmea_rmc_variation": NMEA_RMC_VARIATION,
//				"nmea_rmc_ffaMode":   NMEA_RMC_FFAMODE,
//				"nmea_rmc_navStatus": NMEA_RMC_NAV_STATUS,
//			},
//		},
//		{
//			name: "NMEA_ROT_*",
//			set: map[string]Token{
//				"nmea_rot_rateOfTurn": NMEA_ROT_RATE_OF_TURN,
//				"nmea_rot_valid":      NMEA_ROT_VALID,
//			},
//		},
//		{
//			name: "NMEA_RPM_*",
//			set: map[string]Token{
//				"nmea_rpm_source":       NMEA_RPM_SOURCE,
//				"nmea_rpm_engineNumber": NMEA_RPM_ENGINE_NUMBER,
//				"nmea_rpm_speedRPM":     NMEA_RPM_SPEED_RPM,
//				"nmea_rpm_pitchPercent": NMEA_RPM_PITCH_PERCENT,
//				"nmea_rpm_status":       NMEA_RPM_STATUS,
//			},
//		},
//		{
//			name: "NMEA_RSA_*",
//			set: map[string]Token{
//				"nmea_rsa_starboardRudderAngle":       NMEA_RSA_STARBOARD_RUDDER_ANGLE,
//				"nmea_rsa_starboardRudderAngleStatus": NMEA_RSA_STARBOARD_RUDDER_ANGLE_STATUS,
//				"nmea_rsa_portRudderAngle":            NMEA_RSA_PORT_RUDDER_ANGLE,
//				"nmea_rsa_portRudderAngleStatus":      NMEA_RSA_PORT_RUDDER_ANGLE_STATUS,
//			},
//		},
//		{
//			name: "NMEA_RSD_*",
//			set: map[string]Token{
//				"nmea_rsd_origin1Range":           NMEA_RSD_ORIGIN1_RANGE,
//				"nmea_rsd_origin1Bearing":         NMEA_RSD_ORIGIN1_BEARING,
//				"nmea_rsd_variableRangeMarker1":   NMEA_RSD_VARIABLE_RANGE_MARKET1,
//				"nmea_rsd_bearingLine1":           NMEA_RSD_BEARING_LINE1,
//				"nmea_rsd_origin2Range":           NMEA_RSD_ORIGIN2_RANGE,
//				"nmea_rsd_origin2Bearing":         NMEA_RSD_ORIGIN2_BEARING,
//				"nmea_rsd_variableRangeMarker2":   NMEA_RSD_VARIABLE_RANGE_MARKET2,
//				"nmea_rsd_bearingLine2":           NMEA_RSD_BEARING_LINE2,
//				"nmea_rsd_cursorRangeFromOwnShip": NMEA_RSD_CURSOR_RANGE_FROM_OWN_SHIP,
//				"nmea_rsd_cursorBearingDegrees":   NMEA_RSD_CURSOR_BEARING_DEGREES,
//				"nmea_rsd_rangeScale":             NMEA_RSD_RANGE_SCALE,
//				"nmea_rsd_rangeUnit":              NMEA_RSD_RANGE_UNIT,
//				"nmea_rsd_displayRotation":        NMEA_RSD_DISPLAY_ROTATION,
//			},
//		},
//		{
//			name: "NMEA_RTE_*",
//			set: map[string]Token{
//				"nmea_rte_numberOfSentences":         NMEA_RTE_NUMBER_OF_SENTENCES,
//				"nmea_rte_sentenceNumber":            NMEA_RTE_SENTENCE_NUMBER,
//				"nmea_rte_activeRouteOrWaypointList": NMEA_RTE_ACTIVE_ROUTER_OR_WAYPOINT_LIST,
//				"nmea_rte_name":                      NMEA_RTE_NAME,
//				"nmea_rte_idents":                    NMEA_RTE_IDENTS,
//			},
//		},
//		{
//			name: "NMEA_THS_*",
//			set: map[string]Token{
//				"nmea_ths_heading": NMEA_THS_HEADING,
//				"nmea_ths_status":  NMEA_THS_STATUS,
//			},
//		},
//		{
//			name: "NMEA_TLL_*",
//			set: map[string]Token{
//				"nmea_tll_targetNumber":    NMEA_TLL_TARGET_NUMBER,
//				"nmea_tll_targetLatitude":  NMEA_TLL_TARGET_LATITUDE,
//				"nmea_tll_targetLongitude": NMEA_TLL_TARGET_LONGITUDE,
//				"nmea_tll_targetName":      NMEA_TLL_TARGET_NAME,
//				"nmea_tll_timeUTC":         NMEA_TLL_TIME_UTC,
//				"nmea_tll_targetStatus":    NMEA_TLL_TARGET_STATUS,
//				"nmea_tll_referenceTarget": NMEA_TLL_REFERENCE_TARGET,
//			},
//		},
//		{
//			name: "NMEA_TTM_*",
//			set: map[string]Token{
//				"nmea_ttm_targetNumber":      NMEA_TTM_TARGET_NUMBER,
//				"nmea_ttm_targetDistance":    NMEA_TTM_TARGET_DISTANCE,
//				"nmea_ttm_bearing":           NMEA_TTM_BEARING,
//				"nmea_ttm_bearingType":       NMEA_TTM_BEARING_TYPE,
//				"nmea_ttm_targetSpeed":       NMEA_TTM_TARGET_SPEED,
//				"nmea_ttm_targetCourse":      NMEA_TTM_TARGET_COURSE,
//				"nmea_ttm_courseType":        NMEA_TTM_COURSE_TYPE,
//				"nmea_ttm_distanceCPA":       NMEA_TTM_DISTANCE_CPA,
//				"nmea_ttm_timeCPA":           NMEA_TTM_TIME_CPA,
//				"nmea_ttm_speedUnits":        NMEA_TTM_SPEED_UNITS,
//				"nmea_ttm_targetName":        NMEA_TTM_TARGET_NAME,
//				"nmea_ttm_targetStatus":      NMEA_TTM_TARGET_STATUS,
//				"nmea_ttm_referenceTarget":   NMEA_TTM_REFERENCE_TARGET,
//				"nmea_ttm_timeUTC":           NMEA_TTM_TIME_UTC,
//				"nmea_ttm_typeOfAcquisition": NMEA_TTM_TYPE_OF_ACQUISITION,
//			},
//		},
//		{
//			name: "NMEA_VBW_*",
//			set: map[string]Token{
//				"nmea_vbw_longitudinalWaterSpeedKnots":         NMEA_VBW_LONGITUDINAL_WATER_SPEED_KNOTS,
//				"nmea_vbw_transverseWaterSpeedKnots":           NMEA_VBW_TRANSVERSE_WATER_SPEED_KNOTS,
//				"nmea_vbw_waterSpeedStatusValid":               NMEA_VBW_WATER_SPEED_STATUS_VALID,
//				"nmea_vbw_waterSpeedStatus":                    NMEA_VBW_WATER_SPEED_STATUS,
//				"nmea_vbw_longitudinalGroundSpeedKnots":        NMEA_VBW_LONGITUDINAL_GROUND_SPEED_KNOTS,
//				"nmea_vbw_transverseGroundSpeedKnots":          NMEA_VBW_TRANSVERSE_GROUND_SPEED_KNOTS,
//				"nmea_vbw_groundSpeedStatusValid":              NMEA_VBW_GROUND_SPEED_STATUS_VALID,
//				"nmea_vbw_groundSpeedStatus":                   NMEA_VBW_GROUND_SPEED_STATUS,
//				"nmea_vbw_sternTraverseWaterSpeedKnots":        NMEA_VBW_STERN_TRAVERSE_WATER_SPEED_KNOTS,
//				"nmea_vbw_sternTraverseWaterSpeedStatusValid":  NMEA_VBW_STERN_TRAVERSE_WATER_SPEED_STATUS_VALID,
//				"nmea_vbw_sternTraverseWaterSpeedStatus":       NMEA_VBW_STERN_TRAVERSE_WATER_SPEED_STATUS,
//				"nmea_vbw_sternTraverseGroundSpeedKnots":       NMEA_VBW_STERN_TRAVERSE_GROUND_SPEED_KNOTS,
//				"nmea_vbw_sternTraverseGroundSpeedStatusValid": NMEA_VBW_STERN_TRAVERSE_GROUND_SPEED_STATUS_VALID,
//				"nmea_vbw_sternTraverseGroundSpeedStatus":      NMEA_VBW_STERN_TRAVERSE_GROUND_SPEED_STATUS,
//			},
//		},
//		{
//			name: "NMEA_VDR_*",
//			set: map[string]Token{
//				"nmea_vdr_setDegreesTrue":         NMEA_VDR_SET_DEGREES_TRUE,
//				"nmea_vdr_setDegreesTrueUnit":     NMEA_VDR_SET_DEGREES_TRUE_UNIT,
//				"nmea_vdr_setDegreesMagnetic":     NMEA_VDR_SET_DEGREES_MAGNETIC,
//				"nmea_vdr_setDegreesMagneticUnit": NMEA_VDR_SET_DEGREES_MAGNETIC_UNIT,
//				"nmea_vdr_driftKnots":             NMEA_VDR_DRIFT_KNOTS,
//				"nmea_vdr_driftUnit":              NMEA_VDR_DRIFT_UNIT,
//			},
//		},
//		{
//			name: "NMEA_VHW_*",
//			set: map[string]Token{
//				"nmea_vhw_trueHeading":            NMEA_VHW_TRUE_HEADING,
//				"nmea_vhw_magneticHeading":        NMEA_VHW_MAGNETIC_HEADING,
//				"nmea_vhw_speedThroughWaterKnots": NMEA_VHW_SPEED_THROUGHT_WATER_KNOTS,
//				"nmea_vhw_speedThroughWaterKPH":   NMEA_VHW_SPEED_THROUGHT_WATER_KPH,
//			},
//		},
//		{
//			name: "NMEA_VLW_*",
//			set: map[string]Token{
//				"nmea_vlw_totalInWater":           NMEA_VLW_TOTAL_IN_WATER,
//				"nmea_vlw_totalInWaterUnit":       NMEA_VLW_TOTAL_IN_WATER_UNIT,
//				"nmea_vlw_sinceResetInWater":      NMEA_VLW_SINCE_RESET_IN_WATER,
//				"nmea_vlw_sinceResetInWaterUnit":  NMEA_VLW_SINCE_RESET_IN_WATER_UNIT,
//				"nmea_vlw_totalOnGround":          NMEA_VLW_TOTAL_ON_GROUND,
//				"nmea_vlw_totalOnGroundUnit":      NMEA_VLW_TOTAL_ON_GROUND_UNIT,
//				"nmea_vlw_sinceResetOnGround":     NMEA_VLW_SINCE_RESET_ON_GROUND,
//				"nmea_vlw_sinceResetOnGroundUnit": NMEA_VLW_SINCE_RESET_ON_GROUND_UNIT,
//			},
//		},
//		{
//			name: "NMEA_VPW_*",
//			set: map[string]Token{
//				"nmea_vpw_speedKnots":     NMEA_VPW_SPEED_KNOTS,
//				"nmea_vpw_speedKnotsUnit": NMEA_VPW_SPEED_KNOTS_UNIT,
//				"nmea_vpw_speedMPS":       NMEA_VPW_SPEED_MPS,
//				"nmea_vpw_speedMPSUnit":   NMEA_VPW_SPEED_MPS_UNIT,
//			},
//		},
//		{
//			name: "NMEA_VTG_*",
//			set: map[string]Token{
//				"nmea_vtg_trueTrack":        NMEA_VTG_TRUE_TRACK,
//				"nmea_vtg_magneticTrack":    NMEA_VTG_MAGNETIC_TRACK,
//				"nmea_vtg_groundSpeedKnots": NMEA_VTG_GROUND_SPEED_KNOTS,
//				"nmea_vtg_groundSpeedKPH":   NMEA_VTG_GROUND_SPEED_KPH,
//				"nmea_vtg_ffaMode":          NMEA_VTG_FFAMODE,
//			},
//		},
//		{
//			name: "NMEA_VWR_*",
//			set: map[string]Token{
//				"nmea_vwr_measuredAngle":        NMEA_VWR_MEASURED_ANGLE,
//				"nmea_vwr_measuredDirectionBow": NMEA_VWR_MEASURED_DIRECTION_BOW,
//				"nmea_vwr_speedKnots":           NMEA_VWR_SPEED_KNOTS,
//				"nmea_vwr_speedKnotsUnit":       NMEA_VWR_SPEED_KNOTS_UNIT,
//				"nmea_vwr_speedMPS":             NMEA_VWR_SPEED_MPS,
//				"nmea_vwr_speedMPSUnit":         NMEA_VWR_SPEED_MPS_UNIT,
//				"nmea_vwr_speedKPH":             NMEA_VWR_SPEED_KPH,
//				"nmea_vwr_speedKPHUnit":         NMEA_VWR_SPEED_KPH_UNIT,
//			},
//		},
//		{
//			name: "NMEA_VWT_*",
//			set: map[string]Token{
//				"nmea_vwt_trueAngle":        NMEA_VWT_TRUE_ANGLE,
//				"nmea_vwt_trueDirectionBow": NMEA_VWT_TRUE_DIRECTION_BOW,
//				"nmea_vwt_speedKnots":       NMEA_VWT_SPEED_KNOTS,
//				"nmea_vwt_speedKnotsUnit":   NMEA_VWT_SPEED_KNOTS_UNIT,
//				"nmea_vwt_speedMPS":         NMEA_VWT_SPEED_MPS,
//				"nmea_vwt_speedMPSUnit":     NMEA_VWT_SPEED_MPS_UNIT,
//				"nmea_vwt_speedKPH":         NMEA_VWT_SPEED_KPH,
//				"nmea_vwt_speedKPHUnit":     NMEA_VWT_SPEED_KPH_UNIT,
//			},
//		},
//		{
//			name: "NMEA_WPL_*",
//			set: map[string]Token{
//				"nmea_wpl_latitude":  NMEA_WPL_LATITUDE,
//				"nmea_wpl_longitude": NMEA_WPL_LONGITUDE,
//				"nmea_wpl_ident":     NMEA_WPL_IDENT,
//			},
//		},
//	}
//
//	toStr := func(set map[string]Token) (str string) {
//		for k := range set {
//			str += k + " "
//		}
//		return
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			str := toStr(tc.set)
//			tokenizer := NewTokenizer(strings.NewReader(str))
//			for i := 0; i < len(tc.set); i++ {
//				haveTok, haveLit := tokenizer.Scan()
//				wantTok, ok := tc.set[haveLit]
//				if !ok {
//					t.Fatalf("token not found %s", haveLit)
//				}
//				if haveTok != wantTok {
//					t.Fatalf("have %d, want %d token id", haveTok, wantTok)
//				}
//			}
//		})
//	}
//}

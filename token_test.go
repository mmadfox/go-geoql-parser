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
				"nmea:bec:longitude":                  NMEA_BEC_LONGTITUDE,
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
				"nmea:bwc:longitude":                 NMEA_BWC_LONGTITUDE,
				"nmea:bwc:bearingTrue":               NMEA_BWC_BEARING_TRUE,
				"nmea:bwc:bearingTrueType":           NMEA_BWC_BEARING_TRUE_TYPE,
				"nmea:bwc:bearingMagnetic":           NMEA_BWC_BEARING_MAGNETIC,
				"nmea:bwc:bearingMagneticType":       NMEA_BWC_BEARING_MAGNETIC_TYPE,
				"nmea:bwc:distanceNauticalMiles":     NMEA_BWC_DISTANCE_NAUTICAL_MILES,
				"nmea:bwc:distanceNauticalMilesUnit": NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT,
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

package geoqlparser

import "strings"

const (
	ILLEGAL Token = iota
	UNUSED
	EOF

	keywordsBegin
	TRIGGER    // trigger
	WHEN       // when
	VARS       // vars
	REPEAT     // repeat
	RESET      // reset
	AFTER      // after
	INTERVAL   // interval
	TIMES      // times
	NOTBETWEEN // not between
	BETWEEN    // between
	keywordsEnd

	INT    // 1
	FLOAT  // 1.1
	STRING // "1"

	ASSIGN    // =
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	COMMA     // ,
	RBRACK    // ]
	LBRACK    // [
	RBRACE    // }
	LBRACE    // {
	QUO       // /
	MUL       // *

	operatorBegin
	GEQ  // >=
	LEQ  // <=
	NEQ  // !=
	GTR  // >
	LSS  // <
	LAND // &&
	LOR  // ||
	AND  // and
	OR   // or
	operatorEnd

	selectorBegin
	TRACKER // tracker - latitude, longitude, altitude
	OBJECT  // object - reference to external objects
	SPEED   // speed

	// NMEA AAM - Waypoint Arrival Alarm
	NMEA_AAM_STATUS_ARRIVAL_CIRCLE_ENTERED // nmea:aam:statusArrivalCircleEntered
	NMEA_AAM_STATUS_PERPENDICULAR_PASSED   // nmea:aam:statusPerpendicularPassed
	NMEA_AAM_ARRIVAL_CIRCLE_RADIUS         // nmea:aam:arrivalCircleRadius
	NMEA_AAM_ARRIVAL_CIRCLE_UNIT           // nmea:aam:arrivalCircleRadiusUnit
	NMEA_AAM_DESTINATION_WAYPOINT_ID       // nmea:aam:destinationWaypointID

	// NMEA ALA - System Faults and Alarms
	NMEA_ALA_TIME                 // nmea:ala:time
	NMEA_ALA_SYSTEM_INDICATOR     // nmea:ala:systemIndicator
	NMEA_ALA_SUB_SYSTEM_INDICATOR // nmea:ala:subSystemIndicator
	NMEA_ALA_INSTANCE_NUMBER      // nmea:ala:instanceNumber
	NMEA_ALA_TYPE                 // nmea:ala:type
	NMEA_ALA_CONDITION            // nmea:ala:condition
	NMEA_ALA_ALARM_ACK_STATE      // nmea:ala:alarmAckState
	NMEA_ALA_MESSAGE              // nmea:ala:message

	// NMEA APB - Autopilot Sentence "B" for heading/tracking
	NMEA_APB_STATUS_GENERAL_WARNING        // nmea:apb:statusGeneralWarning
	NMEA_APB_STATUS_LOCK_WARNING           // nmea:apb:statusLockWarning
	NMEA_APB_CROSS_TRACK_ERROR_MAGNITUDE   // nmea:apb:crossTrackErrorMagnitude
	NMEA_APB_DIRECTION_TO_STEER            // nmea:apb:directionToSteer
	NMEA_APB_CROSS_TRACK_UNITS             // nmea:apb:crossTrackUnits
	NMEA_APB_STATUS_ARRIVAL_CIRCLE_ENTERED // nmea:apb:statusArrivalCircleEntered
	NMEA_APB_STATUS_PERPENDICULAR_PASSED   // nmea:apb:statusPerpendicularPassed
	NMEA_APB_BEARING_ORIGIN_TO_DEST        // nmea:apb:bearingOriginToDest
	NMEA_APB_BEARING_ORIGIN_DEST_TYPE      // nmea:apb:bearingOriginToDestType
	NMEA_APB_DESTINATION_WAYPOINT_ID       // nmea:apb:destinationWaypointID
	NMEA_APB_BEARING_PRESENT_TO_DEST       // nmea:apb:bearingPresentToDest
	NMEA_APB_BEARING_PRESENT_TO_DEST_TYPE  // nmea:apb:bearingPresentToDestType
	NMEA_APB_HEADING                       // nmea:apb:heading
	NMEA_APB_HEADING_TYPE                  // nmea:apb:headingType
	NMEA_APB_FFA_MODE                      // nmea:apb:ffaMode

	// NMEA BEC - bearing and distance to waypoint (dead reckoning)
	NMEA_BEC_TIME                          // nmea:bec:time
	NMEA_BEC_LATITUDE                      // nmea:bec:latitude
	NMEA_BEC_LONGTITUDE                    // nmea:bec:longitude
	NMEA_BEC_BEARING_TRUE                  // nmea:bec:bearingTrue
	NMEA_BEC_BEARING_TRUE_VALID            // nmea:bec:bearingTrueValid
	NMEA_BEC_BEARING_MAGNETIC              // nmea:bec:bearingMagnetic
	NMEA_BEC_BEARING_MAGNETIC_VALID        // nmea:bec:bearingMagneticValid
	NMEA_BEC_DISTANCE_NAUTICAL_MILES       // nmea:bec:distanceNauticalMiles
	NMEA_BEC_DISTANCE_NAUTICAL_MILES_VALID // nmea:bec:distanceNauticalMilesValid
	NMEA_BEC_DESTINATION_WAYPOINT_ID       // nmea:bec:destinationWaypointID

	// NMEA BOD - bearing waypoint to waypoint (origin to destination).
	NMEA_BOD_BEARING_TRUE            // nmea:bod:bearingTrue
	NMEA_BOD_BEARING_TRUE_TYPE       // nmea:bod:bearingTrueType
	NMEA_BOD_BEARING_MAGNETIC        // nmea:bod:bearingMagnetic
	NMEA_BOD_BEARING_MAGNETIC_TYPE   // nmea:bod:bearingMagneticType
	NMEA_BOD_DESTINATION_WAYPOINT_ID // nmea:bod:destinationWaypointID
	NMEA_BOD_ORIGIN_WAYPOINT_ID      // nmea:bod:originWaypointID

	// NMEA BWC - bearing and distance to waypoint, great circle
	NMEA_BWC_TIME                         // nmea:bwc:time
	NMEA_BWC_LATITUDE                     // nmea:bwc:latitude
	NMEA_BWC_LONGTITUDE                   // nmea:bwc:longitude
	NMEA_BWC_BEARING_TRUE                 // nmea:bwc:bearingTrue
	NMEA_BWC_BEARING_TRUE_TYPE            // nmea:bwc:bearingTrueType
	NMEA_BWC_BEARING_MAGNETIC             // nmea:bwc:bearingMagnetic
	NMEA_BWC_BEARING_MAGNETIC_TYPE        // nmea:bwc:bearingMagneticType
	NMEA_BWC_DISTANCE_NAUTICAL_MILES      // nmea:bwc:distanceNauticalMiles
	NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT // nmea:bwc:distanceNauticalMilesUnit
	selectorEnd
)

var keywords = map[string]Token{
	"trigger":  TRIGGER,
	"vars":     VARS,
	"when":     WHEN,
	"repeat":   REPEAT,
	"reset":    RESET,
	"after":    AFTER,
	"interval": INTERVAL,
	"times":    TIMES,
	"between":  BETWEEN,

	"not between": NOTBETWEEN,

	"=":   ASSIGN,
	";":   SEMICOLON,
	"(":   LPAREN,
	")":   RPAREN,
	",":   COMMA,
	">=":  GEQ,
	"<=":  LEQ,
	"!=":  NEQ,
	">":   GTR,
	"<":   LSS,
	"&&":  LAND,
	"||":  LOR,
	"or":  OR,
	"and": AND,
	"[":   LBRACK,
	"]":   RBRACK,
	"{":   LBRACE,
	"}":   RBRACE,
	"/":   QUO,
	"*":   MUL,
	":":   COLON,

	"tracker": TRACKER,
	"object":  OBJECT,
	"speed":   SPEED,

	"nmea:aam:statusArrivalCircleEntered": NMEA_AAM_STATUS_ARRIVAL_CIRCLE_ENTERED,
	"nmea:aam:statusPerpendicularPassed":  NMEA_AAM_STATUS_PERPENDICULAR_PASSED,
	"nmea:aam:arrivalCircleRadius":        NMEA_AAM_ARRIVAL_CIRCLE_RADIUS,
	"nmea:aam:arrivalCircleRadiusUnit":    NMEA_AAM_ARRIVAL_CIRCLE_UNIT,
	"nmea:aam:destinationWaypointID":      NMEA_AAM_DESTINATION_WAYPOINT_ID,

	"nmea:ala:time":               NMEA_ALA_TIME,
	"nmea:ala:systemIndicator":    NMEA_ALA_SYSTEM_INDICATOR,
	"nmea:ala:subSystemIndicator": NMEA_ALA_SUB_SYSTEM_INDICATOR,
	"nmea:ala:instanceNumber":     NMEA_ALA_INSTANCE_NUMBER,
	"nmea:ala:type":               NMEA_ALA_TYPE,
	"nmea:ala:condition":          NMEA_ALA_CONDITION,
	"nmea:ala:alarmAckState":      NMEA_ALA_ALARM_ACK_STATE,
	"nmea:ala:message":            NMEA_ALA_MESSAGE,

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

	"nmea:bod:bearingTrue":           NMEA_BOD_BEARING_TRUE,
	"nmea:bod:bearingTrueType":       NMEA_BOD_BEARING_TRUE_TYPE,
	"nmea:bod:bearingMagnetic":       NMEA_BOD_BEARING_MAGNETIC,
	"nmea:bod:bearingMagneticType":   NMEA_BOD_BEARING_MAGNETIC_TYPE,
	"nmea:bod:destinationWaypointID": NMEA_BOD_DESTINATION_WAYPOINT_ID,
	"nmea:bod:originWaypointID":      NMEA_BOD_ORIGIN_WAYPOINT_ID,

	"nmea:bwc:time":                      NMEA_BWC_TIME,
	"nmea:bwc:latitude":                  NMEA_BWC_LATITUDE,
	"nmea:bwc:longitude":                 NMEA_BWC_LONGTITUDE,
	"nmea:bwc:bearingTrue":               NMEA_BWC_BEARING_TRUE,
	"nmea:bwc:bearingTrueType":           NMEA_BWC_BEARING_TRUE_TYPE,
	"nmea:bwc:bearingMagnetic":           NMEA_BWC_BEARING_MAGNETIC,
	"nmea:bwc:bearingMagneticType":       NMEA_BWC_BEARING_MAGNETIC_TYPE,
	"nmea:bwc:distanceNauticalMiles":     NMEA_BWC_DISTANCE_NAUTICAL_MILES,
	"nmea:bwc:distanceNauticalMilesUnit": NMEA_BWC_DISTANCE_NAUTICAL_MILES_UNIT,
}

var keywordStrings = map[Token]string{}

func init() {
	for str, id := range keywords {
		keywordStrings[id] = str
	}
}

func KeywordString(id Token) string {
	str, ok := keywordStrings[id]
	if !ok {
		return type2str(id)
	}
	if id >= keywordsBegin && id <= keywordsEnd {
		str = strings.ToUpper(str)
	}
	return str
}

func isSelector(tok Token) bool {
	return tok >= selectorBegin && tok <= selectorEnd
}

func isOperator(tok Token) bool {
	return tok >= operatorBegin && tok <= operatorEnd
}

func type2str(id Token) (str string) {
	switch id {
	case UNUSED:
		str = "UNUSED"
	case FLOAT:
		str = "FLOATVAL"
	case INT:
		str = "INTVAL"
	case STRING:
		str = "STRINGVAL"
	}
	return
}

package knowledgecommands

import (
	"regexp"
	"sort"
	"strings"
)

func activeInterviewEntries(session map[string]any, layer string) []any {
	entries := []any{}
	for _, raw := range interviewLayerEntries(session, layer) {
		if stringValue(mapValue(raw)["status"]) != "stale" {
			entries = append(entries, raw)
		}
	}
	return entries
}

func deriveStakeholderTones(session map[string]any) []any {
	// The template uses all historical entries here, with formal > terse > informal.
	strengths := map[string]int{}
	for layer := range mapValue(session["layers"]) {
		strength := 0
		switch layer {
		case "dependencies":
			strength = 3
		case "friction", "recurring_decisions":
			strength = 2
		case "institutional_knowledge":
			strength = 1
		}
		// No tone signal is declared for rhythms alone.
		if strength == 0 {
			continue
		}
		for _, raw := range interviewLayerEntries(session, layer) {
			for _, value := range sliceValue(mapValue(raw)["stakeholders"]) {
				name := stringValue(value)
				if strength > strengths[name] {
					strengths[name] = strength
				}
			}
		}
	}
	names := make([]string, 0, len(strengths))
	for name := range strengths {
		names = append(names, name)
	}
	sort.Strings(names)
	result := []any{}
	for _, name := range names {
		result = append(result, map[string]any{"stakeholder": name, "tone": []string{"", "informal", "terse", "formal"}[strengths[name]]})
	}
	return result
}

var interviewWeekdayPattern = regexp.MustCompile(`(?i)\b(mon(?:day)?|tue(?:sday)?|wed(?:nesday)?|thu(?:rsday)?|fri(?:day)?|sat(?:urday)?|sun(?:day)?)\b`)
var interviewCadencePattern = regexp.MustCompile(`(?i)\b(daily|weekly|monthly|every day|each day|every week|each week|every month|each month)\b`)
var interviewClockPattern = regexp.MustCompile(`\b([0-9]{1,2}:[0-9]{2})\b`)
var interviewValidClock = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):[0-5][0-9]$`)

func parseInterviewCadence(text string) (cadence, day, clock string, ok bool) {
	day = "rolling"
	if weekday := interviewWeekdayPattern.FindString(text); weekday != "" {
		day = strings.ToUpper(weekday[:1]) + strings.ToLower(weekday[1:3])
	}
	switch strings.ToLower(interviewCadencePattern.FindString(text)) {
	case "daily", "every day", "each day":
		cadence = "daily"
	case "weekly", "every week", "each week":
		cadence = "weekly"
	case "monthly", "every month", "each month":
		cadence = "monthly"
	default:
		if day != "rolling" {
			cadence = "weekly"
		}
	}
	if cadence == "" {
		return "", "", "", false
	}
	clock = interviewClockPattern.FindString(text)
	if len(clock) == 4 {
		clock = "0" + clock
	}
	if clock != "" && !interviewValidClock.MatchString(clock) {
		return "", "", "", false
	}
	return cadence, day, clock, true
}

func deriveStandingSlots(session map[string]any) []any {
	slots := []any{}
	for _, raw := range activeInterviewEntries(session, "dependencies") {
		item := mapValue(raw)
		details := mapValue(item["details"])
		cadence, day, clock, ok := parseInterviewCadence(stringValue(details["needed_by"]))
		if !ok {
			continue
		}
		slots = append(slots, map[string]any{"label": details["deliverable"], "cadence": cadence, "day": day, "time": clock, "counterparty": details["dependency_owner"], "source_entries": []any{"dependencies." + stringValue(item["entry_id"])}})
	}
	return slots
}

// suppliedAvoidWindow selects a window from the same entry. Text naming a label or
// clock selects that window; a single supplied window is unambiguous. No matching
// window among several is insufficient evidence to invent an avoidance schedule.
func suppliedAvoidWindow(details map[string]any, description string) map[string]any {
	windows := sliceValue(details["time_windows"])
	var selected map[string]any
	for _, raw := range windows {
		window := mapValue(raw)
		start, end := stringValue(window["start"]), stringValue(window["end"])
		if !interviewValidClock.MatchString(start) || !interviewValidClock.MatchString(end) || len(sliceValue(window["days"])) == 0 {
			continue
		}
		label := strings.ToLower(stringValue(window["label"]))
		matches := len(windows) == 1 || (label != "" && strings.Contains(strings.ToLower(description), label)) || strings.Contains(description, start)
		if !matches {
			continue
		}
		if selected != nil {
			return nil
		}
		selected = window
	}
	return selected
}

var interviewInterruptionPattern = regexp.MustCompile(`(?i)\b(interrupt\w*|disrupt\w*|delay\w*|wait\w*|sync\w*|handoff\w*|loss|lost|lose|eats?|takes?|minutes?|hours?)\b`)

func deriveAvoidWindows(session map[string]any) []any {
	result := []any{}
	for _, layer := range []string{"operating_rhythms", "friction"} {
		for _, raw := range activeInterviewEntries(session, layer) {
			item := mapValue(raw)
			details := mapValue(item["details"])
			description := stringValue(details["non_calendar_reality"])
			reason := description
			label := description
			if layer == "friction" {
				if stringValue(details["priority"]) != "high" {
					continue
				}
				description = displayTemplateValue(details["systems_involved"]) + " " + stringValue(details["frequency"]) + " " + stringValue(details["time_cost"])
				reason, label = stringValue(item["summary"]), stringValue(item["title"])
			}
			if _, _, _, recurring := parseInterviewCadence(description); !recurring || !interviewInterruptionPattern.MatchString(description) {
				continue
			}
			window := suppliedAvoidWindow(details, description)
			if window == nil {
				continue
			}
			if layer == "operating_rhythms" {
				label = stringValue(window["label"])
			}
			result = append(result, map[string]any{"label": label, "days": window["days"], "start": window["start"], "end": window["end"], "reason": reason, "source_entries": []any{layer + "." + stringValue(item["entry_id"])}})
		}
	}
	return result
}

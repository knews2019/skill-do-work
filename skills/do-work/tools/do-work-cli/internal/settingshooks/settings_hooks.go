// Package settingshooks composes do-work's hook fragment into a consumer's
// .claude/settings.json without reordering or rewriting anything the consumer owns.
//
// It replaces both incumbent reconcilers — a jq program and an embedded Python composer —
// with one implementation. The whole package exists because Go's encoding/json marshals a
// map with sorted keys: round-tripping a consumer's settings through map[string]any would
// silently reorder unrelated user state on every install. Every value therefore travels as
// an orderedObject that remembers the key order it was decoded in.
package settingshooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// retiredPipelineGuardCommand identifies the Stop hook this suite no longer ships. Only hook
// objects whose string command CONTAINS it are removed, and only they.
const retiredPipelineGuardCommand = ".claude/skills/do-work/hooks/pipeline-guard.sh"

// ComposeSettings removes the retired pipeline guard, then appends every fragment hook entry
// that is not already present, and re-encodes with two-space indent and a trailing newline.
func ComposeSettings(settingsData, fragmentData []byte) ([]byte, error) {
	settings, err := decodeOrderedJSON(settingsData)
	if err != nil {
		return nil, fmt.Errorf("settings are not valid JSON: %w", err)
	}
	fragment, err := decodeOrderedJSON(fragmentData)
	if err != nil {
		return nil, fmt.Errorf("core hook fragment is not valid JSON: %w", err)
	}

	settingsObject, isObject := settings.(*orderedObject)
	if !isObject {
		return nil, errors.New("settings root must be an object")
	}
	hooksObject, err := hooksObjectOf(settingsObject)
	if err != nil {
		return nil, err
	}
	if err := removeRetiredPipelineGuard(hooksObject); err != nil {
		return nil, err
	}
	if err := appendFragmentEvents(hooksObject, fragment); err != nil {
		return nil, err
	}
	return encodeOrderedJSON(settingsObject)
}

// hooksObjectOf returns the settings' hooks object, creating an empty one when the consumer
// has no hooks at all. A new key lands at the end, which is where both incumbent reconcilers
// put it.
func hooksObjectOf(settings *orderedObject) (*orderedObject, error) {
	existing, present := settings.value("hooks")
	if !present || existing == nil {
		created := newOrderedObject()
		settings.set("hooks", created)
		return created, nil
	}
	hooks, isObject := existing.(*orderedObject)
	if !isObject {
		return nil, errors.New("settings hooks must be an object")
	}
	return hooks, nil
}

func removeRetiredPipelineGuard(hooks *orderedObject) error {
	stopEntries, present := hooks.value("Stop")
	if !present || stopEntries == nil {
		return nil
	}
	wrappers, isArray := stopEntries.([]any)
	if !isArray {
		return errors.New("settings Stop hook event must be an array")
	}

	retained := make([]any, 0, len(wrappers))
	for _, wrapper := range wrappers {
		wrapperObject, isObject := wrapper.(*orderedObject)
		if !isObject {
			retained = append(retained, wrapper)
			continue
		}
		wrapperHooks, hasHookArray := wrapperObject.value("hooks")
		hookEntries, isHookArray := wrapperHooks.([]any)
		if !hasHookArray || !isHookArray {
			retained = append(retained, wrapper)
			continue
		}
		keptHooks := make([]any, 0, len(hookEntries))
		for _, hookEntry := range hookEntries {
			if !isRetiredPipelineGuard(hookEntry) {
				keptHooks = append(keptHooks, hookEntry)
			}
		}
		// A wrapper that lost nothing is preserved exactly as it was, empty or not; only a
		// wrapper emptied BY the removal disappears.
		if len(keptHooks) == len(hookEntries) {
			retained = append(retained, wrapper)
			continue
		}
		if len(keptHooks) == 0 {
			continue
		}
		wrapperObject.set("hooks", keptHooks)
		retained = append(retained, wrapperObject)
	}
	if len(retained) == 0 {
		hooks.delete("Stop")
		return nil
	}
	hooks.set("Stop", retained)
	return nil
}

func isRetiredPipelineGuard(hookEntry any) bool {
	hookObject, isObject := hookEntry.(*orderedObject)
	if !isObject {
		return false
	}
	command, present := hookObject.value("command")
	if !present {
		return false
	}
	commandText, isString := command.(string)
	return isString && strings.Contains(commandText, retiredPipelineGuardCommand)
}

// appendFragmentEvents walks the fragment's events in sorted order, matching the jq program
// that was the preferred incumbent branch, so the composed output does not depend on how the
// shipped fragment happens to be written.
func appendFragmentEvents(hooks *orderedObject, fragment any) error {
	fragmentObject, isObject := fragment.(*orderedObject)
	if !isObject {
		return errors.New("core hook fragment root must be an object")
	}
	fragmentHooks, present := fragmentObject.value("hooks")
	if !present || fragmentHooks == nil {
		return nil
	}
	fragmentHooksObject, isHooksObject := fragmentHooks.(*orderedObject)
	if !isHooksObject {
		return errors.New("core hook fragment hooks must be an object")
	}

	eventNames := append([]string(nil), fragmentHooksObject.keys...)
	sort.Strings(eventNames)
	for _, eventName := range eventNames {
		fragmentEntries, isArray := mustValue(fragmentHooksObject, eventName).([]any)
		if !isArray {
			return errors.New("settings hook event must be an array")
		}
		existing, hasEvent := hooks.value(eventName)
		if !hasEvent || existing == nil {
			hooks.set(eventName, append([]any(nil), fragmentEntries...))
			continue
		}
		installedEntries, isInstalledArray := existing.([]any)
		if !isInstalledArray {
			return errors.New("settings hook event must be an array")
		}
		for _, fragmentEntry := range fragmentEntries {
			if !containsJSONValue(installedEntries, fragmentEntry) {
				installedEntries = append(installedEntries, fragmentEntry)
			}
		}
		hooks.set(eventName, installedEntries)
	}
	return nil
}

func mustValue(object *orderedObject, key string) any {
	value, _ := object.value(key)
	return value
}

func containsJSONValue(entries []any, candidate any) bool {
	for _, entry := range entries {
		if jsonValuesEqual(entry, candidate) {
			return true
		}
	}
	return false
}

// jsonValuesEqual compares by JSON value, not by key order: both incumbent reconcilers treat
// two objects carrying the same pairs as the same entry.
func jsonValuesEqual(first, second any) bool {
	switch firstTyped := first.(type) {
	case *orderedObject:
		secondTyped, sameKind := second.(*orderedObject)
		if !sameKind || len(firstTyped.keys) != len(secondTyped.keys) {
			return false
		}
		for _, key := range firstTyped.keys {
			secondValue, present := secondTyped.value(key)
			if !present || !jsonValuesEqual(mustValue(firstTyped, key), secondValue) {
				return false
			}
		}
		return true
	case []any:
		secondTyped, sameKind := second.([]any)
		if !sameKind || len(firstTyped) != len(secondTyped) {
			return false
		}
		for index := range firstTyped {
			if !jsonValuesEqual(firstTyped[index], secondTyped[index]) {
				return false
			}
		}
		return true
	case json.Number:
		secondTyped, sameKind := second.(json.Number)
		if !sameKind {
			return false
		}
		if firstTyped.String() == secondTyped.String() {
			return true
		}
		firstFloat, firstErr := firstTyped.Float64()
		secondFloat, secondErr := secondTyped.Float64()
		return firstErr == nil && secondErr == nil && firstFloat == secondFloat
	default:
		return first == second
	}
}

// orderedObject is a JSON object that remembers its key order. Every decoded object in this
// package is one, which is what keeps a consumer's settings.json byte-stable.
type orderedObject struct {
	keys   []string
	values map[string]any
}

func newOrderedObject() *orderedObject {
	return &orderedObject{values: map[string]any{}}
}

func (object *orderedObject) value(key string) (any, bool) {
	value, present := object.values[key]
	return value, present
}

func (object *orderedObject) set(key string, value any) {
	if _, present := object.values[key]; !present {
		object.keys = append(object.keys, key)
	}
	object.values[key] = value
}

func (object *orderedObject) delete(key string) {
	if _, present := object.values[key]; !present {
		return
	}
	delete(object.values, key)
	for index, existing := range object.keys {
		if existing == key {
			object.keys = append(object.keys[:index], object.keys[index+1:]...)
			return
		}
	}
}

func decodeOrderedJSON(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := decodeOrderedValue(decoder)
	if err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, errors.New("trailing content after the JSON document")
	}
	return value, nil
}

func decodeOrderedValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return token, nil
	}
	switch delimiter {
	case '{':
		object := newOrderedObject()
		for decoder.More() {
			keyToken, keyErr := decoder.Token()
			if keyErr != nil {
				return nil, keyErr
			}
			key, isString := keyToken.(string)
			if !isString {
				return nil, errors.New("object key is not a string")
			}
			value, valueErr := decodeOrderedValue(decoder)
			if valueErr != nil {
				return nil, valueErr
			}
			object.set(key, value)
		}
		if _, err := decoder.Token(); err != nil {
			return nil, err
		}
		return object, nil
	case '[':
		items := []any{}
		for decoder.More() {
			item, itemErr := decodeOrderedValue(decoder)
			if itemErr != nil {
				return nil, itemErr
			}
			items = append(items, item)
		}
		if _, err := decoder.Token(); err != nil {
			return nil, err
		}
		return items, nil
	default:
		return nil, fmt.Errorf("unexpected JSON delimiter %q", delimiter)
	}
}

// encodeOrderedJSON writes two-space-indented JSON with a trailing newline, with HTML
// escaping off so a hook command's `<`, `>` and `&` survive as themselves.
func encodeOrderedJSON(value any) ([]byte, error) {
	var output bytes.Buffer
	if err := encodeOrderedValue(&output, value, 0); err != nil {
		return nil, err
	}
	output.WriteByte('\n')
	return output.Bytes(), nil
}

func encodeOrderedValue(output *bytes.Buffer, value any, indentLevel int) error {
	switch typed := value.(type) {
	case *orderedObject:
		if len(typed.keys) == 0 {
			output.WriteString("{}")
			return nil
		}
		output.WriteString("{\n")
		for index, key := range typed.keys {
			writeIndent(output, indentLevel+1)
			if err := encodeJSONString(output, key); err != nil {
				return err
			}
			output.WriteString(": ")
			if err := encodeOrderedValue(output, mustValue(typed, key), indentLevel+1); err != nil {
				return err
			}
			if index < len(typed.keys)-1 {
				output.WriteByte(',')
			}
			output.WriteByte('\n')
		}
		writeIndent(output, indentLevel)
		output.WriteByte('}')
		return nil
	case []any:
		if len(typed) == 0 {
			output.WriteString("[]")
			return nil
		}
		output.WriteString("[\n")
		for index, item := range typed {
			writeIndent(output, indentLevel+1)
			if err := encodeOrderedValue(output, item, indentLevel+1); err != nil {
				return err
			}
			if index < len(typed)-1 {
				output.WriteByte(',')
			}
			output.WriteByte('\n')
		}
		writeIndent(output, indentLevel)
		output.WriteByte(']')
		return nil
	case string:
		return encodeJSONString(output, typed)
	case json.Number:
		output.WriteString(typed.String())
		return nil
	case bool:
		output.WriteString(strconv.FormatBool(typed))
		return nil
	case nil:
		output.WriteString("null")
		return nil
	default:
		return fmt.Errorf("unsupported JSON value of type %T", value)
	}
}

func writeIndent(output *bytes.Buffer, indentLevel int) {
	for count := 0; count < indentLevel; count++ {
		output.WriteString("  ")
	}
}

// encodeJSONString escapes only what JSON requires. HTML escaping stays off so `<`, `>` and
// `&` survive, matching the jq branch this replaces.
func encodeJSONString(output *bytes.Buffer, value string) error {
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return err
	}
	output.Write(bytes.TrimRight(encoded.Bytes(), "\n"))
	return nil
}

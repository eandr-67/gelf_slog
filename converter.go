package gelf_slog

import (
	"encoding/json"
	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"
var ErrorKeys = []string{"error", "err"}

type Converter func(
	addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string,
	record *slog.Record) (extra map[string]any)

func DefaultConverter(
	addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string,
	record *slog.Record) (extra map[string]any) {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	if addSource {
		attrs = append(attrs, slogcommon.Source(SourceKey, record))
	}
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	// handler formatter
	extra = AttrsToMap(attrs...)

	for _, errorKey := range ErrorKeys {
		if v, ok := extra[errorKey]; ok {
			if err, ok := v.(error); ok {
				extra[errorKey] = slogcommon.FormatError(err)
				break
			}
		}
	}

	return extra
}

func AttrsToMap(attrs ...slog.Attr) map[string]any {
	output := map[string]any{}

	attrsByKey := groupValuesByKey(attrs)
	for k, values := range attrsByKey {
		v := mergeAttrValues(values...)
		if v.Kind() == slog.KindGroup {
			output[k] = AttrsToMap(v.Group()...)
		} else if v.Kind() == slog.KindAny {
			if tmp, err := json.Marshal(v.Any()); err != nil {
				output[k] = "!BAD_SERIALIZE! " + err.Error()
			} else {
				output[k] = string(tmp)
			}
		} else {
			output[k] = v.Any()
		}
	}

	return output
}

func groupValuesByKey(attrs []slog.Attr) map[string][]slog.Value {
	result := map[string][]slog.Value{}

	for _, item := range attrs {
		key := item.Key
		result[key] = append(result[key], item.Value)
	}

	return result
}

func mergeAttrValues(values ...slog.Value) slog.Value {
	v := values[0]

	for i := 1; i < len(values); i++ {
		if v.Kind() != slog.KindGroup || values[i].Kind() != slog.KindGroup {
			v = values[i]
			continue
		}

		v = slog.GroupValue(append(v.Group(), values[i].Group()...)...)
	}

	return v
}

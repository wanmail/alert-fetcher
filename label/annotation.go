package label

import (
	"bytes"
	"html/template"
	"log/slog"
)

func BuildAnnotations(labels map[string]string, relabels map[string]string) map[string]string {
	annotations := make(map[string]string)

	for k, tpl := range relabels {
		t := template.New(k)
		parse, err := t.Parse(tpl)
		if err != nil {
			slog.Error("invalid template", "name", k, "template", tpl, "error", err)
			continue
		}
		bf := bytes.NewBufferString("")
		err = parse.Execute(bf, labels)
		if err != nil {
			slog.Error("failed to execute template", "name", k, "template", tpl, "error", err)
			continue
		} else {
			annotations[k] = bf.String()
		}
	}

	return annotations
}

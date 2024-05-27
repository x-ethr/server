package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/x-ethr/server/logging"
)

func Test(t *testing.T) {
	logging.Verbose(true)

	ctx := context.Background()

	t.Run("New", func(t *testing.T) {
		t.Run("Custom-Log-Levels", func(t *testing.T) {
			t.Run("Trace", func(t *testing.T) {
				const level = slog.Level(-8)
				t.Run("Defaults-Enabled", func(t *testing.T) {
					logging.Level(logging.Trace)

					var w bytes.Buffer

					handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })

					instance := slog.New(handler)
					slog.SetDefault(instance)

					slog.Log(ctx, level, "Test Message")
					if w.Len() == 0 {
						t.Fatalf("Invalid Log Message - Output Should be Enabled: %s", t.Name())
					}

					t.Logf("Success: %s, %s", t.Name(), w.String())
				})

				t.Run("Defaults-Disabled", func(t *testing.T) {
					logging.Level(logging.Info)

					var w bytes.Buffer

					handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })
					instance := slog.New(handler)
					slog.SetDefault(instance)

					slog.Log(ctx, level, "Test Message")
					if w.Len() != 0 {
						t.Fatalf("Invalid Log Message - Output Should be Disabled: %s", t.Name())
					}

					t.Logf("Success: %s, %s", t.Name(), w.String())
				})
			})
		})
		t.Run("Message", func(t *testing.T) {
			t.Run("No-Group", func(t *testing.T) {
				const level = slog.Level(0)

				logging.Level(logging.Info)

				var w bytes.Buffer

				handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })
				instance := slog.New(handler)
				slog.SetDefault(instance)

				slog.Log(ctx, level, "Test Message")

				comparator := fmt.Sprintf("%s\n", "{}")
				if !(strings.HasSuffix(w.String(), comparator)) {
					t.Errorf("Log Output Should Contain Suffice ({}), Received: %s", w.String())
				}

				t.Logf("Success: %s, %s", t.Name(), w.String())
			})

			t.Run("1-Nested-Group", func(t *testing.T) {
				const level = slog.Level(0)

				logging.Level(logging.Info)

				var w bytes.Buffer

				handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })
				instance := slog.New(handler)
				slog.SetDefault(instance)

				slog.Log(ctx, level, "Test Message", slog.Group("group", slog.String("key", "value")))

				partials := strings.SplitN(w.String(), fmt.Sprintf("%s ", "Test Message"), 2)
				if len(partials) != 2 {
					t.Errorf("Log Output Expected Split of Length(2), Received: %d", len(partials))
				}

				content := []byte(partials[1])

				var mapping map[string]map[string]string
				if e := json.Unmarshal(content, &mapping); e != nil {
					t.Fatalf("Error While Unmarshalling: %v", e)
				}

				if mapping["group"]["key"] != "value" {
					t.Errorf("Log Output Should Contain JSON (group.key = value), Received: %s", w.String())
				} else {
					t.Logf("Success: %s, %s", t.Name(), w.String())
				}
			})

			t.Run("2-Nested-Groups", func(t *testing.T) {
				const level = slog.Level(0)

				logging.Level(logging.Info)

				var w bytes.Buffer

				handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })
				instance := slog.New(handler)
				slog.SetDefault(instance)

				slog.Log(ctx, level, "Test Message", slog.Group("group-1", slog.Group("group-2", slog.String("key", "value"))))

				partials := strings.SplitN(w.String(), fmt.Sprintf("%s ", "Test Message"), 2)
				if len(partials) != 2 {
					t.Errorf("Log Output Expected Split of Length(2), Received: %d", len(partials))
				}

				content := []byte(partials[1])

				var mapping map[string]map[string]map[string]string
				if e := json.Unmarshal(content, &mapping); e != nil {
					t.Fatalf("Fatal Error While Unmarshalling: %v", e)
				}

				if mapping["group-1"]["group-2"]["key"] != "value" {
					t.Errorf("Log Output Should Contain JSON (group-1.group-2.key = value), Received: %s", w.String())
				} else {
					t.Logf("Success: %s, %s", t.Name(), w.String())
				}
			})

			t.Run("3-Nested-Groups", func(t *testing.T) {
				const level = slog.Level(0)

				logging.Level(logging.Info)

				var w bytes.Buffer

				handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })
				instance := slog.New(handler)
				slog.SetDefault(instance)

				slog.Log(ctx, level, "Test Message", slog.Group("group-1", slog.Group("group-2", slog.Group("group-3", slog.String("key", "value")))))

				partials := strings.SplitN(w.String(), fmt.Sprintf("%s ", "Test Message"), 2)
				if len(partials) != 2 {
					t.Errorf("Log Output Expected Split of Length(2), Received: %d", len(partials))
				}

				content := []byte(partials[1])

				var mapping map[string]map[string]map[string]map[string]string
				if e := json.Unmarshal(content, &mapping); e != nil {
					t.Fatalf("Fatal Error While Unmarshalling: %v", e)
				}

				if mapping["group-1"]["group-2"]["group-3"]["key"] != "value" {
					t.Errorf("Log Output Should Contain JSON (group-1.group-2.group-3.key = value), Received: %s", w.String())
				} else {
					t.Logf("Success: %s, %s", t.Name(), w.String())
				}
			})

			t.Run("4-Nested-Groups", func(t *testing.T) {
				const level = slog.Level(0)

				logging.Level(logging.Info)

				var w bytes.Buffer

				handler := logging.Logger(func(o *logging.Options) { o.Writer = &w })
				instance := slog.New(handler)
				slog.SetDefault(instance)

				slog.Log(ctx, level, "Test Message", slog.Group("group-1", slog.Group("group-2", slog.Group("group-3", slog.Group("group-4", slog.String("key", "value"))))))

				partials := strings.SplitN(w.String(), fmt.Sprintf("%s ", "Test Message"), 2)
				if len(partials) != 2 {
					t.Errorf("Log Output Expected Split of Length(2), Received: %d", len(partials))
				}

				content := []byte(partials[1])

				var mapping map[string]map[string]map[string]map[string]string
				if e := json.Unmarshal(content, &mapping); e != nil {
					t.Fatalf("Fatal Error While Unmarshalling: %v", e)
				}

				if mapping["group-1"]["group-2"]["group-3"]["group-4"] != "[key=value]" {
					t.Errorf("Log Output Should Contain JSON (group-1.group-2.group-3.key = value), Received: %s", w.String())
				} else {
					t.Logf("Success: %s, %s", t.Name(), w.String())
				}
			})
		})
	})
}

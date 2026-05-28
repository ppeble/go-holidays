package holidays_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	holidays "github.com/ppeble/go-holidays/pkg"
	_ "github.com/ppeble/go-holidays/pkg/definitions"
)

func writeYAML(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

func TestLoadCustom_BasicRuleVisible(t *testing.T) {
	path := writeYAML(t, "myteam.yaml", `
months:
  3:
  - name: Team Anniversary
    regions: [myteam]
    mday: 15
`)
	if err := holidays.LoadCustom(path); err != nil {
		t.Fatalf("LoadCustom: %v", err)
	}
	defer holidays.UnloadCustom(path)

	d := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(d, holidays.Options{Regions: []string{"myteam"}})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if len(hs) != 1 || hs[0].Name != "Team Anniversary" {
		t.Fatalf("expected Team Anniversary; got %+v", hs)
	}
}

func TestLoadCustom_DoesNotBreakBuiltins(t *testing.T) {
	path := writeYAML(t, "extras.yaml", `
months:
  6:
  - name: Custom June Holiday
    regions: [extras_test]
    mday: 10
`)
	if err := holidays.LoadCustom(path); err != nil {
		t.Fatalf("LoadCustom: %v", err)
	}
	defer holidays.UnloadCustom(path)

	d := time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(d, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if len(hs) == 0 || hs[0].Name != "Independence Day" {
		t.Fatalf("built-in us holiday lookup broken; got %+v", hs)
	}
}

func TestLoadCustom_AddsToExistingRegion(t *testing.T) {
	// A user can extend an existing region (e.g. "us") by writing a rule with
	// regions: [us]. The custom rule appears alongside the built-in ones.
	path := writeYAML(t, "us_extras.yaml", `
months:
  4:
  - name: Custom April Holiday
    regions: [us]
    mday: 1
`)
	if err := holidays.LoadCustom(path); err != nil {
		t.Fatalf("LoadCustom: %v", err)
	}
	defer holidays.UnloadCustom(path)

	d := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(d, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	found := false
	for _, h := range hs {
		if h.Name == "Custom April Holiday" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected Custom April Holiday in us results; got %+v", hs)
	}
}

func TestLoadCustom_UnregisteredFunctionErrors(t *testing.T) {
	path := writeYAML(t, "bad.yaml", `
months:
  3:
  - name: Bad
    regions: [bad_test]
    function: not_a_real_method(year)
`)
	if err := holidays.LoadCustom(path); err == nil {
		holidays.UnloadCustom(path)
		t.Fatal("expected error for unregistered method, got nil")
	}
}

func TestLoadCustom_RegisterMethodThenLoad(t *testing.T) {
	// Register a custom method that always returns Feb 29 (in leap years).
	holidays.RegisterMethod("test_feb_29", func(a holidays.MethodArgs) (time.Time, error) {
		return time.Date(a.Year, time.February, 29, 0, 0, 0, 0, time.UTC), nil
	})

	path := writeYAML(t, "leapday.yaml", `
months:
  2:
  - name: Leap Day
    regions: [leap_test]
    function: test_feb_29(year)
`)
	if err := holidays.LoadCustom(path); err != nil {
		t.Fatalf("LoadCustom: %v", err)
	}
	defer holidays.UnloadCustom(path)

	d := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(d, holidays.Options{Regions: []string{"leap_test"}})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if len(hs) != 1 || hs[0].Name != "Leap Day" {
		t.Fatalf("expected Leap Day; got %+v", hs)
	}
}

func TestLoadCustom_NoPaths(t *testing.T) {
	if err := holidays.LoadCustom(); err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestLoadCustom_MissingFile(t *testing.T) {
	if err := holidays.LoadCustom("/nonexistent/path.yaml"); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadCustom_ReloadReplaces(t *testing.T) {
	// Same basename, different contents → second load replaces first.
	dir := t.TempDir()
	first := filepath.Join(dir, "swap.yaml")
	if err := os.WriteFile(first, []byte(`
months:
  5:
  - name: First Version
    regions: [swap_test]
    mday: 1
`), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := holidays.LoadCustom(first); err != nil {
		t.Fatalf("LoadCustom 1: %v", err)
	}

	if err := os.WriteFile(first, []byte(`
months:
  5:
  - name: Second Version
    regions: [swap_test]
    mday: 1
`), 0o644); err != nil {
		t.Fatalf("write 2: %v", err)
	}
	if err := holidays.LoadCustom(first); err != nil {
		t.Fatalf("LoadCustom 2: %v", err)
	}
	defer holidays.UnloadCustom(first)

	d := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(d, holidays.Options{Regions: []string{"swap_test"}})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if len(hs) != 1 || hs[0].Name != "Second Version" {
		t.Fatalf("expected Second Version only; got %+v", hs)
	}
}

func TestUnloadCustom_RemovesRules(t *testing.T) {
	path := writeYAML(t, "temp.yaml", `
months:
  9:
  - name: Will Be Removed
    regions: [removable_test]
    mday: 5
`)
	if err := holidays.LoadCustom(path); err != nil {
		t.Fatalf("LoadCustom: %v", err)
	}

	d := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)
	hs, _ := holidays.On(d, holidays.Options{Regions: []string{"removable_test"}})
	if len(hs) != 1 {
		t.Fatalf("expected 1 before unload, got %d", len(hs))
	}

	holidays.UnloadCustom(path)

	hs, _ = holidays.On(d, holidays.Options{Regions: []string{"removable_test"}})
	if len(hs) != 0 {
		t.Fatalf("expected 0 after unload, got %d", len(hs))
	}
}

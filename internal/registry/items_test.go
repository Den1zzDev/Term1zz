package registry

import (
	"testing"

	"github.com/Den1zzDev/Term1zz/internal/distro"
)

func TestItemsIntegrity(t *testing.T) {
	items := Items()
	if len(items) == 0 {
		t.Fatal("Items() returned empty catalog")
	}

	seenIDs := make(map[string]bool)
	for _, it := range items {
		if it.ID == "" {
			t.Errorf("item with Name %q has empty ID", it.Name)
		}
		if seenIDs[it.ID] {
			t.Errorf("duplicate item ID %q found", it.ID)
		}
		seenIDs[it.ID] = true

		if it.Name == "" {
			t.Errorf("item ID %q has empty Name", it.ID)
		}
		if it.Description == "" {
			t.Errorf("item ID %q has empty Description", it.ID)
		}
		if it.Category == "" {
			t.Errorf("item ID %q has empty Category", it.ID)
		}
	}
}

func TestAllCategories(t *testing.T) {
	cats := AllCategories()
	if len(cats) != 7 {
		t.Errorf("expected 7 categories, got %d", len(cats))
	}
}

func TestNerdFontsFallback(t *testing.T) {
	items := Items()
	var nf *Item
	for i := range items {
		if items[i].ID == "nerdfonts" {
			nf = &items[i]
			break
		}
	}
	if nf == nil {
		t.Fatal("nerdfonts item not found")
	}
	if !nf.CustomScriptFallback {
		t.Error("nerdfonts must have CustomScriptFallback set to true")
	}
	if nf.CustomScript == "" {
		t.Error("nerdfonts must provide fallback CustomScript")
	}
	// Check package mappings
	if pkg, ok := nf.PackageFor(distro.PMPacman); !ok || pkg == "" {
		t.Error("nerdfonts missing pacman mapping")
	}
	if pkg, ok := nf.PackageFor(distro.PMBrew); !ok || pkg == "" {
		t.Error("nerdfonts missing brew mapping")
	}
}

func TestMossExclusions(t *testing.T) {
	items := Items()
	for _, it := range items {
		if it.ID == "navi" || it.ID == "xh" {
			if pkg, ok := it.PackageFor(distro.PMMoss); ok && pkg != "" {
				t.Errorf("%s must not be mapped to moss because it is not in the moss repository", it.ID)
			}
		}
	}
}

func TestFisherSetup(t *testing.T) {
	items := Items()
	var fisher *Item
	for i := range items {
		if items[i].ID == "fisher" {
			fisher = &items[i]
			break
		}
	}
	if fisher == nil {
		t.Fatal("fisher item not found")
	}
	if fisher.CustomScript == "" {
		t.Error("fisher must have CustomScript")
	}
}

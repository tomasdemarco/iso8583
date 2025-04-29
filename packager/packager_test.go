package packager

import (
	"testing"
)

// TestLoadFromJson calls message.UnpackPadding
func TestLoadFromJson(t *testing.T) {
	fileName := "iso87BPackager.json"
	_, err := LoadFromJson("./", "iso87BPackager.json")
	if err != nil {
		t.Fatalf(`TestLoadFromJson - error loading packager "%s": %v`, fileName, err)
	}

	t.Logf(`TestLoadFromJson - packager "%s" uploaded successfully`, fileName)

	fileName = "iso93EAmexPackager.json"

	_, err = LoadFromJson("./", fileName)
	if err != nil {
		t.Fatalf(`TestLoadFromJson - error loading packager "%s": %v`, fileName, err)
	}

	t.Logf(`TestLoadFromJson - packager "%s" uploaded successfully`, fileName)
}

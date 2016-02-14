package translate

import (
	"flag"
	"os"
	"testing"
)

var validAPIKey = os.Getenv("YANDEX_TRANSLATE_API_KEY")

func TestMain(m *testing.M) {
	flag.Parse()

	if validAPIKey == "" {
		// skip if no key is provided for tests
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func newTranslator() *Translator {
	return New(validAPIKey)
}

func TestGetLangs(t *testing.T) {
	tr := newTranslator()

	langs, err := tr.Languages("fr")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if langs == nil {
		t.Fatalf("expected non-nil response")
	}

	t.Logf("got languages %+v", langs)
}

func TestDetect(t *testing.T) {
	tr := newTranslator()

	lang, err := tr.Detect("Hello world!")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if lang == "" {
		t.Fatalf("unexpected empty string")
	}

	t.Logf("detected language %q", lang)
}

func TestTranslate(t *testing.T) {
	tr := newTranslator()

	translation, err := tr.Translate("en", "fr", "Hello world!")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if translation == nil {
		t.Fatalf("unexpected nil translation")
	}

	t.Logf("got translation %+v", translation)
}

func TestTranslateDetect(t *testing.T) {
	tr := newTranslator()

	translation, err := tr.Translate("", "fr", "Hello world!")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if translation == nil {
		t.Fatalf("unexpected nil translation")
	}

	t.Logf("got translation %+v", translation)
}

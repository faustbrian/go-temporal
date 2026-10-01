package temporalwire

import (
	"encoding/json"
	"testing"
)

func TestJSONStringSizeMatchesStandardEncodingWithoutChangingRepresentation(t *testing.T) {
	for _, value := range []string{"", "08:00", "quoted\"\\", "\b\f\n\r\t\x01", "<>&", "é\u2028\u2029", "\ufffd"} {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if got := jsonStringSize(value); got != len(encoded) {
			t.Fatalf("encoded size = %d, want %d", got, len(encoded))
		}
	}
}

func TestCollectionSizeAdmissionPreservesNilAndEmptyEncoding(t *testing.T) {
	for name, values := range map[string][]string{"nil": nil, "empty": {}} {
		t.Run(name, func(t *testing.T) {
			document := CollectionDocument{Version: Version1, Kind: KindInstantSet, Values: values}
			encoded, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			for _, maximum := range []int{len(encoded), len(encoded) - 1} {
				size, exceeds := collectionDocumentSize(document, maximum)
				if size != len(encoded) || exceeds != (maximum < len(encoded)) {
					t.Fatalf("size admission at %d = %d, %t; encoded size %d", maximum, size, exceeds, len(encoded))
				}
			}
		})
	}
}

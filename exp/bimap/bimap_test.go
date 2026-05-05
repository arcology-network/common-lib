package bimap

import "testing"

func TestBiMapSetGetDelete(t *testing.T) {
	this := NewBiMap[string, int]()
	this.Set("alpha", 1)

	if value, ok := this.GetByKey("alpha"); !ok || value != 1 {
		t.Error("Error: GetByKey should return the inserted value")
	}

	if key, ok := this.GetByValue(1); !ok || key != "alpha" {
		t.Error("Error: GetByValue should return the inserted key")
	}

	this.Set("alpha", 2)
	if _, ok := this.GetByValue(1); ok {
		t.Error("Error: replacing a key should remove the old reverse mapping")
	}

	this.Set("beta", 2)
	if _, ok := this.GetByKey("alpha"); ok {
		t.Error("Error: assigning an existing value to a new key should evict the old key")
	}

	if key, ok := this.GetByValue(2); !ok || key != "beta" {
		t.Error("Error: reverse mapping should point to the latest key")
	}

	this.DeleteByKey("beta")
	if _, ok := this.GetByValue(2); ok {
		t.Error("Error: DeleteByKey should remove the reverse mapping")
	}

	this.Set("gamma", 3)
	this.DeleteByValue(3)
	if _, ok := this.GetByKey("gamma"); ok {
		t.Error("Error: DeleteByValue should remove the forward mapping")
	}
}
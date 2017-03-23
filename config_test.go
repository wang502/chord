package chord

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	conf := DefaultConfig("localhost")
	if conf.Host != "localhost" {
		t.Errorf("wrong host")
	}
	if conf.HashFunc == nil {
		t.Errorf("bad hash func")
	}
	if conf.HashBits != 3 {
		t.Errorf("wrong hash bits")
	}
	if conf.NumNodes != 8 {
		t.Errorf("wrong number of nodes")
	}
}

func TestHashFunc(t *testing.T) {
	conf := DefaultConfig("localhost")
	hash := conf.HashFunc
	_, err := hash.Write([]byte(conf.Host))
	if err != nil {
		t.Errorf("bad hash function in config")
	}
	bytes := hash.Sum(nil)
	if len(bytes) == 0 {
		t.Errorf("bad result of hash.Sum()")
	}
}

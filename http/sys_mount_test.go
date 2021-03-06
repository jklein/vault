package http

import (
	"reflect"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/vault"
)

func TestSysMounts(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysMount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysMount_put(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	// The TestSysMount test tests the thing is actually created. See that test
	// for more info.
}

func TestSysRemount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "foo",
		"to":   "bar",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"bar/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysUnmount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/mounts/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysTuneMount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}

	// Shorter than system default
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": time.Duration(time.Hour * 72),
	})
	testResponseStatus(t, resp, 204)

	// Longer than system max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": time.Duration(time.Hour * 72000),
	})
	testResponseStatus(t, resp, 400)

	// Longer than system default
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"max_lease_ttl": time.Duration(time.Hour * 72000),
	})
	testResponseStatus(t, resp, 204)

	// Longer than backend max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": time.Duration(time.Hour * 72001),
	})
	testResponseStatus(t, resp, 400)

	// Shorter than backend default
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"max_lease_ttl": time.Duration(time.Hour * 1),
	})
	testResponseStatus(t, resp, 400)

	// Shorter than backend max, longer than system max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": time.Duration(time.Hour * 71999),
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")
	expected = map[string]interface{}{
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(time.Duration(time.Hour * 71999)),
				"max_lease_ttl":     float64(time.Duration(time.Hour * 72000)),
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Check simple configuration endpoint
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/foo/tune")
	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"default_lease_ttl": float64(time.Duration(time.Hour * 71999)),
		"max_lease_ttl":     float64(time.Duration(time.Hour * 72000)),
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Set a low max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"default_lease_ttl": time.Duration(time.Second * 40),
		"max_lease_ttl":     time.Duration(time.Second * 80),
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"default_lease_ttl": float64(time.Duration(time.Second * 40)),
		"max_lease_ttl":     float64(time.Duration(time.Second * 80)),
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// First try with lease above backend max
	resp = testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
		"ttl":  "28347h",
	})
	testResponseStatus(t, resp, 204)

	// read secret
	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	var result struct {
		LeaseID       string `json:"lease_id" structs:"lease_id"`
		LeaseDuration int    `json:"lease_duration" structs:"lease_duration"`
	}

	testResponseBody(t, resp, &result)

	expected = map[string]interface{}{
		"lease_duration": int(time.Duration(time.Second * 80).Seconds()),
		"lease_id":       result.LeaseID,
	}

	if !reflect.DeepEqual(structs.Map(result), expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, structs.Map(result))
	}

	// Now with lease TTL unspecified
	resp = testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 204)

	// read secret
	resp = testHttpGet(t, token, addr+"/v1/secret/foo")

	testResponseBody(t, resp, &result)

	expected = map[string]interface{}{
		"lease_duration": int(time.Duration(time.Second * 40).Seconds()),
		"lease_id":       result.LeaseID,
	}

	if !reflect.DeepEqual(structs.Map(result), expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, structs.Map(result))
	}
}

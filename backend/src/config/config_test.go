package config

import (
	"os"
	"path/filepath"
	"testing"
)

const testConfig = `
[COMMON]
port = 8180
open_tls = false

[MySQL]
url = "127.0.0.1:3307"
user = "root"
pwd = "1234"
db_name = "offer_hub"

[MongoDB]
url = "mongodb://127.0.0.1:27017"
user = ""
pwd = ""
database = "offer_hub"
max_pool = 100
min_pool = 5

[Redis]
url = "127.0.0.1:6379"
pwd = ""
db = 0

[JWT]
secret = "test-secret"
expire = 24
token_cache_expire = 24
enable = true

[RateLimit]
enable = true
window_seconds = 60
max_requests = 20
`

func TestInitWithExplicitPath(t *testing.T) {
	path := writeTestConfig(t, "custom.toml")

	if err := Init(path); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if Conf.Common.Port != 8180 {
		t.Fatalf("Common.Port = %d, want 8180", Conf.Common.Port)
	}
	if Conf.MySQL.DSN() != "root:1234@tcp(127.0.0.1:3307)/offer_hub?charset=utf8mb4&parseTime=True&loc=Local" {
		t.Fatalf("unexpected MySQL DSN: %s", Conf.MySQL.DSN())
	}
	if !Conf.JWT.Enable || Conf.JWT.Expire != 24 || Conf.JWT.TokenCacheExpire != 24 {
		t.Fatalf("unexpected JWT config: %+v", Conf.JWT)
	}
	if !Conf.RateLimit.Enable || Conf.RateLimit.MaxRequests != 20 {
		t.Fatalf("unexpected RateLimit config: %+v", Conf.RateLimit)
	}
}

func TestInitUsesAppEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	directory := filepath.Dir(writeTestConfig(t, "config-production.toml"))

	if err := Init(directory); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if Conf.MongoDB.Database != "offer_hub" {
		t.Fatalf("MongoDB.Database = %q, want offer_hub", Conf.MongoDB.Database)
	}
}

func TestConfigFileNameDefaultsToTest(t *testing.T) {
	t.Setenv("APP_ENV", "")

	if name := configFileName(); name != "config-test.toml" {
		t.Fatalf("configFileName() = %q, want config-test.toml", name)
	}
}

func writeTestConfig(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(testConfig), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

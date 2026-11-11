package config

import "testing"

func cleanEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"APP_ENV", "HTTP_ADDR", "DATABASE_URL", "LOG_LEVEL", "HTTP_READ_TIMEOUT", "HTTP_WRITE_TIMEOUT", "HTTP_IDLE_TIMEOUT", "SHUTDOWN_TIMEOUT"} {
		t.Setenv(key, "")
	}
}

func TestProductionRequiresDatabase(t *testing.T) {
	cleanEnv(t)
	t.Setenv("APP_ENV", "production")
	if _, err := Load(); err == nil {
		t.Fatal("production accepted missing database configuration")
	}
}

func TestRejectsInvalidConfiguration(t *testing.T) {
	for _, tc := range []struct{ key, value string }{
		{"APP_ENV", "prodution"}, {"HTTP_ADDR", "8080"},
		{"LOG_LEVEL", "everything"}, {"HTTP_READ_TIMEOUT", "0s"},
		{"SHUTDOWN_TIMEOUT", "-1s"}, {"HTTP_IDLE_TIMEOUT", "forever"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			cleanEnv(t)
			t.Setenv(tc.key, tc.value)
			if _, err := Load(); err == nil {
				t.Fatalf("accepted %s=%s", tc.key, tc.value)
			}
		})
	}
}

func TestDevelopmentDefaults(t *testing.T) {
	cleanEnv(t)
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.HTTPAddr != ":8080" || c.DatabaseURL != "" {
		t.Fatalf("unexpected development defaults: %+v", c)
	}
}

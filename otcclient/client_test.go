package otcclient

import (
	"strings"
	"testing"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
)

// noopLogger is a no-op implementation of internal.ILogger used as the default
// logger so that callers of c.Logger.Info() etc. never panic on a nil receiver.
type noopLogger struct{}

func (n noopLogger) Info(_ string, _ ...interface{})              {}
func (n noopLogger) Debug(_ string, _ ...interface{})             {}
func (n noopLogger) Warn(_ string, _ ...interface{})              {}
func (n noopLogger) Error(_ string, _ ...interface{})             {}
func (n noopLogger) Panic(_ string, _ ...interface{})             {}
func (n noopLogger) Sync() error                                  { return nil }
func (n noopLogger) WithFields(_ ...interface{}) internal.ILogger { return n }

func validUserPassConfig() Config {
	return Config{
		Username:   "user",
		Password:   "pass",
		DomainName: "domain",
		ProjectID:  "project-123",
		Region:     "eu-de",
	}
}

func validAKSKConfig() Config {
	return Config{
		AccessKey: "ak",
		SecretKey: "sk",
		ProjectID: "project-123",
		Region:    "eu-nl",
	}
}

func TestNewClientRejectsEmptyConfig(t *testing.T) {
	_, err := New(Config{}, noopLogger{})
	if err == nil {
		t.Fatal("expected error for empty config, got nil")
	}
}

func TestNewClientRejectsInvalidRegion(t *testing.T) {
	cfg := validUserPassConfig()
	cfg.Region = "us-east-1"
	_, err := New(cfg, noopLogger{})
	if err == nil {
		t.Fatal("expected error for invalid region, got nil")
	}
	if !strings.Contains(err.Error(), "invalid region") {
		t.Fatalf("expected 'invalid region' in error, got: %s", err.Error())
	}
}

func TestNewClientRejectsMixedAuth(t *testing.T) {
	cfg := Config{
		Username:   "user",
		Password:   "pass",
		AccessKey:  "ak",
		SecretKey:  "sk",
		DomainName: "domain",
		ProjectID:  "project-123",
		Region:     "eu-de",
	}
	_, err := New(cfg, noopLogger{})
	if err == nil {
		t.Fatal("expected error for mixed auth, got nil")
	}
	if !strings.Contains(err.Error(), "not both") {
		t.Fatalf("expected 'not both' in error, got: %s", err.Error())
	}
}

func TestNewClientRejectsNoAuth(t *testing.T) {
	cfg := Config{
		ProjectID: "project-123",
		Region:    "eu-de",
	}
	_, err := New(cfg, noopLogger{})
	if err == nil {
		t.Fatal("expected error for missing auth, got nil")
	}
	if !strings.Contains(err.Error(), "must provide") {
		t.Fatalf("expected 'must provide' in error, got: %s", err.Error())
	}
}

func TestNewClientRejectsMissingDomainForUserPass(t *testing.T) {
	cfg := validUserPassConfig()
	cfg.DomainName = ""
	_, err := New(cfg, noopLogger{})
	if err == nil {
		t.Fatal("expected error for missing domain name, got nil")
	}
	if !strings.Contains(err.Error(), "DomainName is required") {
		t.Fatalf("expected 'DomainName is required' in error, got: %s", err.Error())
	}
}

func TestValidateConfigAcceptsValidUserPass(t *testing.T) {
	if err := validateConfig(validUserPassConfig()); err != nil {
		t.Fatalf("expected valid user/pass config to pass validation, got: %s", err)
	}
}

func TestValidateConfigAcceptsValidAKSK(t *testing.T) {
	if err := validateConfig(validAKSKConfig()); err != nil {
		t.Fatalf("expected valid AK/SK config to pass validation, got: %s", err)
	}
}

func TestValidateConfigRejectsMissingProjectID(t *testing.T) {
	cfg := validUserPassConfig()
	cfg.ProjectID = ""
	err := validateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing ProjectID, got nil")
	}
	if !strings.Contains(err.Error(), "ProjectID") {
		t.Fatalf("expected 'ProjectID' in error, got: %s", err.Error())
	}
}

func TestValidateConfigRejectsMissingPassword(t *testing.T) {
	cfg := validUserPassConfig()
	cfg.Password = ""
	err := validateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing password, got nil")
	}
	if !strings.Contains(err.Error(), "password") {
		t.Fatalf("expected 'password' in error, got: %s", err.Error())
	}
}

func TestValidateConfigRejectsMissingUsername(t *testing.T) {
	cfg := validUserPassConfig()
	cfg.Username = ""
	err := validateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing username, got nil")
	}
	if !strings.Contains(err.Error(), "username") {
		t.Fatalf("expected 'username' in error, got: %s", err.Error())
	}
}

func TestValidateConfigRejectsPartialAKSK(t *testing.T) {
	cfg := validAKSKConfig()
	cfg.SecretKey = ""
	err := validateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing secret key, got nil")
	}
	if !strings.Contains(err.Error(), "secret-key") {
		t.Fatalf("expected 'secret-key' in error, got: %s", err.Error())
	}
}

func TestIamEndpoint(t *testing.T) {
	got := iamEndpoint("eu-de")
	want := "https://iam.eu-de.otc.t-systems.com/v3"
	if got != want {
		t.Fatalf("iamEndpoint(eu-de) = %q, want %q", got, want)
	}
}

func TestValidateConfigAcceptsSwissRegion(t *testing.T) {
	cfg := validUserPassConfig()
	cfg.Region = "eu-ch2"
	if err := validateConfig(cfg); err != nil {
		t.Fatalf("expected eu-ch2 config to pass validation, got: %s", err)
	}
}

func TestIamEndpointSwiss(t *testing.T) {
	got := iamEndpoint("eu-ch2")
	want := "https://iam-pub.eu-ch2.sc.otc.t-systems.com/v3"
	if got != want {
		t.Fatalf("iamEndpoint(eu-ch2) = %q, want %q", got, want)
	}
}

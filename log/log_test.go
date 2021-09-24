package log

import "testing"

func TestNewLog(t *testing.T) {
	logger, err := New(&Config{
		Level: "info",
	})
	if err != nil {
		t.Fatal(err)
	}

	sugger := logger.Sugar()
	sugger.Infow("hello wold", "key", "ray", "value", "good")
	sugger.Infof("this is test: %s, %v", "good", nil)
	sugger.Info("hello world")
}

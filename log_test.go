package logger

import (
	"bufio"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	Debug("234-aae", "test Debug", map[string]interface{}{"env": "test", "customer": "op"})
	traceInfoID := "123-abc"
	infoMessage := "test Info"
	Info(traceInfoID, infoMessage, map[string]interface{}{"env": "test", "customer": "op"})
	traceErrorID := "456-aae"
	errorMessage := "test Error"
	Error(traceErrorID, errorMessage, map[string]interface{}{"env": "test", "customer": "op"})

	outFile := "logs.json"
	assert.FileExists(t, outFile)
	file, err := os.Open(outFile)
	assert.NoError(t, err)

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	testChecks := []struct {
		traceID string
		message string
		level   string
	}{
		{traceID: traceInfoID, message: infoMessage, level: getLevelAsString(InfoLevel)},
		{traceID: traceErrorID, message: errorMessage, level: getLevelAsString(ErrorLevel)},
	}

	i := 0
	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		attributes := make(map[string]interface{})
		err = json.Unmarshal(line, &attributes)
		if err != nil {
			t.Errorf("Error unmarshalling logs.json line: %v", err)
		}
		assert.Equal(t, testChecks[i].traceID, attributes["traceID"])
		assert.Equal(t, testChecks[i].message, attributes["message"])
		assert.Equal(t, testChecks[i].level, attributes["level"])
		i++
	}

	file.Close()
	os.Remove(outFile)
}

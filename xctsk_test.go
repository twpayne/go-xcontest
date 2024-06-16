package xcontest_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/twpayne/go-xctrack"

	"github.com/twpayne/go-xcontest"
)

func TestClientXCTskSaveLoad(t *testing.T) {
	t.Skip("uses xcontest.org servers") // Comment out this line to enable this test.

	ctx := context.Background()
	client := xcontest.NewClient()
	task := loadTask(t, "testdata/zugerberg-zurich.xctsk")
	qrCodeTask := task.QRCodeTask()
	author := "Tom Payne"

	taskCode, err := client.XCTskSave(ctx, qrCodeTask, author)
	assert.NoError(t, err)

	t.Logf("taskCode=%s", taskCode)

	t.Run("load", func(t *testing.T) {
		savedTask, err := client.XCTskLoad(ctx, taskCode)
		assert.NoError(t, err)
		expectedTask := *task
		expectedTask.SSS.Direction = xctrack.DirectionExit // Returned task has SSS direction exit.
		assert.Equal(t, &expectedTask, savedTask.Task)
		assert.Equal(t, author, savedTask.Author)
	})

	t.Run("loadV2", func(t *testing.T) {
		savedTaskV2, err := client.XCTskLoadV2(ctx, taskCode)
		assert.NoError(t, err)
		expectedTask := *qrCodeTask
		expectedTask.SSS.Direction = xctrack.QRCodeDirectionExit // Returned task has SSS direction exit.
		assert.Equal(t, qrCodeTask, savedTaskV2.Task)
		assert.Equal(t, author, savedTaskV2.Author)
	})

	t.Run("qr", func(t *testing.T) {
		qrCodeData, err := client.XCTskQR(ctx, qrCodeTask)
		assert.NoError(t, err)
		resetGoldenFiles := false
		if resetGoldenFiles {
			os.WriteFile("testdata/zugerberg-zurich.svg", qrCodeData, 0o666)
		} else {
			expectedQRCodeData, err := os.ReadFile("testdata/zugerberg-zurich.svg")
			assert.NoError(t, err)
			assert.Equal(t, expectedQRCodeData, qrCodeData)
		}
	})
}

func loadTask(t *testing.T, filename string) *xctrack.Task {
	t.Helper()
	data, err := os.ReadFile(filename)
	assert.NoError(t, err)
	var task xctrack.Task
	assert.NoError(t, json.Unmarshal(data, &task))
	return &task
}

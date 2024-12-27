package lib

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/spf13/viper"
)

func TestStorageDirectory(t *testing.T) {
	viper.Set("STORAGE_DIRECTORY", "demo")
	viper.Set("STORAGE_CREATE", true)
	utils.AssertEqual(t, "demo", StorageDirectory())

	viper.Set("STORAGE_DIRECTORY", "")
	workingDirectory, _ := os.Getwd()
	utils.AssertEqual(t, workingDirectory+"/uploads", StorageDirectory())
	os.RemoveAll(workingDirectory + "/uploads")
	os.RemoveAll(workingDirectory + "/demo")
}

func TestDirExists(t *testing.T) {
	utils.AssertEqual(t, false, DirExists("non-existing-directory"))
}

func TestFileExists(t *testing.T) {
	utils.AssertEqual(t, false, FileExists("non-existing-file"))
}

func TestGetMimeFile(t *testing.T) {
	utils.AssertEqual(t, "text/plain", GetMimeFile("file_manager.go").MIME.Value)

	sample := "example.tar.gz"
	os.WriteFile(sample, []byte(""), 0755)
	utils.AssertEqual(t, "tar.gz", GetMimeFile(sample).Extension)
	os.RemoveAll(sample)
	mime := GetMimeFile("non-existing-file")
	utils.AssertEqual(t, "", mime.MIME.Type)
}

func TestGetImageScaleSize(t *testing.T) {
	base64Image := `iVBORw0KGgoAAAANSUhEUgAAAA4AAAAOCAYAAAAfSC3RAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAAB3RJTUUH5QgYCCMuQ+w/NgAAAo1JREFUKM9Nkk1rVGcYhq/3Y2bOZKyRJDOWTKepH2CaoqILSUulFoVCF93EMlBF0k2hkE1/Q/+AlNk3ixKIKPiB4IiK1oWB0tak1aaiY21nUs2EsY7jzJlz3uc9LpKKz+bePDf3tbgUr93DlbXiUlNPL6zq8uOeHndasSVQy/tH/PzBgpvdOTrU+P9XASTSMhcWZepSjcr5ms63nMEbjTeGWGtygeHwW9L8pBjPfD0RnlHZgliAi4tuau5XmavWEjMQKHLa4hKD+IS00rgo4VzN5ld7as4n6gvglH5QXyveuBdVbj0IzXDKsTMXstX2sc5hRTAiWBGyOG6tKFP9x1burbSK+u+WTF9aCvM7NkVMDIZk8Rwc7WPFrZedbKQji+N6Q+WvNOy0XnwUlZ+1HUMZYU/BsdJWHJ2AN4ME/aokWCeknBCGwi9NyrrVlnGD50UIe0cNhSDh5v2Ek58qAgTlBCNunWCD4lk3Gdc68Vgv/NuCTl/z5aTl9CLUnzzn+89ituViUrK+bDYItHjM0RPffP7Tn+FW56EdQvlAmsFA8d0NTxy94Ph7HQZtTC9WdCJFLIrJUe7aPWPp+eEBdrcjz+2/HJXLHb46lGY4SPjh5wzV+wGbMp4eHi+ezWnHvkJqXpdGUrNH9g40fSRklXD1juPbsz3azyOO7e4y9W6Pd7aASjQu8nxUovnxmJq120r5RnWhPrP6NJ778W7PBCn4o6GorWkKgyneCKDrPP1uzPtvWzm8Xc/sKI001pVLnprqQmfq5u/dyrXfuvn/euCNwiuDV4qBrOWDXUHzw12ZmeOTmTMqMyTqdckf1ZvF5Xp/eulhv7zWSca9gs05szxRSs/vG8vMbi+NvJL8Je6JQY9chg00AAAAAElFTkSuQmCC`
	imgData, err := base64.StdEncoding.DecodeString(base64Image)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, false, nil == imgData)

	sample := "sample.png"
	os.WriteFile(sample, imgData, 0755)

	width, height, err := GetImageScaleSize(sample)
	os.RemoveAll(sample)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 14, width)
	utils.AssertEqual(t, 14, height)

	_, _, err = GetImageScaleSize("file_manager.go")
	utils.AssertEqual(t, false, nil == err)
}

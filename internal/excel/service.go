package excel

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	START_ROW     = 4
	MAX_ROWS      = 1000
	FILE_NAME     = "Implants.xlsx"
	DIR_NAME      = "SmileOpsTables"
	SHEET         = "Sheet1"
	TEMPLATE_FILE = "templates/Implants.xlsx"
)

// ExcelService handles all Excel operations for patient data
type ExcelService struct{}

// NewExcelService creates a new ExcelService instance
func NewExcelService() *ExcelService {
	return &ExcelService{}
}

// getExcelPath returns the full path to the Excel file on Desktop
func (s *ExcelService) getExcelPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	dirPath := filepath.Join(home, "Desktop", DIR_NAME)
	filePath := filepath.Join(dirPath, FILE_NAME)

	return filePath, nil
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (s *ExcelService) ensureDirectoryExists() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	dirPath := filepath.Join(home, "Desktop", DIR_NAME)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return nil
}

// copyTemplateFile copies the Excel template file to the destination path using streaming
func (s *ExcelService) copyTemplateFile(templatePath, destPath string) error {
	// Open source file for reading
	src, err := os.Open(templatePath)
	if err != nil {
		return fmt.Errorf("failed to open template file %s: %w", templatePath, err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer dst.Close()

	// Copy contents using streaming (memory efficient)
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy template from %s to %s: %w", templatePath, destPath, err)
	}

	// Ensure data is written to disk
	if err := dst.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// ensureFileExists checks if the Excel file exists and creates it if not
func (s *ExcelService) ensureFileExists() error {
	filePath, err := s.getExcelPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := s.ensureDirectoryExists(); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, copy from template
		// Get executable path to locate bundled template
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}
		exeDir := filepath.Dir(exe)

		templatePath := filepath.Join(exeDir, TEMPLATE_FILE)

		// Verify template exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("template file not found at %s", templatePath)
		}

		// Copy template to destination
		if err := s.copyTemplateFile(templatePath, filePath); err != nil {
			return fmt.Errorf("failed to copy Excel template: %w", err)
		}
	}

	return nil
}

// generateID finds the maximum ID in the Excel file and generates the next one
func (s *ExcelService) generateID() (string, error) {
	filePath, err := s.getExcelPath()
	if err != nil {
		return "", err
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	maxID := 0

	// Read rows starting from START_ROW
	rows, err := f.GetRows(SHEET)
	if err != nil {
		return "1", nil // If error reading, start with 1
	}

	for i := START_ROW - 1; i < len(rows) && i < MAX_ROWS; i++ {
		if len(rows[i]) > 0 {
			idStr := strings.TrimSpace(rows[i][0])
			if idStr != "" {
				if id, err := strconv.Atoi(idStr); err == nil {
					if id > maxID {
						maxID = id
					}
				}
			}
		}
	}

	return strconv.Itoa(maxID + 1), nil
}

// findPatientRow finds the row number for a patient with the given ID, or returns the first empty row
func (s *ExcelService) findPatientRow(f *excelize.File, patientID string) (int, error) {
	rows, err := f.GetRows(SHEET)
	if err != nil {
		return START_ROW, nil // Return start row if error
	}

	// If patient has an ID, find existing row
	if patientID != "" {
		for i := START_ROW - 1; i < len(rows) && i < MAX_ROWS; i++ {
			if len(rows[i]) > 0 {
				existingID := strings.TrimSpace(rows[i][0])
				if existingID == patientID {
					return i + 1, nil // Excel rows are 1-based
				}
			}
		}
	}

	// Find first empty row
	for i := START_ROW - 1; i < MAX_ROWS; i++ {
		if i >= len(rows) || len(rows[i]) == 0 || strings.TrimSpace(rows[i][0]) == "" {
			return i + 1, nil // Excel rows are 1-based
		}
	}

	return 0, fmt.Errorf("no available rows (max %d reached)", MAX_ROWS)
}

// SavePatient saves a patient to the Excel file
func (s *ExcelService) SavePatient(patient Patient) error {
	// Ensure file exists
	if err := s.ensureFileExists(); err != nil {
		return err
	}

	filePath, err := s.getExcelPath()
	if err != nil {
		return err
	}

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Generate ID if empty
	if patient.ID == "" {
		patient.ID, err = s.generateID()
		if err != nil {
			return err
		}
	}

	// Find the row to write to
	row, err := s.findPatientRow(f, patient.ID)
	if err != nil {
		return err
	}

	// Helper function to join string slices
	joinStrings := func(arr []string) string {
		return strings.Join(arr, ", ")
	}

	// Write patient data to the row
	f.SetCellValue(SHEET, fmt.Sprintf("A%d", row), patient.ID)
	f.SetCellValue(SHEET, fmt.Sprintf("B%d", row), patient.FIO)
	f.SetCellValue(SHEET, fmt.Sprintf("C%d", row), patient.ImplantNumber)

	// Write operations
	if removal, ok := patient.Operations["removal"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("D%d", row), joinStrings(removal.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("E%d", row), joinStrings(removal.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("F%d", row), removal.Comment)
	}

	if sinusLift, ok := patient.Operations["sinusLift"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("G%d", row), joinStrings(sinusLift.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("H%d", row), joinStrings(sinusLift.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("I%d", row), sinusLift.Comment)
	}

	if boneGrafting, ok := patient.Operations["boneGrafting"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("J%d", row), joinStrings(boneGrafting.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("K%d", row), joinStrings(boneGrafting.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("L%d", row), boneGrafting.Comment)
	}

	if installationFormation, ok := patient.Operations["installationFormation"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("M%d", row), joinStrings(installationFormation.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("N%d", row), joinStrings(installationFormation.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("O%d", row), installationFormation.Comment)
	}

	if reinstallationImplant, ok := patient.Operations["reinstallationImplant"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("P%d", row), joinStrings(reinstallationImplant.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("Q%d", row), joinStrings(reinstallationImplant.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("R%d", row), reinstallationImplant.Comment)
	}

	if permanentProsthetics, ok := patient.Operations["permanentProsthetics"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("S%d", row), joinStrings(permanentProsthetics.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("T%d", row), joinStrings(permanentProsthetics.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("U%d", row), permanentProsthetics.Comment)
	}

	if temporaryProsthetics, ok := patient.Operations["temporaryProsthetics"]; ok {
		f.SetCellValue(SHEET, fmt.Sprintf("V%d", row), joinStrings(temporaryProsthetics.Dates))
		f.SetCellValue(SHEET, fmt.Sprintf("W%d", row), joinStrings(temporaryProsthetics.Numbers))
		f.SetCellValue(SHEET, fmt.Sprintf("X%d", row), temporaryProsthetics.Comment)
	}

	// Write control fields
	f.SetCellValue(SHEET, fmt.Sprintf("Y%d", row), patient.ControlHalfYear)
	f.SetCellValue(SHEET, fmt.Sprintf("Z%d", row), patient.ControlYear)
	f.SetCellValue(SHEET, fmt.Sprintf("AA%d", row), patient.OccupationalHygiene)

	// Save the file
	if err := f.Save(); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// GetAllPatients reads all patients from the Excel file
func (s *ExcelService) GetAllPatients() ([]Patient, error) {
	// Ensure file exists
	if err := s.ensureFileExists(); err != nil {
		return nil, err
	}

	filePath, err := s.getExcelPath()
	if err != nil {
		return nil, err
	}

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(SHEET)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	var patients []Patient

	// Helper function to split strings
	splitStrings := func(s string) []string {
		if s == "" {
			return []string{}
		}
		parts := strings.Split(s, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	// Helper function to safely get cell value
	getCell := func(row []string, col int) string {
		if col < len(row) {
			return strings.TrimSpace(row[col])
		}
		return ""
	}

	// Read data rows starting from START_ROW
	for i := START_ROW - 1; i < len(rows) && i < MAX_ROWS; i++ {
		row := rows[i]
		if len(row) == 0 || getCell(row, 0) == "" {
			continue // Skip empty rows
		}

		// Parse implant number
		implantNumber := 0
		if numStr := getCell(row, 2); numStr != "" {
			if num, err := strconv.Atoi(numStr); err == nil {
				implantNumber = num
			}
		}

		// Create operations map
		operations := make(map[string]Operation)

		operations["removal"] = Operation{
			Dates:   splitStrings(getCell(row, 3)),
			Numbers: splitStrings(getCell(row, 4)),
			Comment: getCell(row, 5),
		}

		operations["sinusLift"] = Operation{
			Dates:   splitStrings(getCell(row, 6)),
			Numbers: splitStrings(getCell(row, 7)),
			Comment: getCell(row, 8),
		}

		operations["boneGrafting"] = Operation{
			Dates:   splitStrings(getCell(row, 9)),
			Numbers: splitStrings(getCell(row, 10)),
			Comment: getCell(row, 11),
		}

		operations["installationFormation"] = Operation{
			Dates:   splitStrings(getCell(row, 12)),
			Numbers: splitStrings(getCell(row, 13)),
			Comment: getCell(row, 14),
		}

		operations["reinstallationImplant"] = Operation{
			Dates:   splitStrings(getCell(row, 15)),
			Numbers: splitStrings(getCell(row, 16)),
			Comment: getCell(row, 17),
		}

		operations["permanentProsthetics"] = Operation{
			Dates:   splitStrings(getCell(row, 18)),
			Numbers: splitStrings(getCell(row, 19)),
			Comment: getCell(row, 20),
		}

		operations["temporaryProsthetics"] = Operation{
			Dates:   splitStrings(getCell(row, 21)),
			Numbers: splitStrings(getCell(row, 22)),
			Comment: getCell(row, 23),
		}

		patient := Patient{
			ID:                  getCell(row, 0),
			FIO:                 getCell(row, 1),
			ImplantNumber:       implantNumber,
			Operations:          operations,
			ControlHalfYear:     getCell(row, 24),
			ControlYear:         getCell(row, 25),
			OccupationalHygiene: getCell(row, 26),
		}

		patients = append(patients, patient)
	}

	return patients, nil
}

package devscripts

// Backward compatible wrapper - delegates to new implementation
// Deprecated: Use DevScriptsReadmeUpdater instead

// UpdateReadmeIfNeeded updates the README file with the scripts section if there are changes
func (sr *scriptRunner) UpdateReadmeIfNeeded(readmePath string) (bool, error) {
	updater := NewDevScriptsReadmeUpdater(sr.scriptsDir)
	return updater.UpdateReadmeIfNeeded(readmePath)
}

// GenerateReadmeSection generates a markdown section for README
func (sr *scriptRunner) GenerateReadmeSection() (string, error) {
	updater := NewDevScriptsReadmeUpdater(sr.scriptsDir)
	return updater.GenerateScriptsSection()
}

// GetScriptNames obtiene los nombres de los scripts .sh en el directorio
func (sr *scriptRunner) GetScriptNames() ([]string, error) {
	parser := NewScriptParser(sr.scriptsDir)
	return parser.GetScriptNames()
}

package devscripts

import (
	"os"
)

// ExecuteWithArgs ejecuta una función después de cambiar al directorio especificado
// El primer argumento siempre debe ser el directorio de trabajo
// Los argumentos restantes se pasan a la función
func ExecuteWithArgs(fn func(...string)) {
	// El primer argumento siempre debe ser el directorio de ejecución
	if len(os.Args) > 1 {
		workDir := os.Args[1]
		err := os.Chdir(workDir)
		if err != nil {
			println("Error cambiando al directorio:", err.Error())
			os.Exit(1)
		}
	}

	// Los argumentos adicionales están en os.Args[2:]
	extraArgs := os.Args[2:]

	// Ejecutar la función con los argumentos adicionales
	fn(extraArgs...)
}

// GetWorkingDirectory retorna el directorio de trabajo del primer argumento
func GetWorkingDirectory() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}

// GetExtraArgs retorna los argumentos adicionales (después del directorio)
func GetExtraArgs() []string {
	if len(os.Args) > 2 {
		return os.Args[2:]
	}
	return []string{}
}

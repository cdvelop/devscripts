# ISSUE AUTO BADGES - Documentación del Sistema de Badges Automáticos

## Resumen del Proyecto
Sistema automatizado para generar y actualizar badges de estado en archivos README.md de proyectos Go, integrado en el flujo de trabajo de desarrollo.

## Arquitectura del Sistema

### Flujo de Ejecución
```
gomodcheck.sh → gotest.sh → gobadge.sh
                    ↓
               license.sh (auxiliar)
               gomodutils.sh (auxiliar)
```

### Scripts Principales

#### 1. `gotest.sh` (Independiente)
**Propósito**: Ejecutar todas las pruebas de Go y recopilar métricas
**Responsabilidades**:
- Ejecutar `go vet ./...`
- Ejecutar `go test ./...`
- Ejecutar `go test -race ./...`
- Calcular cobertura de tests con `go test -cover ./...`
- Recopilar resultados y llamar a `gobadge.sh`
- Continuar ejecución aunque fallen las pruebas
- Retornar código de salida 1 si hay errores (después de actualizar README)

#### 2. `gobadge.sh` (Independiente)
**Propósito**: Actualizar badges en README.md
**Parámetros recibidos** (en orden):
1. `$1` - Nombre del módulo Go
2. `$2` - Estado de tests ("Passing"/"Failed")
3. `$3` - Porcentaje de cobertura (número entero)
4. `$4` - Estado de race conditions ("Clean"/"Detected")
5. `$5` - Estado de go vet ("OK"/"Issues")
6. `$6` - Tipo de licencia (opcional, si no se pasa busca LICENSE)

**Comportamiento**:
- Si no existe README.md: mostrar warning y continuar
- Si existe README.md sin título `#`: agregar título con nombre del módulo
- Buscar primer título `# Título` y colocar badges inmediatamente después
- Reemplazar badges existentes si ya existen (detectar por comentario)

#### 3. `license.sh` (Auxiliar)
**Propósito**: Detectar tipo de licencia
**Comportamiento**:
- Buscar archivos en orden: LICENSE.txt, LICENSE, LICENSE.md
- Extraer primera palabra después de remover "License"
- Ejemplos:
  - "MIT License" → "MIT"
  - "Apache License 2.0" → "Apache"
  - "GNU General Public License" → "GNU"
- Si no encuentra archivo: retornar "MIT" por defecto

### Integración con gomodutils.sh
**Función nueva**: `get_go_version`
- Extraer versión de go.mod (línea `go 1.xx`)
- Retornar solo el número: "1.22" (sin el "+")

## Formato de Badges

### HTML Template
```html
<!-- Generated dynamically by gotest.sh from github.com/cdvelop/devscripts -->
<div style="display: flex;">
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px;">License</span><span style="background-color: #007acc; color: white; padding: 4px 8px; font-size: 12px;">MIT</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Go</span><span style="background-color: #00add8; color: white; padding: 4px 8px; font-size: 12px;">1.22</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Tests</span><span style="background-color: #28a745; color: white; padding: 4px 8px; font-size: 12px;">Passing</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Coverage</span><span style="background-color: #ffc107; color: white; padding: 4px 8px; font-size: 12px;">85%</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Race</span><span style="background-color: #28a745; color: white; padding: 4px 8px; font-size: 12px;">Clean</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Vet</span><span style="background-color: #28a745; color: white; padding: 4px 8px; font-size: 12px;">OK</span>
</div>
```

### Reglas de Formato
- **Sin espacios** entre título y valor de cada badge
- **Espacio de 1ch** entre diferentes badges (usando margin-left)
- **Títulos**: Fondo gris (#6c757d), texto blanco
- **Valores**: Color de fondo según estado, texto blanco siempre

### Colores por Estado

#### License
- **MIT/Apache/GNU/etc**: Azul (#007acc)

#### Go Version
- **Cualquier versión**: Azul Go (#00add8)

#### Tests
- **Passing**: Verde (#28a745)
- **Failed**: Rojo (#dc3545)

#### Coverage
- **100%**: Verde (#28a745)
- **<100%**: Amarillo (#ffc107)
- **Sin cobertura/Error**: Rojo (#dc3545)

#### Race Detection
- **Clean**: Verde (#28a745)
- **Detected**: Rojo (#dc3545)

#### Go Vet
- **OK**: Verde (#28a745)
- **Issues**: Rojo (#dc3545)

## Cálculo de Cobertura
- Usar `go test -cover ./...` para obtener cobertura por paquete
- Si hay múltiples paquetes: calcular promedio
- Ejemplo: 3 paquetes con 80%, 60%, 40% → (80+60+40)/3 = 60%
- Redondear resultado a número entero

## Manejo de Errores
- **Filosofía**: Continuar ejecución aunque haya errores
- **Actualizar README**: Siempre actualizar badges con el estado real
- **Código de salida**: Retornar 1 al final si hubo errores
- **Mensajes**: Acumular errores y mostrarlos al final (similar a pu.sh)

## Modificaciones en Scripts Existentes

### gomodcheck.sh
- Remover ejecución directa de go vet, go test, go test -race
- Llamar a `gotest.sh` en su lugar
- Mantener otras responsabilidades (go mod tidy, syscall, etc.)

## Testing
- **gotest_test.go**: Pruebas para gotest.sh
- **gobadge_test.go**: Pruebas para gobadge.sh
- **Referencia**: Usar doingmdfile_test.go como template
- **Método**: Crear wrappers que llamen a los scripts en directorios temporales

## Ubicación de Scripts
- **Directorio**: `c:\Users\Cesar\Packages\Internal\devscripts\`
- **Disponibilidad**: Scripts globalmente disponibles en Git Bash

## Ejemplo de README.md Resultante

### Antes
```markdown
# MyProject

Este es un proyecto de ejemplo...
```

### Después
```markdown
# MyProject
<!-- Generated dynamically by gotest.sh from github.com/cdvelop/devscripts -->
<div style="display: flex;">
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px;">License</span><span style="background-color: #007acc; color: white; padding: 4px 8px; font-size: 12px;">MIT</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Go</span><span style="background-color: #00add8; color: white; padding: 4px 8px; font-size: 12px;">1.22</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Tests</span><span style="background-color: #28a745; color: white; padding: 4px 8px; font-size: 12px;">Passing</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Coverage</span><span style="background-color: #ffc107; color: white; padding: 4px 8px; font-size: 12px;">85%</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Race</span><span style="background-color: #28a745; color: white; padding: 4px 8px; font-size: 12px;">Clean</span>
    <span style="background-color: #6c757d; color: white; padding: 4px 8px; font-size: 12px; margin-left: 1ch;">Vet</span><span style="background-color: #28a745; color: white; padding: 4px 8px; font-size: 12px;">OK</span>
</div>

Este es un proyecto de ejemplo...
```

---
**Fecha de creación**: 2025-06-06  
**Autor**: Sistema automatizado de desarrollo  
**Versión**: 1.0

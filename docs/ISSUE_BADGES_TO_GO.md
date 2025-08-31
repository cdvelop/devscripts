Respuestas preliminares del autor

- Para la pregunta sobre `README update`: "sectionUpdate.sh es un graper de sectionUpdate.go que usa #file:la implementacion de #file:GO_BASH_SCRIPTS.md" — entendido como: existe una implementación Go para la actualización de secciones (`sectionUpdate.go`) que el script shell invoca/usa. Podemos reutilizar esa implementación Go desde el CLI.

1) API: confirmado — se seguirá la recomendación por defecto: el paquete `badge` será "puro" y expondrá funciones para parsear badges y generar SVG (devuelve string/[]byte, count, warnings, error). El CLI / shim `badges.sh` será responsable de:
   - comprobar `.git` (si corresponde),
   - decidir la ruta de salida (docs/img vs custom),
   - escribir el archivo solo si cambió,
   - invocar la función de actualización de README reutilizando `sectionUpdate.go`.

2) README update: confirmado reutilizar la implementación Go existente (`sectionUpdate.go`) desde el CLI para insertar/actualizar la línea Markdown en el README. Si hace falta, añadiré una pequeña función exportada en `sectionUpdate.go` para facilitar esta integración.

3) Dependencias y heurística de anchura de texto: confirmado — se mantiene la aproximación simple (chars * font_size * 0.6) para evitar dependencias.

4) Requisito `.git`: confirmado — la comprobación de `.git` se mantendrá en el CLI/shim, no en la biblioteca, para que el paquete `badge` sea reutilizable fuera de repositorios git.

5) Formato de color: confirmado — sólo hex `#rrggbb` (se requiere el prefijo `#`).

6) Tests: confirmado — debemos soportar los tests existentes (`badges_test.go`). El CLI/shim reproducirá los mensajes y salidas que los tests esperan (mensajes de éxito/errores, rutas por defecto `docs/img/badges.svg`, comportamiento ante badges inválidos, etc.).

Detalles confirmados / decisiones tomadas

- Mensajes con colores ANSI: se mantendrán los colores en salida humana (como ahora). Los tests siguen buscando texto, por lo que no se alterarán las cadenas críticas.
- Import path / publicación: durante la integración en este repo el paquete residirá en `./badge` (import path local). Al extraer a un repo independiente ajustaremos el module path entonces.
- `badges.sh` shim: mantendremos `badges.sh` como shim que invoca el binario Go (si está disponible) y seguirá funcionando como script si el binario no existe — así mantenemos compatibilidad.
- Validación de color: se requiere el prefijo `#` (formato `#rrggbb`).

- Ruta de salida por defecto: `docs/img`. Esta ruta se mantendrá como valor por defecto pero será configurable tanto en la API (opción/argumento `output`) como en el CLI (flag `--output`).

# ISSUE: Convert badges.sh to Go package `badge`

Resumen

- Objetivo: extraer la lógica de `badges.sh` y convertirla en un paquete Go autónomo `badge` para poder publicarlo como repositorio independiente y reutilizarlo desde `devscripts` y otros proyectos.
- Decisión solicitada: revisa y aprueba este plan antes de la implementación.

Requisitos extraídos del prompt

- Crear el directorio `badge` en la raíz del repo.
- Implementar la lógica de generación de badges (parsing, validación, cálculo de anchuras, generación SVG, control de escritura de archivo si no cambia el contenido).
- Mantener el comportamiento observable por `badges_test.go` (salidas/errores, mensajes y formatos) o acordar adaptación de tests.
- Reducir dependencias a mínimo; preferir stdlib.
- Proponer API (funciones públicas), comportamiento y formatos de salida.

Contrato propuesto (mínimo)

- Entrada: lista de badges como strings `label:value:color` y opciones (output file path optional, readme path optional, generator info).
- Salida: (1) SVG como `string` o `[]byte`, (2) número de badges generados, (3) lista de warnings/errors.
- Modo de error: validaciones por badge devuelven warnings y siguen con badges válidos; si no hay badges válidos, error.
- Efectos secundarios opcionales: escribir archivo SVG y/o actualizar README (esto debería estar en el CLI o paquete?)

Preguntas / decisiones pendientes (necesito tu respuesta para empezar)

1) API: ¿prefieres que el paquete `badge` sólo genere y devuelva el SVG (pure function), dejando al consumidor la decisión de escribir archivos/actualizar README, o que el paquete tenga funciones de conveniencia que escriban el archivo y actualicen README?  
   - Opción A (recomendada): paquete puro que devuelve SVG y metadatos; separado un pequeño CLI/funciones auxiliares para escritura y README.


2) README update: `badges.sh` usa `sectionUpdate.sh` para insertar el markdown en README. ¿Quieres que la versión Go reimplemente esa lógica (parsing y actualización del README) en Go o que deje esa acción fuera del paquete y la gestione un wrapper/CLI?

    sectionUpdate.sh es un graper de sectionUpdate.go que usa #file:la implementacion de #file:GO_BASH_SCRIPTS.md


3) Dependencias y heurística de anchura de texto: el script actual calcula anchura con una fórmula aproximada (chars * font_size * 0.6). ¿Aceptar esa aproximación o usar una medida más precisa (p. ej. usar una librería de fuentes que mida texto, o incluir una tabla de anchos por carácter)?  
   - Mantener aproximación: simple y sin dependencias. Recomendado.  

4) Requisito .git: `badges.sh` falla si `.git` no existe. ¿Debemos mantener ese comportamiento en la biblioteca? (parece más apropiado en la CLI).  

5) Formato de color: el script acepta hex como `#rrggbb`. ¿Deseas soporte para nombres CSS o validaciones extras?  
no

6) Tests: ¿quieres que portemos los tests existentes (`badges_test.go`) para validar el nuevo paquete, o prefieres crear tests nuevos más orientados a la API del paquete?  
debemos soportar los tests existentes
Plan de trabajo propuesto (pasos)

1. Acordar respuestas a las preguntas anteriores.
2. Diseñar API pública del paquete `badge` (ej. `GenerateSVG(opts, badges) (svg string, count int, warnings []string, err error)`).
3. Implementar el paquete `badge` con funciones core: parseBadge, calcWidth, renderSVG, and optional helpers writeFile/updateReadme. Añadir docs y comentarios.
4. Añadir tests unitarios para el paquete (happy path + invalids + unchanged content). Adaptar `badges_test.go` a usar paquete si procede.

Siguientes pasos tras tu revisión

- Responde las preguntas marcadas como pendientes (1, 2 follow-up, 4, A-D). Si confirmas las opciones por defecto (paquete puro, CLI responsable de `.git`, reutilizar `sectionUpdate.go`, mantener colors, exigir `#` en hex), empezaré la implementación del paquete `badge`, añadiré tests y un CLI shim, luego ejecutaré `go test` y te reportaré los resultados.

Estado actual y siguiente acción

- Según tu indicación, las instrucciones de `docs/GO_BASHSCRIPTS_PROMPT.md` ya están implementadas. He actualizado este documento para fijar todas las decisiones necesarias para avanzar.
- Próximo paso: implementar (o enlazar) el paquete `badge` en `./badge`, portar la lógica de `badges.sh` a Go, añadir tests y crear `cmd/badgeCli` (y un `badges.sh` shim que invoque el binario). Ejecutaré `go test` y reportaré resultados.

Si necesitas que haga los cambios ahora, confirma y comienzo la implementación inmediata.

5. Ejecutar tests y arreglar fallos, añadir small README en `badge/README.md` y un ejemplo.
6. (Opcional) Añadir un pequeño CLI wrapper `cmd/badgeCli` o una función en `cmd/` que reproduzca `badges.sh` behaviour y que sea lo que `badges.sh` llame (podríamos renombrar o mantener el script como shim que llama al binario Go).

Alternativas y trade-offs

- Implementación completa en paquete (sin dependencias): más portable, fácil de probar. Recomendado.



# GoTest Badge Format Examples

## Opción 1: HTML con emojis y colores
```html
<div align="center">
    <span style="background-color: #28a745; color: white; padding: 2px 8px; border-radius: 4px;">
        ✅ Tests: Passing
    </span>
    <span style="background-color: #ffc107; color: black; padding: 2px 8px; border-radius: 4px;">
        ⚠️ Coverage: 85%
    </span>
    <span style="background-color: #28a745; color: white; padding: 2px 8px; border-radius: 4px;">
        ✅ Race: Clean
    </span>
    <span style="background-color: #28a745; color: white; padding: 2px 8px; border-radius: 4px;">
        ✅ Vet: Clean
    </span>
</div>
```

## Opción 2: Markdown simple con emojis
```markdown
**Tests:** ✅ Passing | **Coverage:** ⚠️ 85% | **Race:** ✅ Clean | **Vet:** ✅ Clean
```

## Opción 3: Badges de estilo GitHub (texto plano)
```markdown
![Tests](https://img.shields.io/badge/Tests-✅%20Passing-green)
![Coverage](https://img.shields.io/badge/Coverage-⚠️%2085%25-yellow)
![Race](https://img.shields.io/badge/Race-✅%20Clean-green)
![Vet](https://img.shields.io/badge/Vet-✅%20Clean-green)
```

## Opción 4: Tabla simple
```markdown
| Tests | Coverage | Race | Vet |
|-------|----------|------|-----|
| ✅ Passing | ⚠️ 85% | ✅ Clean | ✅ Clean |
```

## Estados posibles:

### Tests:
- ✅ Passing (verde)
- ❌ Failed (rojo)

### Coverage:
- ✅ 100% (verde)
- ⚠️ XX% (amarillo, cuando < 100%)
- ❌ No coverage (rojo, cuando no se puede obtener)

### Race Detection:
- ✅ Clean (verde)
- ❌ Race detected (rojo)

### Go Vet:
- ✅ Clean (verde)
- ❌ Issues found (rojo)

## Ejemplo de README.md con badges:

```markdown
# MyProject

**Tests:** ✅ Passing | **Coverage:** ⚠️ 85% | **Race:** ✅ Clean | **Vet:** ✅ Clean

Este es un proyecto de ejemplo...
```

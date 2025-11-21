# 5. Checklist de Validación de la Documentación - U-Go

| Campo                 | Valor                                                 |
| :-------------------- | :---------------------------------------------------- |
| **Proyecto**          | U-Go: Plataforma de Viajes Compartidos Universitarios |
| **Responsable**       | Grupo 12                                              |
| **Fecha de Creación** | 07/11/2025                                            |
| **Versión**           | 1.0                                                   |

---

## Checklist de Calidad para Revisión

Este checklist se utiliza para realizar una revisión por pares (o autoevaluación) de la calidad y compleción de los artefactos de requerimientos generados.

| Ítem                                                                                                                           | Estado    | Comentarios/Verificación                                                                                                               |
| :----------------------------------------------------------------------------------------------------------------------------- | :-------- | :------------------------------------------------------------------------------------------------------------------------------------- |
| **1. Requisitos (SRS)**                                                                                                        |           |                                                                                                                                        |
| [x] Los RF/RNF son **medibles, no ambiguos y verificables**.                                                                   | **HECHO** | Se usan métricas claras (ej. RNF-PER-01: < 500ms; RNF-PER-02: < 3s) y criterios verificables (ej. RF-01: activación por email).        |
| [x] El alcance (Scope) define claramente lo que **NO** incluye el sistema.                                                     | **HECHO** | El SRS 1.1.2 excluye explícitamente pasarelas de pago y _pooling_.                                                                     |
| [x] La terminología es consistente y se provee un glosario (acrónimos).                                                        | **HECHO** | SRS 1.1.3 define RF, RNF, Pasajero, Conductor, etc.                                                                                    |
| **2. Historias de Usuario (HU)**                                                                                               |           |                                                                                                                                        |
| [x] Las historias cumplen los criterios **INVEST** (Independientes, Negociables, con Valor, Estimables, Pequeñas, Testeables). | **HECHO** | Las HUs (HU-01 a HU-08) están enfocadas en valor (columna "Para") y son pequeñas.                                                      |
| [x] Se han documentado los criterios de aceptación en formato **Gherkin**.                                                     | **HECHO** | Todas las HUs en `HISTORIAS_USUARIO.md` tienen sus escenarios Gherkin (Dado, Cuando, Entonces).                                        |
| **3. Casos de Uso (UC)**                                                                                                       |           |                                                                                                                                        |
| [x] Los Casos de Uso incluyen **flujos alternos (Extensiones)** y reglas de negocio.                                           | **HECHO** | `CASOS_DE_USO.md` detalla Extensiones (ej. 4a en UC-04, 3a en UC-03) y Reglas de Negocio.                                              |
| [x] El diagrama (PlantUML) es coherente con los actores y UC especificados.                                                    | **HECHO** | Los 3 actores (Pasajero, Conductor, Admin) y los 8 UCs están correctamente mapeados.                                                   |
| **4. Trazabilidad y Coherencia**                                                                                               |           |                                                                                                                                        |
| [x] La **matriz de trazabilidad** cubre todos los RF críticos.                                                                 | **HECHO** | `TRAZABILIDAD.md` mapea RF $\leftrightarrow$ HU $\leftrightarrow$ UC $\leftrightarrow$ TC, asegurando que no hay requisitos huérfanos. |
| [x] Cada archivo incluye **Fecha, Versión y Responsable**.                                                                     | **HECHO** | Todos los archivos `.md` generados incluyen la cabecera de metadatos solicitada.                                                       |

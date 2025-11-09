# 📗 Documento 2: Historias de Usuario (HU) y Criterios Gherkin – _U-Go_

| Campo                 | Valor                                                 |
| :-------------------- | :---------------------------------------------------- |
| **Proyecto**          | U-Go: Plataforma de Viajes Compartidos Universitarios |
| **Responsable**       | Equipo 12 – Ingeniería de Sistemas                    |
| **Fecha de Creación** | 07/11/2025                                            |
| **Versión**           | 1.2 — Historias Ampliadas                             |

---

## 2.1 Tabla de Historias de Usuario

> Las siguientes historias definen las funcionalidades clave de **U-Go** desde la perspectiva del **valor al usuario**.

| ID        | Como (rol)     | Quiero (objetivo)                                                 | Para (beneficio)                                                                   | Prioridad | RF/RNF Relacionados     |
| :-------- | :------------- | :---------------------------------------------------------------- | :--------------------------------------------------------------------------------- | :-------- | :---------------------- |
| **HU-01** | **Estudiante** | registrarme únicamente con mi correo institucional                | garantizar que el servicio sea exclusivo y seguro para la comunidad universitaria. | Alta      | RF-01, RNF-SEC-01       |
| **HU-02** | **Pasajero**   | solicitar un viaje indicando mi origen y destino en el mapa       | obtener una ruta clara y un tiempo estimado de llegada (ETA).                      | Alta      | RF-02, RNF-PER-01       |
| **HU-03** | **Conductor**  | activar o desactivar mi disponibilidad                            | controlar cuándo deseo recibir solicitudes de viaje.                               | Alta      | RF-03                   |
| **HU-04** | **Pasajero**   | visualizar en tiempo real el vehículo asignado                    | conocer la ubicación exacta y el tiempo restante para mi recogida.                 | Media     | RF-05, RNF-PER-02       |
| **HU-05** | **Usuario**    | calificar al otro usuario (1 a 5 estrellas) al finalizar el viaje | mantener la calidad y la confianza dentro del servicio.                            | Media     | RF-06                   |
| **HU-06** | **Usuario**    | editar mi información básica de perfil (nombre, foto)             | mantener mis datos actualizados y ser fácilmente identificable.                    | Baja      | RF-Implícito, (MongoDB) |
| **HU-07** | **Conductor**  | registrar los datos y fotos de mi vehículo                        | permitir que mi vehículo sea verificado e identificado por los pasajeros.          | Alta      | RF-Implícito, (MongoDB) |
| **HU-08** | **Pasajero**   | cancelar un viaje solicitado antes de la recogida                 | evitar tomar un viaje que ya no necesito.                                          | Media     | RF-Implícito            |

---

## 2.2 Criterios de Aceptación (Formato Gherkin)

Los criterios de aceptación siguen el formato **Given–When–Then**, garantizando una validación clara, medible y automatizable.

---

### 🟦 HU-01: Registro con Correo Institucional

**Funcionalidad:** Verificación de Correo Institucional (RF-01)

**Escenario 1: Registro exitoso**

```
Dado que ingreso un correo con dominio institucional válido (ej. @uni.edu)
Y una contraseña que cumple las reglas de seguridad (RNF-SEC-01)
Cuando presiono "Registrar"
Entonces el sistema envía un enlace de activación al correo
Y la cuenta permanece inactiva hasta que hago clic en dicho enlace.
```

**Escenario 2: Falla por dominio incorrecto**

```
Dado que ingreso un correo con dominio no institucional (ej. @gmail.com)
Cuando presiono "Registrar"
Entonces el sistema muestra el mensaje "El correo debe ser institucional"
Y la cuenta no es creada.
```

---

### 🟩 HU-02: Solicitud de Viaje en el Mapa

**Funcionalidad:** Trazado de Ruta y ETA (RF-02)

**Escenario: Solicitud válida y rápida**

```
Dado que me encuentro en la pantalla de solicitud de viaje como Pasajero
Y la geolocalización está activa
Cuando selecciono mi Origen y Destino en el mapa
Entonces el sistema traza la ruta en menos de 2 segundos (RNF-PER-01)
Y muestra la distancia y el ETA del conductor más cercano.
```

---

### 🟧 HU-03: Gestión de Disponibilidad del Conductor

**Funcionalidad:** Estado de Conductor (RF-03)\*\*

**Escenario 1: Conductor se conecta**

```
Dado que soy un Conductor autenticado con estado “Desconectado”
Cuando activo el interruptor "Disponibilidad"
Entonces mi estado cambia a "Disponible"
Y empiezo a recibir solicitudes de viaje.
```

**Escenario 2: Falla al desconectarse con viaje activo**

```
Dado que soy un Conductor con estado "Disponible" y un viaje en curso
Cuando intento desactivar "Disponibilidad"
Entonces el sistema mantiene mi estado en "Disponible"
Y muestra el mensaje "No puedes desconectarte durante un viaje".
```

---

### 🟨 HU-04: Seguimiento en Tiempo Real

**Funcionalidad:** Mapa de Seguimiento (RF-05)\*\*

**Escenario: Actualización periódica de ubicación**

```
Dado que soy un Pasajero con viaje aceptado
Cuando observo la pantalla del mapa del viaje
Entonces veo el ícono del vehículo del Conductor
Y su posición se actualiza al menos cada 5 segundos.
```

---

### 🟦 HU-05: Calificación de Viaje

**Funcionalidad:** Sistema de Calificación (RF-06)\*\*

**Escenario 1: Calificación completa**

```
Dado que he finalizado un viaje (como Pasajero o Conductor)
Cuando abro la pantalla de calificación
Y selecciono “5 estrellas” con un comentario
Entonces la calificación y comentario se registran en el perfil del otro usuario.
```

**Escenario 2: Calificación rápida (sin comentario)**

```
Dado que he finalizado un viaje
Cuando califico con “3 estrellas” sin escribir comentario
Entonces la calificación se registra correctamente.
```

---

### 🟩 HU-06: Edición de Perfil de Usuario

**Funcionalidad:** Gestión de Perfil Básico

**Escenario: Actualización de foto de perfil**

```
Dado que estoy autenticado y en la sección “Mi Perfil”
Cuando selecciono “Cambiar foto de perfil”
Y subo una nueva imagen válida desde mi galería
Entonces la foto se actualiza visualmente en la app
Y se almacena en MongoDB/S3 asociada a mi ID de usuario.
```

---

### 🟧 HU-07: Registro de Vehículo (Conductor)

**Funcionalidad:** Gestión de Vehículos

**Escenario: Conductor registra un nuevo vehículo**

```
Dado que soy un Conductor autenticado y en la sección “Mis Vehículos”
Cuando ingreso la placa, modelo y color
Y subo una foto clara del vehículo
Y presiono “Guardar”
Entonces el vehículo se registra con estado “Pendiente de Aprobación”
Y las imágenes se almacenan en MongoDB/S3.
```

---

### 🟥 HU-08: Cancelación de Viaje

**Funcionalidad:** Cancelación de Solicitud

**Escenario: Cancelación antes de la recogida**

```
Dado que he solicitado un viaje y un Conductor ha sido asignado
Y el Conductor aún no ha iniciado el viaje
Cuando presiono “Cancelar viaje” y confirmo la acción
Entonces el viaje se marca como “Cancelado por pasajero”
Y el Conductor recibe una notificación de la cancelación.
```

---

## 2.3 Observaciones Finales

- Cada HU se vincula con sus **RF/RNF** en el documento SRS.
- Los escenarios Gherkin pueden transformarse en **tests automáticos** dentro del pipeline de QA.
- Se recomienda mantener esta tabla actualizada conforme se desarrollen nuevas iteraciones del producto.

---

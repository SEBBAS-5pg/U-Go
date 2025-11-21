# 📘 Documento de Especificación de Requisitos de Software (SRS)

### Basado en el Estándar IEEE 830

| Campo                 | Valor                                                 |
| :-------------------- | :---------------------------------------------------- |
| **Proyecto**          | U-Go: Plataforma de Viajes Compartidos Universitarios |
| **Responsable**       | Grupo 12                                              |
| **Fecha de Creación** | 07/11/2025                                            |
| **Versión**           | 1.0                                                   |

---

## 1. Introducción

### 1.1 Propósito

Este documento SRS tiene como objetivo definir de manera detallada los **requisitos funcionales y no funcionales** del sistema **U-Go**, una aplicación móvil de viajes compartidos diseñada exclusivamente para la comunidad universitaria.
Sirve de referencia para los equipos de **desarrollo**, **QA** y **stakeholders** al establecer el alcance, comportamiento esperado y criterios de aceptación del sistema.

### 1.2 Alcance

#### Funcionalidades Incluidas

U-Go permite:

- Registro e inicio de sesión con **verificación de correo institucional**.
- **Geolocalización** de usuarios.
- **Solicitud, asignación y seguimiento en tiempo real** de viajes.
- **Calificación mutua** entre conductores y pasajeros.

> La aplicación está dirigida exclusivamente a miembros de la comunidad universitaria (estudiantes y personal) con correo institucional verificado.

#### Funcionalidades Excluidas

U-Go **no incluye**:

- Integración con **pasarelas de pago** ni manejo de transacciones económicas internas.
- Funcionalidades de **viajes compartidos múltiples (_pooling_)** o **gestión de boletos**.
- **Gestión de publicidad** o monetización mediante anuncios.

### 1.3 Definiciones, Acrónimos y Abreviaturas

| Acrónimo      | Definición                                                                   |
| :------------ | :--------------------------------------------------------------------------- |
| **RF**        | Requisito Funcional                                                          |
| **RNF**       | Requisito No Funcional                                                       |
| **HU**        | Historia de Usuario                                                          |
| **UC**        | Caso de Uso                                                                  |
| **SRS**       | Software Requirements Specification                                          |
| **Pasajero**  | Usuario que solicita un viaje dentro de la app.                              |
| **Conductor** | Usuario verificado que ofrece un viaje (requiere verificación del vehículo). |

---

## 2. Descripción General

### 2.1 Perspectiva del Producto

U-Go es una aplicación móvil desarrollada en **Ionic React** (arquitectura híbrida/nativa), que consume una **API REST** implementada en **Go (Golang)**.
El backend emplea:

- **PostgreSQL** para el almacenamiento de datos relacionales y transaccionales.
- **MongoDB** para la gestión de imágenes y datos no estructurados.

### 2.2 Funciones del Producto (Resumen de Alto Nivel)

1. Registro y autenticación de usuarios con correo institucional.
2. Gestión de perfiles (conductor/pasajero).
3. Solicitud, asignación y seguimiento de viajes en tiempo real.
4. Calificación mutua al finalizar un viaje.

### 2.3 Características de los Usuarios

Los usuarios de U-Go son **miembros de la comunidad universitaria** (estudiantes, docentes o personal administrativo) con conocimientos básicos de tecnología móvil y acceso a un correo institucional válido.

### 2.4 Restricciones del Sistema

- **Tecnológicas:**

  - Uso obligatorio de **Ionic React**, **Go**, **PostgreSQL** y **MongoDB**.

- **Legales y de negocio:**

  - Acceso restringido a usuarios con correo institucional verificado.

- **Operacionales:**

  - Conectividad a Internet requerida para el uso completo de las funciones.

---

## 3. Requisitos Específicos

### 3.1 Requisitos Funcionales (RF)

| ID        | Descripción                                                                              | Prioridad | Criterio de Aceptación (Verificable)                                                                                              |
| :-------- | :--------------------------------------------------------------------------------------- | :-------- | :-------------------------------------------------------------------------------------------------------------------------------- |
| **RF-01** | El sistema permitirá el **registro de usuarios** con correo institucional y contraseña.  | Alta      | Dado un correo institucional válido y contraseña ≥ 8 caracteres, la cuenta se crea y se activa solo tras confirmación por correo. |
| **RF-02** | El usuario podrá **solicitar un viaje** indicando origen y destino en el mapa.           | Alta      | Al seleccionar origen/destino, el sistema traza la ruta y muestra el ETA en ≤ 2 segundos (ver RNF-PER-01).                        |
| **RF-03** | El conductor podrá **publicar su disponibilidad** para aceptar solicitudes de viaje.     | Alta      | Al activar el estado “Disponible”, el conductor puede recibir solicitudes; no las recibe si tiene un viaje activo.                |
| **RF-04** | El sistema **asignará automáticamente** un viaje al conductor más cercano disponible.    | Media     | La solicitud notifica al conductor más cercano según geolocalización.                                                             |
| **RF-05** | El sistema permitirá la **visualización en tiempo real** del vehículo asignado.          | Alta      | El pasajero visualiza la ubicación actualizada cada 5 segundos hasta finalizar el viaje.                                          |
| **RF-06** | Conductores y pasajeros podrán **calificar el viaje** (1 a 5 estrellas con comentarios). | Media     | Tras finalizar el viaje, ambos usuarios pueden asignar puntuación y comentario opcional.                                          |

---

### 3.2 Requisitos No Funcionales (RNF)

| ID             | Categoría                | Descripción (Medible y Verificable)                                                                                                       |
| :------------- | :----------------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| **RNF-SEC-01** | **Seguridad**            | Las contraseñas se almacenarán mediante funciones hash seguras (**bcrypt** o superior). Todo el tráfico deberá cifrarse con **TLS 1.2+**. |
| **RNF-PER-01** | **Rendimiento (Carga)**  | La latencia promedio de la API para la solicitud/asignación de viajes no superará **500 ms** en el p95.                                   |
| **RNF-PER-02** | **Rendimiento (Inicio)** | El tiempo de arranque de la aplicación no superará los **3 segundos** en dispositivos de gama media (p95).                                |
| **RNF-MAN-01** | **Mantenibilidad**       | El código backend (Go) deberá seguir una arquitectura modular y mantener una cobertura de pruebas unitarias mínima del **70%**.           |
| **RNF-DIS-01** | **Disponibilidad**       | La aplicación deberá ser compatible con las dos versiones más recientes de **Android** e **iOS** al momento del lanzamiento.              |
| **RNF-ACC-01** | **Accesibilidad**        | La interfaz cumplirá con el nivel **AA** de WCAG, y todos los elementos táctiles tendrán un tamaño mínimo de **44×44 píxeles**.           |

---

### 3.3 Trazabilidad y Referencias

Cada requisito (RF y RNF) será trazable hacia su correspondiente **HU (Historia de Usuario)** y **UC (Caso de Uso)** en la matriz de trazabilidad del proyecto (**TRAZABILIDAD.md**).
Las referencias citadas corresponden a fuentes internas del documento y normas IEEE para especificación de requisitos.

---

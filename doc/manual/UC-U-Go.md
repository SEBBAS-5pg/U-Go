# 3. Casos de Uso (Especificación) — U-Go 🚗🎓

| Campo | Valor |
| :--- | :--- |
| **Proyecto** | **U-Go: Plataforma de Viajes Compartidos Universitarios** |
| **Responsable** | Grupo 12 |
| **Fecha de Creación** | 07/11/2025 |
| **Versión** | 1.2 — *Especificación de 8 Casos de Uso* |

---

## 📑 Tabla de Contenidos

1. [UC-01: Gestionar Cuenta y Perfil](#-uc-01-gestionar-cuenta-y-perfil)  
2. [UC-02: Gestionar Vehículos](#-uc-02-gestionar-vehículos)  
3. [UC-03: Gestionar Disponibilidad](#-uc-03-gestionar-disponibilidad)  
4. [UC-04: Solicitar Viaje](#-uc-04-solicitar-viaje)  
5. [UC-05: Realizar Seguimiento de Viaje](#-uc-05-realizar-seguimiento-de-viaje)  
6. [UC-06: Cancelar Viaje](#-uc-06-cancelar-viaje)  
7. [UC-07: Calificar Viaje](#-uc-07-calificar-viaje)  
8. [UC-08: Aprobar Vehículo (Administrador)](#-uc-08-aprobar-vehículo-administrador)  
9. [Notas Generales](#-notas-generales)

---

## 3.2 Especificación de Casos de Uso

A continuación, se describen los casos de uso principales del sistema **U-Go**, según la arquitectura funcional definida y las historias de usuario relacionadas.

---

### 🧑‍💻 UC-01: Gestionar Cuenta y Perfil

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-01 |
| **Nombre** | Gestionar Cuenta y Perfil |
| **Actor Primario** | Pasajero, Conductor |
| **Descripción** | Permite el registro, autenticación y edición del perfil de usuario, incluyendo la actualización de nombre y foto (HU-01, HU-06). |
| **Precondiciones** | - Para registro: correo institucional válido.<br>- Para edición: usuario autenticado. |
| **Postcondiciones (Éxito)** | - Cuenta creada y activada (RF-01).<br>- Perfil actualizado y almacenado (foto en MongoDB/S3). |
| **Flujo Principal (Registro)** | 1. Usuario ingresa correo y clave institucional (HU-01).<br>2. El sistema valida el dominio del correo.<br>3. Se envía un correo de activación.<br>4. El usuario confirma su cuenta.<br>5. El sistema activa el perfil. |
| **Flujo Principal (Edición)** | 1. Usuario autenticado accede a “Mi Perfil”.<br>2. Modifica nombre o sube nueva foto (HU-06).<br>3. El sistema guarda los cambios.<br>4. Perfil actualizado correctamente. |
| **Extensiones (Flujo Alterno)** | **2a)** Correo no institucional → *El sistema muestra:* “El correo debe ser institucional”. |
| **RF/HU Relacionados** | RF-01 / HU-01, HU-06 |

---

### 🚘 UC-02: Gestionar Vehículos

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-02 |
| **Nombre** | Gestionar Vehículos |
| **Actor Primario** | Conductor |
| **Descripción** | Permite al conductor registrar los vehículos que usará para ofrecer viajes (HU-07). |
| **Precondiciones** | Usuario autenticado con rol **Conductor**. |
| **Postcondiciones (Éxito)** | Vehículo registrado con estado **Pendiente de Aprobación**. |
| **Flujo Principal** | 1. Conductor ingresa a “Mis Vehículos”.<br>2. Selecciona “Añadir Vehículo”.<br>3. Ingresa datos: placa, modelo y color.<br>4. Sube fotos del vehículo (frontal, placa).<br>5. El sistema guarda la información y marca el estado “Pendiente”. |
| **RF/HU Relacionados** | (RF implícito) / HU-07 |

---

### 🟢 UC-03: Gestionar Disponibilidad

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-03 |
| **Nombre** | Gestionar Disponibilidad |
| **Actor Primario** | Conductor |
| **Descripción** | Permite cambiar el estado del conductor (Disponible / No disponible) para recibir solicitudes de viaje (HU-03). |
| **Precondiciones** | Conductor autenticado con al menos un vehículo “Aprobado”. |
| **Postcondiciones (Éxito)** | Estado actualizado correctamente en la base de datos (RF-03). |
| **Flujo Principal** | 1. Conductor accede al panel principal.<br>2. Activa/desactiva el *toggle* de disponibilidad.<br>3. El sistema verifica que no haya viajes activos.<br>4. Actualiza el estado en PostgreSQL.<br>5. Confirma el nuevo estado (“Estás en línea”). |
| **Extensiones (Flujo Alterno)** | **3a)** Conductor tiene viaje activo → *El sistema muestra:* “No puedes desconectarte durante un viaje”. |
| **RF/HU Relacionados** | RF-03, RF-04 / HU-03 |

---

### 📍 UC-04: Solicitar Viaje

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-04 |
| **Nombre** | Solicitar Viaje |
| **Actor Primario** | Pasajero |
| **Descripción** | Permite solicitar un viaje indicando origen y destino en el mapa (HU-02). |
| **Precondiciones** | Pasajero autenticado (UC-01) y geolocalización activa. |
| **Postcondiciones (Éxito)** | Se asigna un conductor y se inicia el seguimiento (UC-05). |
| **Postcondiciones (Fallo)** | El sistema informa: “No hay conductores disponibles”. |
| **Flujo Principal** | 1. Pasajero selecciona **origen** y **destino** (RF-02).<br>2. El sistema calcula la ruta y ETA.<br>3. El pasajero confirma la solicitud.<br>4. Se buscan conductores disponibles (UC-03).<br>5. El más cercano acepta.<br>6. Se inicia el seguimiento del viaje (RF-05). |
| **Extensiones (Flujo Alterno)** | **4a)** Sin conductores disponibles → *Mensaje:* “Inténtelo de nuevo más tarde”. |
| **RF/HU Relacionados** | RF-02, RF-04, RF-05 / HU-02, HU-04, HU-08 |

---

### 🗺️ UC-05: Realizar Seguimiento de Viaje

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-05 |
| **Nombre** | Realizar Seguimiento de Viaje |
| **Actor Primario** | Pasajero, Conductor |
| **Descripción** | Permite visualizar la ubicación y ruta del viaje en tiempo real (HU-04). |
| **Precondiciones** | El viaje está en estado **Aceptado**. |
| **Postcondiciones (Éxito)** | Viaje finalizado correctamente; se habilita UC-07 (Calificación). |
| **Flujo Principal** | 1. Se muestra la ubicación del conductor en el mapa (RF-05).<br>2. Se actualiza la posición en tiempo real.<br>3. El conductor llega al destino.<br>4. Presiona “Finalizar viaje”.<br>5. El sistema marca el viaje como **Finalizado**.<br>6. Redirige a la pantalla de calificación (UC-07). |
| **RF/HU Relacionados** | RF-05 / HU-04 |

---

### ❌ UC-06: Cancelar Viaje

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-06 |
| **Nombre** | Cancelar Viaje |
| **Actor Primario** | Pasajero |
| **Descripción** | Permite cancelar un viaje en estado “Solicitado” o “Aceptado” (HU-08). |
| **Precondiciones** | Viaje no iniciado por el conductor. |
| **Postcondiciones (Éxito)** | Viaje actualizado como “Cancelado por Pasajero” y conductor notificado. |
| **Flujo Principal** | 1. Pasajero presiona “Cancelar Viaje”.<br>2. El sistema solicita confirmación.<br>3. Pasajero confirma la cancelación.<br>4. El sistema cambia el estado a “Cancelado”.<br>5. Notifica al conductor (si aplica). |
| **RF/HU Relacionados** | (RF implícito) / HU-08 |

---

### ⭐ UC-07: Calificar Viaje

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-07 |
| **Nombre** | Calificar Viaje |
| **Actor Primario** | Pasajero, Conductor |
| **Descripción** | Permite calificar al otro usuario con estrellas y comentario (HU-05). |
| **Precondiciones** | El viaje está en estado **Finalizado**. |
| **Postcondiciones (Éxito)** | Calificación almacenada en la base de datos PostgreSQL. |
| **Flujo Principal** | 1. Se muestra pantalla de calificación.<br>2. Usuario selecciona de 1 a 5 estrellas.<br>3. (Opcional) Agrega comentario.<br>4. Presiona “Enviar”.<br>5. Sistema guarda la calificación. |
| **RF/HU Relacionados** | RF-06 / HU-05 |

---

### 🛠️ UC-08: Aprobar Vehículo (Administrador)

| Campo | Descripción |
| :--- | :--- |
| **ID** | UC-08 |
| **Nombre** | Aprobar Vehículo |
| **Actor Primario** | Administrador |
| **Descripción** | Caso de uso administrativo para verificar y aprobar los vehículos registrados. |
| **Precondiciones** | Administrador autenticado en el panel web. |
| **Postcondiciones (Éxito)** | Vehículo aprobado; conductor notificado. |
| **Postcondiciones (Fallo)** | Vehículo rechazado con motivo enviado al conductor. |
| **Flujo Principal** | 1. Administrador abre “Vehículos Pendientes”.<br>2. Selecciona un vehículo.<br>3. Revisa datos e imágenes (HU-07).<br>4. Presiona “Aprobar”.<br>5. El sistema cambia el estado y notifica al conductor. |
| **Extensiones (Flujo Alterno)** | **4a)** Datos incorrectos → *Administrador rechaza y detalla el motivo* → Conductor recibe notificación para corrección. |
| **RF/HU Relacionados** | (RF implícito de Admin) / HU-07 |

---

## 📘 Notas Generales

- Las imágenes y documentos del vehículo y perfil se almacenan en **MongoDB/S3**, mientras que los datos transaccionales (viajes, calificaciones, estados) se registran en **PostgreSQL**.  
- Los casos **UC-01 a UC-07** corresponden al *Front/Back principal (Pasajero–Conductor)*, mientras que **UC-08** pertenece al *Módulo Administrativo (Back Office)*.  
- Este documento sigue la trazabilidad **RF ↔ HU ↔ UC**, complementado con el archivo [`TRAZABILIDAD.md`](./TRAZABILIDAD.md).

---

> 🧭 **Documento oficial de Especificación de Casos de Uso — Proyecto U-Go (2025).**

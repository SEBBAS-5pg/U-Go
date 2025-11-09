# 4. Matriz de Trazabilidad (REQ <-> HU <-> UC <-> Test) - U-Go

| Campo                 | Valor                                                 |
| :-------------------- | :---------------------------------------------------- |
| **Proyecto**          | U-Go: Plataforma de Viajes Compartidos Universitarios |
| **Responsable**       | Grupo 12                                              |
| **Fecha de Creación** | 07/11/2025                                            |
| **Versión**           | 1.0                                                   |

---

## Matriz de Trazabilidad de Requerimientos

Esta matriz asegura que cada Requisito (REQ) esté cubierto por al menos una Historia de Usuario (HU), un Caso de Uso (UC) y un Caso de Prueba (TC), garantizando que el producto final cumpla con lo especificado en el SRS.

| Requisito (RF)               | Historias de Usuario (HU)        | Casos de Uso (UC)                                        | Pruebas (Test Cases - TC)                                                         |
| :--------------------------- | :------------------------------- | :------------------------------------------------------- | :-------------------------------------------------------------------------------- |
| **RF-01** (Registro Inst.)   | HU-01 (Registro institucional)   | UC-01 (Gestionar Cuenta)                                 | TC-01 (Registro válido); TC-02 (Email NO institucional)                           |
| **RF-02** (Solicitud Mapa)   | HU-02 (Solicitar viaje mapa)     | UC-04 (Solicitar Viaje)                                  | TC-10 (Trazado ruta < 2s); TC-11 (Origen/Destino fuera de límites)                |
| **RF-03** (Disponibilidad)   | HU-03 (Gestionar Disponibilidad) | UC-03 (Gestionar Disponibilidad)                         | TC-20 (Toggle ON/OFF); TC-21 (Intento desactivar con viaje activo)                |
| **RF-04** (Asignación Auto)  | HU-02 (Solicitar viaje mapa)     | UC-04 (Solicitar Viaje)                                  | TC-30 (Asignación a conductor más cercano); TC-31 (Falla si no hay conductores)   |
| **RF-05** (Seguimiento T.R.) | HU-04 (Ver en tiempo real)       | UC-05 (Realizar Seguimiento)                             | TC-40 (Actualización de posición cada 5s); TC-41 (Fin de seguimiento)             |
| **RF-06** (Calificación)     | HU-05 (Calificar usuario)        | UC-07 (Calificar Viaje)                                  | TC-50 (Calificación 5 estrellas); TC-51 (Calificación sin comentario)             |
| **(RF-Impl) / RNF**          | HU-06 (Editar Perfil)            | UC-01 (Gestionar Cuenta)                                 | TC-60 (Subir foto de perfil a MongoDB); TC-61 (Cambiar nombre)                    |
| **(RF-Impl) / RNF**          | HU-07 (Registrar Vehículo)       | UC-02 (Gestionar Vehículos);<br>UC-08 (Aprobar Vehículo) | TC-70 (Registro vehículo pendiente); TC-71 (Admin aprueba); TC-72 (Admin rechaza) |
| **(RF-Impl) / RNF**          | HU-08 (Cancelar Viaje)           | UC-06 (Cancelar Viaje)                                   | TC-80 (Cancelar antes de aceptar); TC-81 (Cancelar después de aceptar)            |

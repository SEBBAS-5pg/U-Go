import React, { useEffect, useRef, useState } from 'react';
import { IonCard, IonCardContent, IonNote, IonIcon } from '@ionic/react';
import { locate, car } from 'ionicons/icons';

// Coordenadas iniciales de prueba (por ejemplo, Barcelona, España)
const DEFAULT_CENTER: [number, number] = [41.3851, 2.1734];
const DEFAULT_ZOOM = 13;

// 1. Exportamos la interfaz para que otros componentes puedan usarla si es necesario
export interface MapProps {
  center: [number, number];
  zoom: number;
  driverLocation?: [number, number] | null;
  passengerLocation?: [number, number] | null;
}

// 2. Componente interno que dibuja el mapa (MapRenderer)
export const MapRenderer: React.FC<MapProps> = ({ 
    center, 
    zoom, 
    driverLocation, 
    passengerLocation 
}) => {
    const mapContainerRef = useRef<HTMLDivElement>(null);
    const [isMapLoaded, setIsMapLoaded] = useState(false);

    useEffect(() => {
        // Simulamos un pequeño delay de carga
        const timer = setTimeout(() => {
            setIsMapLoaded(true);
        }, 500); 
        return () => clearTimeout(timer);
    }, [center, zoom]);

    // Renderiza marcador del conductor
    const renderDriverMarker = () => {
        if (!driverLocation) return null;
        return (
            <div className="map-marker driver-marker" key="driver">
                <IonIcon icon={car} color="danger" />
                <IonNote color="danger">Conductor</IonNote>
            </div>
        );
    };

    // Renderiza marcador del pasajero
    const renderPassengerMarker = () => {
        if (!passengerLocation) return null;
        return (
            <div className="map-marker passenger-marker" key="passenger">
                <IonIcon icon={locate} color="primary" />
                <IonNote color="primary">Tú</IonNote>
            </div>
        );
    };

    return (
        <div className="map-container" ref={mapContainerRef}>
            {!isMapLoaded && (
                <div className="map-loading">
                    Cargando Mapa Interactivo...
                </div>
            )}
            {isMapLoaded && (
                <div className="map-simulation">
                    <div className="map-bg"></div>
                    <div className="marker-layer">
                        {renderDriverMarker()}
                        {renderPassengerMarker()}
                    </div>
                    <IonNote className="map-legend" color="medium">
                       Zoom: {zoom}, Centro: {center[0].toFixed(4)}, {center[1].toFixed(4)}
                    </IonNote>
                </div>
            )}
        </div>
    );
};

/**
 * 3. Componente Contenedor (TripMap) - AQUÍ ESTABA EL ERROR
 * * Solución: Usamos Partial<MapProps> para indicarle a TypeScript que 
 * este componente PUEDE recibir las propiedades de MapProps, 
 * pero que son opcionales (porque si no vienen, usamos la simulación).
 */
const TripMap: React.FC<Partial<MapProps>> = (props) => {
    
    // Si nos pasan 'center', asumimos que vienen datos reales desde Tab2
    if (props.center && props.zoom) {
        // Forzamos el tipado porque ya verificamos que existen
        return <MapRenderer {...(props as MapProps)} />;
    }

    // --- MODO SIMULACIÓN (Si no hay props) ---
    // Esto se usa si pones <TripMap /> sin argumentos en algún lado
    const [mockDriverLat, setMockDriverLat] = useState(DEFAULT_CENTER[0] + 0.005);
    const mockDriverLng = DEFAULT_CENTER[1] + 0.002;

    useEffect(() => {
        const interval = setInterval(() => {
            setMockDriverLat(prevLat => prevLat + 0.0001); 
        }, 5000);
        return () => clearInterval(interval);
    }, []);

    return (
        <MapRenderer
            center={DEFAULT_CENTER}
            zoom={DEFAULT_ZOOM}
            driverLocation={[mockDriverLat, mockDriverLng]}
            passengerLocation={DEFAULT_CENTER}
        />
    );
};

export default TripMap;
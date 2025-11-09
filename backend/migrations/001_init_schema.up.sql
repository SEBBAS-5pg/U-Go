-- Migración 001: Esquema Inicial de U-Go

-- (Borramos en orden inverso a la creación)

DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS trips;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS users;

-- Borramos los ENUMs
DROP TYPE IF EXISTS trip_status;
DROP TYPE IF EXISTS vehicle_status;
DROP TYPE IF EXISTS driver_status;

-- Crear ENUMs (Tipos de datos personalizados)
CREATE TYPE driver_status AS ENUM (
    'offline',
    'online',
    'en_viaje'
);

CREATE TYPE vehicle_status AS ENUM (
    'pendiente',
    'aprobado',
    'rechazado'
);

CREATE TYPE trip_status AS ENUM (
    'solicitado',
    'aceptado',
    'en_curso',
    'finalizado',
    'cancelado'
);

-- Tabla 1: users (Basada en tu MER)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    activation_token VARCHAR(100),
    is_driver BOOLEAN DEFAULT false,
    driver_status driver_status DEFAULT 'offline',
    profile_image_url TEXT, -- De Mongo
    average_rating DECIMAL(3, 2) DEFAULT 0.0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Tabla 2: vehicles (Basada en tu MER)
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conductor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plate VARCHAR(20) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    color VARCHAR(50) NOT NULL,
    status vehicle_status NOT NULL DEFAULT 'pendiente',
    vehicle_image_url TEXT, -- De Mongo
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Tabla 3: trips (Basada en tu MER)
CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pasajero_id UUID NOT NULL REFERENCES users(id),
    conductor_id UUID REFERENCES users(id),
    vehicle_id UUID REFERENCES vehicles(id),
    status trip_status NOT NULL DEFAULT 'solicitado',
    origin_lat DECIMAL(10, 8) NOT NULL,
    origin_lng DECIMAL(11, 8) NOT NULL,
    origin_name VARCHAR(255),
    destination_lat DECIMAL(10, 8) NOT NULL,
    destination_lng DECIMAL(11, 8) NOT NULL,
    destination_name VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

-- Tabla 4: ratings (Basada en tu MER)
CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trips(id),
    rater_id UUID NOT NULL REFERENCES users(id),
    rated_id UUID NOT NULL REFERENCES users(id),
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(trip_id, rater_id)
);
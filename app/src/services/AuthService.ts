import axios, { isAxiosError } from "axios";

interface RegisterPayload {
  email: string;
  fullName: string;
  password: string;
}

interface LoginPayload {
  email: string;
  password: string;
}

interface LoginResponse {
  token: string;
}

// lee la URL base del archivo .env
const BASE_URL = import.meta.env.VITE_API_BASE_URL;

// Esta es la "llave" con la que se guarda el token
const TOKEN_KEY = "ugo-auth-token";

// Llama al endpoint de registro usando AXIOS

async function register(payload: RegisterPayload) {
  try {
    const response = await axios.post(`${BASE_URL}/auth/register`, payload);
    return response.data; // Devuelve el usuario creado
  } catch (error) {
    if (isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || "Registration error");
    }
    throw new Error("Network or unknown error");
  }
}

// Llama al endpoint de Login T-08 usando AXIOS
async function login(payload: LoginPayload): Promise<string> {
  try {
    const response = await axios.post<LoginResponse>(
      `${BASE_URL}/auth/login`,
      payload
    );

    const token = response.data.token;

    // Guarda el token en el almacenamiento local
    localStorage.setItem(TOKEN_KEY, token);

    return token;
  } catch (error) {
    if (isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || "Error en el login");
    }
    throw new Error("Error de red o desconocido");
  }
}

// Borra el token (Cerrar Sesion)
function logout(){
    localStorage.removeItem(TOKEN_KEY)
}

function getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
}

// Convierte el archivo a un "modulo"
export const AuthService = {
    register,
    login,
    logout,
    getToken
};
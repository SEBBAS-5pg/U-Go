import { IonButton, IonContent, IonIcon, IonInput, IonLoading, IonPage, IonText, useIonRouter } from "@ionic/react";
import { useState } from "react";

// Importa conector
import { AuthService } from "../services/AuthService";
import { logInOutline, personCircleOutline } from "ionicons/icons";

const LoginPage: React.FC = () => {
  // Estados de React
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false); // para el spiner

  const router = useIonRouter(); // Hook de ionic para navegar

  // se llama al presionar el boton "entrar"
  const handleLogin = async () => {
    setError(""); // Limpia errores antiguos
    setIsLoading(true); //Muestra el spinner

    // Validacion simple
    if (!email || !password) {
      setError("Please enter your email and password.");
      setIsLoading(false); // oculta el spinner
      return;
    }

    try {
      // llamamos al authservice de axios
      await AuthService.login({ email, password });

      setIsLoading(false); // Oculta el spinner

      // redirige el usuario a la app principal
      // replace borra la pagina de login del historial

      router.push("/tabs", "forward", "replace");
    } catch (err: any) {
      // error
      setIsLoading(false); // oculta el spinner
      // err.message es el error que se creo en el AuthService
      setError(err.message || "Unknown error");
    }
  };

  return (
    <IonPage>
      {/* Usamos 'bg-background' (el gris claro que definimos)
        y 'flex' de Tailwind para centrar todo.
      */}
      <IonContent fullscreen className="ion-padding bg-background">
        <div className="flex flex-col items-center justify-center min-h-full">
          {/* Tarjeta del Formulario */}
          <div className="w-full max-w-md px-8 py-10 bg-white rounded-lg shadow-xl">
            {/* Icono de Cabecera */}
            <div className="flex justify-center mb-6">
              <IonIcon
                icon={personCircleOutline}
                className="text-7xl text-primary" // <-- ¡Clase de Tema!
              />
            </div>

            <h2 className="text-2xl font-bold text-center text-text-primary mb-8">
              Iniciar Sesión U-Go
            </h2>

            {/* Formulario */}
            <form
              onSubmit={(e) => {
                e.preventDefault();
                handleLogin();
              }}
            >
              <IonInput
                label="Correo Institucional"
                labelPlacement="floating"
                fill="outline"
                type="email"
                value={email}
                onIonInput={(e) => setEmail(e.detail.value!)}
                className="mb-4"
                placeholder="tu@uni.edu"
              />

              <IonInput
                label="Contraseña"
                labelPlacement="floating"
                fill="outline"
                type="password"
                value={password}
                onIonInput={(e) => setPassword(e.detail.value!)}
                className="mb-6"
                placeholder="Tu contraseña"
              />

              {/* Mensaje de Error (si existe) */}
              {error && (
                <IonText color="danger" className="text-center block mb-4">
                  <p>{error}</p>
                </IonText>
              )}

              {/* Botón de Entrar */}
              <IonButton
                type="submit"
                expand="block"
                className="mb-4"
                disabled={isLoading} // Se deshabilita mientras carga
              >
                <IonIcon icon={logInOutline} slot="start" />
                {isLoading ? "Ingresando..." : "Entrar"}
              </IonButton>

              {/* Spinner de Carga (oculto) */}
              <IonLoading isOpen={isLoading} message={"Por favor espera..."} />
            </form>

            {/* TODO: En un paso futuro, haremos que este link
              lleve a la página de registro.
            */}
            <div className="text-center mt-6">
              <a href="#" className="text-sm text-primary hover:underline">
                ¿No tienes cuenta? Regístrate
              </a>
            </div>
          </div>
        </div>
      </IonContent>
    </IonPage>
  );
};

export default LoginPage;
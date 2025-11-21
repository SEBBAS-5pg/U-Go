import { IonButton, IonIcon, IonInput, IonLoading, IonText } from "@ionic/react";
import { logInOutline, personCircleOutline } from "ionicons/icons";
import { useState } from "react";

interface AuthFormProps {
  formType: "login" | "register"; // para saber si mostrar el campo "Nombre"
  headerText: string; // Iniciar Sesion U-Go
  buttonText: string; // Entrar
  linkTo: string; // register
  linkText: string; // ¿no tienes cuenta?

  // funcion que se ejecutara al enviar el formulario
  // promesa para que se pueda manejar el 'isLoading'

  onSubmit: (data: {
    email: string;
    password: string;
    fullName?: string; // Opcional, solo para registro
  }) => Promise<void>;
}

const AuthForm: React.FC<AuthFormProps> = ({
  formType,

  headerText,
  buttonText,
  linkTo,
  linkText,
  onSubmit,
}) => {
  // Estados Internos
  // formulario maneja su propio estado para los campos de texto
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState(""); // Solo se usará si formType == 'register'

  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  // Se llama al presionar el boton de submit

  const handleSubmitInternal = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    // Validacion simple
    if (!email || !password || (formType === 'register' && !fullName)) {
      setError('Por favor, completa todos los campos.');
      setIsLoading(false);
      return;
    }
    
    try {
      // Llama a la función "onSubmit" que nos pasó la página padre
      // (sea handleLogin o handleRegister)
      await onSubmit({ email, password, fullName });
      
      setIsLoading(false);
      
    } catch (err: any) {
      // Si la promesa (onSubmit) falla, mostramos el error
      setIsLoading(false);
      setError(err.message || 'Error desconocido');
    }
  };

  return (
    // Copiamos la UI de Tailwind que ya habíamos hecho
    <div className="flex flex-col items-center justify-center min-h-full">
      <div className="w-full max-w-md px-8 py-10 bg-white rounded-lg shadow-xl">

        <div className="flex justify-center mb-6">
          <IonIcon 
            icon={personCircleOutline} 
            className="text-7xl text-primary" 
          />
        </div>

        <h2 className="text-2xl font-bold text-center text-text-primary mb-8">
          {headerText} {/* <-- Prop! */}
        </h2>

        <form onSubmit={handleSubmitInternal}>
          
          {/* --- Campo Condicional --- */}
          {/* Solo muestra "Nombre Completo" si es el formulario de Registro */}
          {formType === 'register' && (
            <IonInput
              label="Nombre Completo"
              labelPlacement="floating"
              fill="outline"
              type="text"
              value={fullName}
              onIonInput={(e) => setFullName(e.detail.value!)}
              className="mb-4"
              placeholder="Tu Nombre Completo"
            />
          )}

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

          {error && (
            <IonText color="danger" className="text-center block mb-4">
              <p>{error}</p>
            </IonText>
          )}

          <IonButton
            type="submit"
            expand="block"
            className="mb-4"
            disabled={isLoading}
          >
            <IonIcon icon={logInOutline} slot="start" />
            {isLoading ? 'Cargando...' : buttonText} {/* <-- Prop! */}
          </IonButton>

          <IonLoading isOpen={isLoading} message={"Por favor espera..."} />
        </form>

        <div className="text-center mt-6">
          <a href={linkTo} className="text-sm text-primary hover:underline">
            {linkText} {/* <-- Prop! */}
          </a>
        </div>

      </div>
    </div>
  );
};

export default AuthForm;
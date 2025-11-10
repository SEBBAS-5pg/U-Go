import { IonContent, IonPage, useIonRouter, useIonToast } from "@ionic/react";
import { AuthService } from "../services/AuthService";
import AuthForm from "../components/AuthForm";

const RegisterPage: React.FC = () => {
  const router = useIonRouter();
  const [presentToast] = useIonToast();
  // esta funcion maneja la logica de ¿que? hacer cuando el formulario de registro se envia
  const handleRegister = async (formData: {
    email: string;
    password: string;
    fullName?: string;
  }) => {
    // 'AuthForm' validó que los campos no estan vacios
    // Lanza un error si falla

    // Llamamos al AuthService (Backend HU-01)
    await AuthService.register({
      email: formData.email,
      password: formData.password,
      fullName: formData.fullName || "", // Pasamos el nombre
    });

    // Mostramos una notificacion
    presentToast({
      message: "Registration successful! Please log in. ✅",
      duration: 2500,
      color: "succes",

      // se ejecuta despues de que el toast desaparece
      onDidDismiss: () => {
        //Redirige el usuario a la pagina de Login
        router.push("/login", "back", "replace");
      },
    });
  };

  return (
    <IonPage>
      <IonContent fullscreen className="ion-padding bg-background">
        {/* --- 
          ¡AQUÍ LA MAGIA! 
          Renderizamos el MISMO componente, pero con props
          diferentes para el REGISTRO.
        --- */}
        <AuthForm
          formType="register" // <-- ¡LA CLAVE! (Muestra el campo "Nombre")
          headerText="Crear Cuenta U-Go"
          buttonText="Registrarse"
          linkTo="/login"
          linkText="¿Ya tienes cuenta? Inicia sesión"
          onSubmit={handleRegister} // Le pasamos nuestra lógica de registro
        />
      </IonContent>
    </IonPage>
  );
};

export default RegisterPage;

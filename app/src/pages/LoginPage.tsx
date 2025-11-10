// Ruta: /app/src/pages/LoginPage.tsx

import React from 'react';
import {
  IonPage,
  IonContent,
  useIonRouter,
} from '@ionic/react';

// 1. Importamos nuestro "conector"
import { AuthService } from '../services/AuthService';
// 2. Importamos nuestro NUEVO componente de UI
import AuthForm from '../components/AuthForm';

const LoginPage: React.FC = () => {
  const router = useIonRouter(); // Hook de Ionic para navegar

  /**
   * Esta función maneja la lógica de "QUÉ" hacer cuando
   * el formulario de Login se envía.
   */
  const handleLogin = async (formData: { email: string, password: string }) => {
    // El 'AuthForm' ya validó que los campos no están vacíos.
    // Lanzará un error si el login falla, que será
    // atrapado y mostrado por el AuthForm.
    
    // 1. Llamamos al AuthService
    await AuthService.login({ 
      email: formData.email, 
      password: formData.password 
    });

    // 2. ¡Éxito! Redirigimos al usuario a la app principal
    router.push('/tabs', 'forward', 'replace');
  };

  return (
    <IonPage>
      <IonContent fullscreen className="ion-padding bg-background">
        
        {/* --- 
          ¡AQUÍ LA MAGIA! 
          Renderizamos el componente reutilizable y le pasamos las 
          props específicas para el LOGIN.
        --- */}
        <AuthForm
          formType="login"
          headerText="Iniciar Sesión U-Go"
          buttonText="Entrar"
          linkTo="/register" // (Aún no existe, pero lo crearemos)
          linkText="¿No tienes cuenta? Regístrate"
          onSubmit={handleLogin} // Le pasamos nuestra lógica de login
        />
        
      </IonContent>
    </IonPage>
  );
};

export default LoginPage;
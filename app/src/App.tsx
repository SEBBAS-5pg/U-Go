import { Redirect, Route } from 'react-router-dom';
import {
  IonApp,
  IonIcon,
  IonLabel,
  IonRouterOutlet,
  IonTabBar,
  IonTabButton,
  IonTabs,
  setupIonicReact
} from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';

// Iconos que ya estabas usando
import { ellipse, square, triangle } from 'ionicons/icons';

// Páginas de Tabs que ya tenías
import Tab1 from './pages/Tab1';
import Tab2 from './pages/Tab2';
import Tab3 from './pages/Tab3';

// --- (1) IMPORTAR LA NUEVA PÁGINA ---
import LoginPage from './pages/LoginPage'; // ¡Importante!

/* ... (Todos los imports de CSS se quedan igual) ... */
/* Core CSS required for Ionic components to work properly */
import '@ionic/react/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';

/* Optional CSS utils that can be commented out */
import '@ionic/react/css/padding.css';
import '@ionic/react/css/float-elements.css';
import '@ionic/react/css/text-alignment.css';
import '@ionic/react/css/text-transformation.css';
import '@ionic/react/css/flex-utils.css';
import '@ionic/react/css/display.css';

/**
 * Ionic Dark Mode
 * -----------------------------------------------------
 */
/*
import '@ionic/react/css/palettes/dark.system.css';
*/

/* Theme variables */
import './theme/variables.css';

setupIonicReact();

/* --- (2) EXTRAER LAS TABS A SU PROPIO COMPONENTE ---
   El contenido de tu App.tsx (las pestañas) ahora vivirá aquí.
   Lo ponemos en el mismo archivo para no crear más archivos.
*/
const MainTabs: React.FC = () => {
  return (
    <IonTabs>
      <IonRouterOutlet>
        <Route exact path="/tabs/tab1">
          <Tab1 />
        </Route>
        <Route exact path="/tabs/tab2">
          <Tab2 />
        </Route>
        <Route path="/tabs/tab3">
          <Tab3 />
        </Route>
        {/* Redirección por defecto DENTRO de las tabs */}
        <Route exact path="/tabs">
          <Redirect to="/tabs/tab1" />
        </Route>
      </IonRouterOutlet>
      
      {/* Esta es la barra de pestañas de abajo */}
      <IonTabBar slot="bottom">
        <IonTabButton tab="tab1" href="/tabs/tab1">
          <IonIcon aria-hidden="true" icon={triangle} />
          <IonLabel>Tab 1</IonLabel>
        </IonTabButton>
        <IonTabButton tab="tab2" href="/tabs/tab2">
          <IonIcon aria-hidden="true" icon={ellipse} />
          <IonLabel>Tab 2</IonLabel>
        </IonTabButton>
        <IonTabButton tab="tab3" href="/tabs/tab3">
          <IonIcon aria-hidden="true" icon={square} />
          <IonLabel>Tab 3</IonLabel>
        </IonTabButton>
      </IonTabBar>
    </IonTabs>
  );
};

/* --- (3) MODIFICAR EL ROUTER PRINCIPAL (App) ---
   Este es el nuevo componente App.
   Ahora es un "router" principal que decide si mostrar
   el Login o las Tabs.
*/
const App: React.FC = () => (
  <IonApp>
    <IonReactRouter>
      {/* El IonRouterOutlet principal maneja todas las rutas */}
      <IonRouterOutlet>
        
        {/* Ruta para el Login */}
        <Route exact path="/login">
          <LoginPage />
        </Route>
        
        {/* Ruta para el resto de la app (las Tabs) */}
        {/* Si la URL empieza con /tabs, carga el componente MainTabs */}
        <Route path="/tabs" component={MainTabs} />
        
        {/* Redirección por defecto de la App */}
        {/* Ahora, / (la raíz) te manda a /login, no a /tab1 */}
        <Route exact path="/">
          <Redirect to="/login" />
        </Route>
        
      </IonRouterOutlet>
    </IonReactRouter>
  </IonApp>
);

export default App;
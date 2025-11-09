/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    // 1. Escanea todos los archivos TSX y JS en la carpeta src
    "./src/**/*.{js,jsx,ts,tsx}",
    
    // 2. Escanea el index.html
    "./public/index.html",
  ],
  theme: {
    extend: {
      // Aquí conectamos la paleta de Ionic con las clases de Tailwind
      colors: {
        'primary': 'var(--ion-color-primary)',
        'secondary': 'var(--ion-color-secondary)',
        'tertiary': 'var(--ion-color-tertiary)',
        'success': 'var(--ion-color-success)',
        'warning': 'var(--ion-color-warning)',
        'danger': 'var(--ion-color-danger)',
        
        // Renombramos los colores de Ionic para que sean más intuitivos en Tailwind
        'background': 'var(--ion-background-color)', // Fondo principal de la app
        'text-primary': 'var(--ion-text-color)',    // Color de texto principal
        'text-medium': 'var(--ion-color-medium)',  // Color de texto secundario (gris)
        'text-light': 'var(--ion-color-light)',    // Color de texto claro
      }
    },
  },
  plugins: [],
}
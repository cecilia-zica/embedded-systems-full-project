// Vite = o "servidor de dev + empacotador" do mundo React, equivalente ao que
// o `flutter run` + `flutter build web` fazem juntos pro Flutter: serve o app
// com hot-reload em dev (`npm run dev`) e gera os arquivos estáticos finais
// em dev (`npm run build`), prontos pra hospedar em qualquer lugar.
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// o plugin `react()` ensina o Vite a entender JSX (a sintaxe <div>...</div>
// dentro do .jsx) e a habilitar Fast Refresh (hot-reload preservando estado)
export default defineConfig({
  plugins: [react()],
})

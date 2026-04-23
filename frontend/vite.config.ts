import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// export default defineConfig({
// 	plugins: [sveltekit()]
// });

export default defineConfig(({ mode }) => {
  const BACKEND_PORT = 8080;
  const FRONTEND_PORT = 5555;

  return {
    plugins: [sveltekit()],
    base: '/',
    preview: {
      allowedHosts: [
        '.ddns.net',
        'localhost',
      ],
      port: FRONTEND_PORT,
      host: '0.0.0.0',
      proxy: {
        '/api': {
          target: `http://backend:${BACKEND_PORT}`,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, '')
        }
      }
    },
    server: {
      proxy: {
        '/api': {
          target: mode === 'development' ? `http://localhost:${BACKEND_PORT}` : `http://backend:${BACKEND_PORT}`,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, '')
        }
      },
      allowedHosts: [
        '.ddns.net',
        'localhost'
      ],
      port: FRONTEND_PORT,
      host: '0.0.0.0',
    },
    build: {
      assetsInlineLimit: 0
    }
  }
});

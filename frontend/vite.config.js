import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import react from '@vitejs/plugin-react-swc';
import { defineConfig } from 'vite';
import {readFileSync} from 'fs'
//import mkcert from 'vite-plugin-mkcert'
//import basicSsl from '@vitejs/plugin-basic-ssl';

// https://vitejs.dev/config/
export default defineConfig({
  base: '/',
  plugins: [
    // Allows using React dev server along with building a React application with Vite.
    // https://npmjs.com/package/@vitejs/plugin-react-swc
    react(),
    // Allows using self-signed certificates to run the dev server using HTTPS.
    // https://www.npmjs.com/package/@vitejs/plugin-basic-ssl
    // basicSsl({
    //   name: 'localhost',
    //   domains: ['*'],
    //   certDir: '/workspaces/basic-nodejs/cert'
    // }),
    //mkcert()
  ],
  publicDir: './public',
  server: {
    // Exposes your dev server and makes it accessible for the devices in the same network.
    port: 4444,
    host: "0.0.0.0",
    /*https: {
      key: readFileSync('../.cert/localhost-key.pem'),
      cert: readFileSync('../.cert/localhost.pem'),
    },*/
  },
  resolve: {
    alias: {
      '@': resolve(dirname(fileURLToPath(import.meta.url)), './src'),
    }
  },
});


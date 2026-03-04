import legacyPlugin from '@vitejs/plugin-legacy'
import react from '@vitejs/plugin-react'
import * as path from 'path'
import * as dotenv from 'dotenv'
import * as fs from 'fs'

export default () => {
  const nodeEnv = process.env.NODE_ENV || 'development'
  const envFile = `.env.${nodeEnv}`
  if (fs.existsSync(envFile)) {
    const envConfig = dotenv.parse(fs.readFileSync(envFile))
    for (const k in envConfig) {
      process.env[k] = envConfig[k]
    }
  }

  return {
    base: './',
    root: './',
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      }
    },
    define: {
      'process.env': {}
    },
    server: {
      open: true,
      port: Number(process.env.VITE_CLI_PORT) || 8080,
      fs: {
        strict: false
      },
      proxy: {
        [process.env.VITE_BASE_API]: {
          target: `${process.env.VITE_BASE_PATH}:${process.env.VITE_SERVER_PORT}/`,
          changeOrigin: true,
          rewrite: (p) => p.replace(new RegExp('^' + process.env.VITE_BASE_API), '')
        }
      }
    },
    build: {
      target: 'es2017',
      minify: 'esbuild',
      sourcemap: false,
      outDir: 'dist'
    },
    plugins: [
      legacyPlugin({
        targets: ['Android > 39', 'Chrome >= 60', 'Safari >= 10.1', 'iOS >= 10.3', 'Firefox >= 54', 'Edge >= 15']
      }),
      react()
    ]
  }
}

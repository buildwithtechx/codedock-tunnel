import { defineConfig } from 'tsup';

export default defineConfig({
  clean: true,
  dts: true,
  entry: ['src/index.ts'],
  external: ['express'],
  format: ['esm'],
  sourcemap: true,
});

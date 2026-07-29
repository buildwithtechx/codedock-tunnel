import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const env = createEnv({
  server: {
    CODEDOCK_SERVER_URL: z.url().optional(),
  },

  clientPrefix: "VITE_",

  client: {
    VITE_CODEDOCK_APP_TITLE: z.string().min(1).optional(),
    VITE_CODEDOCK_TUNNEL_SERVER_URL: z.url().optional(),
    VITE_CODEDOCK_API_BASE_URL: z.url().optional(),
    VITE_CODEDOCK_POSTHOG_KEY: z.string().optional(),
    VITE_CODEDOCK_POSTHOG_HOST: z.url().optional(),
    VITE_CODEDOCK_POLAR_PRODUCT_RAY: z.string().optional(),
    VITE_CODEDOCK_POLAR_PRODUCT_BEAM: z.string().optional(),
    VITE_CODEDOCK_POLAR_PRODUCT_PULSE: z.string().optional(),
    VITE_CODEDOCK_POLAR_PRODUCT_RAY_YEARLY: z.string().optional(),
    VITE_CODEDOCK_POLAR_PRODUCT_BEAM_YEARLY: z.string().optional(),
    VITE_CODEDOCK_POLAR_PRODUCT_PULSE_YEARLY: z.string().optional(),
  },

  runtimeEnv: import.meta.env,

  emptyStringAsUndefined: true,
});

/// <reference types="vitest/config" />
import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import tanstackRouter from "@tanstack/router-plugin/vite";
import tailwindcss from "@tailwindcss/vite";
import path from "path";
import Icons from "unplugin-icons/vite";

// https://vitejs.dev/config/
export default ({ mode }: { mode: string }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return defineConfig({
    envPrefix: "VITE",
    plugins: [
      tanstackRouter({
        target: "react",
        autoCodeSplitting: true,
        routeToken: "layout",
      }),
      react(),
      tailwindcss(),
      Icons({
        compiler: "jsx",
        jsx: "react",
      }),
    ],
    test: {
      environment: "jsdom",
      setupFiles: ["./src/setupTests.ts"],
      globals: true,
    },
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    esbuild: {
      drop: mode === "production" ? ["console", "debugger"] : [],
    },
    server: {
      allowedHosts: process.env.VITE_ALLOWED_HOSTS
        ? JSON.parse(process.env.VITE_ALLOWED_HOSTS ?? "")
        : "",
    },
  });
};

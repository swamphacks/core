import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClientProvider } from "@tanstack/react-query";

import "./index.css";
import { routeTree } from "./routeTree.gen";
import { queryClient } from "./lib/query";
import { ThemeProvider } from "./components/ThemeProvider";
import { auth } from "./lib/authClient";
import { ToastContainer } from "react-toastify";

const router = createRouter({
  routeTree,
  context: {
    userQuery: undefined!,
  },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function App() {
  return (
    <ThemeProvider defaultTheme="system">
      <QueryClientProvider client={queryClient}>
        <InnerApp />
        <ToastContainer />
      </QueryClientProvider>
    </ThemeProvider>
  );
}

function InnerApp() {
  const userQuery = auth.useUser();

  return (
    <RouterProvider
      router={router}
      context={{
        userQuery,
      }}
    />
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);

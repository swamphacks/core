import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import NotFoundPage from "@/features/NotFound/NotFoundPage";
import "./index.css";
import { routeTree } from "./routeTree.gen";
import { ThemeProvider } from "./components/ThemeProvider";
import { auth } from "./lib/authClient";
import { ToastContainer } from "react-toastify";
import * as TanStackQueryProvider from "./integrations/tanstack-query/root-provider.tsx";
import "@smastrom/react-rating/style.css";

const TanStackQueryProviderContext = TanStackQueryProvider.getContext();
const router = createRouter({
  routeTree,
  context: {
    ...TanStackQueryProviderContext,
    userQuery: undefined!,
  },
  defaultNotFoundComponent: NotFoundPage,
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function App() {
  return (
    <ThemeProvider defaultTheme="system">
      <TanStackQueryProvider.Provider {...TanStackQueryProviderContext}>
        <InnerApp />
        <ToastContainer />
      </TanStackQueryProvider.Provider>
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

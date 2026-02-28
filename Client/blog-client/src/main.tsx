import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import { BrowserRouter } from "react-router-dom";
import CssBaseline from "@mui/material/CssBaseline";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { AuthProviderContainer } from "./context/AuthProviderContainer.tsx";

const querryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={querryClient}>
      <AuthProviderContainer>
        <BrowserRouter>
          <CssBaseline />
          <App />
          <ReactQueryDevtools />
        </BrowserRouter>
      </AuthProviderContainer>
    </QueryClientProvider>
  </StrictMode>
);

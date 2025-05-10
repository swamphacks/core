// src/setupTests.js (or .ts)
import "@testing-library/jest-dom";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";

// This is often implicitly handled by @testing-library/react with Vitest,
// but explicitly calling cleanup can be a good practice to ensure isolation
// between tests in case of specific scenarios or custom renderers.
afterEach(() => {
  cleanup();
});

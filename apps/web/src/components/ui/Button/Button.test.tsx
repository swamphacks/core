// src/components/Button.test.tsx
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Button } from ".";

describe("Button component", () => {
  it("renders with correct text", () => {
    render(<Button color="primary">Click Me!</Button>);
    const button = screen.getByText(/Click Me!/i);
    expect(button).toBeInTheDocument();
  });
});

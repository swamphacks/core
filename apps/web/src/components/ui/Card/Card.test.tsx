import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Card } from ".";

describe("Card component", () => {
  it("renders with correct title", () => {
    render(
      <Card>
        <p>Title</p>
      </Card>,
    );
    const card = screen.getByText(/Title/i);
    expect(card).toBeInTheDocument();
  });
});

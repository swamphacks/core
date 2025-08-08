import { describe, expect, it } from "vitest";
import { FormSchema } from "@/features/FormBuilder/formSchema";
import validJSON from "./valid.json";
import invalidJSON from "./invalid.json";
import invalidMissingMetadata from "./invalidMissingMetadata.json";
import invalidMoreThanTwoQuestionsInLayout from "./invalidMoreThanTwoQuestionsInLayout.json";
import invalidNestedLayouts from "./invalidNestedLayouts.json";
import invalidNestedSections from "./invalidNestedSections.json";
import invalidUnknownFieldType from "./invalidUnknownFieldType.json";

describe("JSON Form Builder", () => {
  it("parses valid form json correctly", () => {
    const { success, error, data } = FormSchema.safeParse(validJSON);

    expect(success).toBe(true);
    expect(data).toBeDefined();
    expect(error).toBeUndefined();
  });

  it("parses invalid form json correctly", () => {
    const { success, error, data } = FormSchema.safeParse(invalidJSON);

    expect(success).toBe(false);
    expect(data).toBeUndefined();
    expect(error).toBeDefined();
  });

  it("error for missing metadata", () => {
    const { success, error, data } = FormSchema.safeParse(
      invalidMissingMetadata,
    );

    expect(success).toBe(false);
    expect(data).toBeUndefined();
    expect(error).toBeDefined();
  });

  it("error for more than two questions in layout", () => {
    const { success, error, data } = FormSchema.safeParse(
      invalidMoreThanTwoQuestionsInLayout,
    );

    expect(success).toBe(false);
    expect(data).toBeUndefined();
    expect(error).toBeDefined();
  });

  it("error for having nested layouts", () => {
    const { success, error, data } = FormSchema.safeParse(invalidNestedLayouts);

    expect(success).toBe(false);
    expect(data).toBeUndefined();
    expect(error).toBeDefined();
  });

  it("error for having nested sections", () => {
    const { success, error, data } = FormSchema.safeParse(
      invalidNestedSections,
    );

    expect(success).toBe(false);
    expect(data).toBeUndefined();
    expect(error).toBeDefined();
  });

  it("error for having unknown field types (other than `section`, `layout`, `question`)", () => {
    const { success, error, data } = FormSchema.safeParse(
      invalidUnknownFieldType,
    );

    expect(success).toBe(false);
    expect(data).toBeUndefined();
    expect(error).toBeDefined();
  });
});

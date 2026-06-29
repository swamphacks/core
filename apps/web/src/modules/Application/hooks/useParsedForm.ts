import { useEffect, useState } from "react";
import applicationFormJSON from "../application.json";

export type ParsedForm = {
  [key: string]: Array<{ label: string; name: string }>;
};

export default function useParsedForm() {
  const [parsedForm, setParsedForm] = useState<ParsedForm | null>(null);

  useEffect(() => {
    const parsedForm: ParsedForm = {};

    function parseForm(content: any) {
      if (!content) return [];

      let questions: any = [];

      for (const item of content) {
        if (item.type === "section" || item.type === "layout") {
          if (item.type === "section") {
            parsedForm[item.label] = parseForm(item["content"]);
          } else {
            questions = [...questions, ...parseForm(item["content"])];
          }
        } else {
          if (item.name === "resume") {
            continue;
          }
          questions.push({ label: item.label, name: item.name });
        }
      }

      return questions;
    }

    parseForm(applicationFormJSON["content"]);
    setParsedForm(parsedForm);
  }, []);

  return parsedForm;
}

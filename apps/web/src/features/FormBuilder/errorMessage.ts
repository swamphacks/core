import { QuestionTypes } from "./types";

export const errorMessage = {
  [QuestionTypes.shortAnswer]: {
    required: "Required",
    tooLong: "Value is too long",
  },

  [QuestionTypes.paragraph]: {
    required: "Required",
    tooLong: "Value is too long",
    tooShort: "Value is too short",
  },

  [QuestionTypes.number]: {
    required: "Required",
    tooLow: "Value is too low",
    tooHigh: "Value is too short",
  },

  [QuestionTypes.multipleChoice]: {
    required: "Required",
  },

  [QuestionTypes.checkbox]: {
    required: "Required",
  },

  [QuestionTypes.select]: {
    required: "Required",
  },

  [QuestionTypes.multiselect]: {
    required: "Required",
  },

  [QuestionTypes.upload]: {
    required: "Required",
    invalidFileType: "Invalid file type",
    invalidSize: "File size is not within range",
  },

  [QuestionTypes.date]: {
    required: "Required",
  },

  [QuestionTypes.url]: {
    required: "Required",
    invalidURL: "Invalid URL",
    tooLong: "URL is too long",
  },
};

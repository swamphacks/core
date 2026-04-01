import { lazy } from "react";

export const textFieldIcons = {
  linkedin: lazy(() => import("~icons/tabler/brand-linkedin")),
  github: lazy(() => import("~icons/tabler/brand-github")),
  at: lazy(() => import("~icons/tabler/at")),
  phone: lazy(() => import("~icons/tabler/phone")),
  // add more if needed...
};

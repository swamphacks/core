# Styling & UI

The web app uses TailwindCSS v4 for utility-class styling, a CSS custom property-based theme system for light/dark mode, React Aria Components as the accessible component primitive layer, and Storybook for component development.

---

## TailwindCSS v4

Tailwind is integrated via the official Vite plugin — there is no `tailwind.config.js`. Configuration lives entirely inside CSS files using the v4 `@import` / `@theme` API.

**`vite.config.ts`**

```ts
import tailwindcss from "@tailwindcss/vite";

plugins: [tailwindcss()]
```

**`src/index.css`**

```css
@import "tailwindcss";
@plugin "tailwindcss-react-aria-components";
@plugin "tailwindcss-animate";

@import "./theme.css";

@tailwind utilities;

@custom-variant dark (&:where(.dark, .dark *));
```

Key points:

- `tailwindcss-react-aria-components` adds state variants (`pressed:`, `selected:`, `invalid:`, etc.) that map directly to React Aria's render props — used extensively in component styling.
- `tailwindcss-animate` provides animation utilities.
- The dark mode variant is **class-based**: `@custom-variant dark (&:where(.dark, .dark *))`. The `.dark` class is toggled on `<html>` by `ThemeProvider`.
- No separate Tailwind config file exists; all custom tokens are declared inside `@theme inline { ... }` in `theme.css`.

**Class merging** uses `clsx` + `tailwind-merge` via a shared utility:

```ts
// src/utils/cn.ts
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

**`tailwind-variants`** (`tv()`) is used in place of raw class strings for components that have multiple variants — it handles variant composition, compound variants, and default variants cleanly.

---

## Theme System

All design tokens are defined as CSS custom properties in `src/theme.css`. The file has three sections:

1. `:root` — light mode values
2. `.dark` — dark mode overrides
3. `@theme inline { ... }` — registers custom properties as Tailwind utility classes

### Core Tokens (light / dark)

| Token | Light | Dark |
|---|---|---|
| `--background` | `oklch(99.405% 0.00011 271.152)` | `var(--color-neutral-900)` |
| `--surface` | `oklch(97.015% 0.00011 271.152)` | `#202020` |
| `--border` | `var(--color-zinc-300)` | `var(--color-zinc-600)` |
| `--accent-main` | `var(--color-blue-500)` | `var(--color-blue-400)` |

### Typography

| Token | Light | Dark |
|---|---|---|
| `--text-main` | `var(--color-zinc-900)` | `var(--color-zinc-200)` |
| `--text-secondary` | `var(--color-zinc-600)` | `var(--color-zinc-400)` |
| `--text-link` | `var(--color-blue-500)` | `var(--color-blue-300)` |

The global font is **Figtree** (loaded from Google Fonts), applied to `html`, `body`, and `#root`. Components reference it via `font-figtree` (registered as `--font-figtree: "Figtree", sans-serif` in `@theme inline`).

### Button Tokens

Four button variants each have `default`, `hover`, `pressed`, and `disabled` states:

```
--button-primary / --button-primary-hover / --button-primary-pressed / --button-primary-disabled
--button-secondary / ...
--button-danger / ...
--button-success / ...
```

An `outline` variant uses `oklch` with alpha for a translucent blue tint.

### Badge Tokens

Badges express application-level status concepts directly as tokens. Each status has a `bg` and `text` pair:

```
--badge-bg-accepted    / --badge-text-accepted     → green
--badge-bg-rejected    / --badge-text-rejected     → red
--badge-bg-waitlisted  / --badge-text-waitlisted   → violet
--badge-bg-attending   / --badge-text-attending    → blue
--badge-bg-under-review / --badge-text-under-review → orange
--badge-bg-staff       / --badge-text-staff        → sky
--badge-bg-admin       / --badge-text-admin        → fuchsia
--badge-bg-completed   / --badge-text-completed    → indigo
```

Dark mode badge backgrounds use raw `oklch()` values (e.g. `oklch(0.3378 0.0595 20.68)` for rejected) to achieve appropriate contrast on dark surfaces.

### Input Tokens

```
--input-bg / --input-bg-disbaled
--input-border / --input-border-focused / --input-border-invalid / --input-border-disabled
--input-text-error / --input-text-disabled
```

### Tailwind Integration

Every custom property is re-exported through `@theme inline` so it becomes a Tailwind utility:

```css
@theme inline {
  --color-background: var(--background);
  --color-surface:    var(--surface);
  --color-text-main:  var(--text-main);
  --color-button-primary: var(--button-primary);
  /* ... */
}
```

This means you can write `bg-background`, `text-text-main`, `bg-button-primary`, `border-input-border`, etc. as regular Tailwind classes and they respond to the `.dark` class automatically.

---

## Dark Mode

Dark mode is managed by `ThemeProvider` (`src/components/ThemeProvider.tsx`). It supports three modes — `"light"`, `"dark"`, and `"system"` — persisted to `localStorage` under the key `"ui-theme"`.

On mount it reads the stored preference (or falls back to the system media query) and toggles the `dark` class on `<html>`. The companion `ThemeSwitch` component renders the Light / Dark toggle buttons visible in the navbar.

---

## Component Library

### React Aria Components

All interactive UI primitives are built on **`react-aria-components`** (v1.10). Components in `src/components/ui/` wrap React Aria primitives and apply Tailwind classes that respond to React Aria's render-prop state.

The `tailwindcss-react-aria-components` Tailwind plugin makes this ergonomic — it exposes state as variants:

```tsx
// src/components/ui/Button/Button.tsx
import { Button as RACButton, composeRenderProps } from "react-aria-components";
import { tv } from "tailwind-variants";

export const button = tv({
  base: "inline-flex cursor-pointer items-center justify-center rounded-md font-medium focus:outline-none gap-2",
  variants: {
    variant: {
      primary: "bg-button-primary hover:bg-button-primary-hover pressed:bg-button-primary-pressed text-white",
      secondary: "bg-button-secondary hover:bg-button-secondary-hover pressed:bg-button-secondary-pressed",
      danger:  "bg-button-danger hover:bg-button-danger-hover pressed:bg-button-danger-pressed text-white",
      icon:    "border-0 p-1 hover:bg-black/[5%] pressed:bg-black/10 dark:hover:bg-white/10",
      // ...
    },
    isDisabled: {
      true: "cursor-not-allowed bg-gray-200 dark:bg-neutral-700 text-text-main/30",
    },
    size: { sm: "py-2 px-4 text-sm", md: "py-2 px-4 text-base", lg: "py-2 px-4 text-lg" },
  },
  defaultVariants: { variant: "primary", size: "md" },
});

export function Button(props: ButtonProps) {
  return (
    <RACButton
      {...props}
      className={composeRenderProps(props.className, (className, renderProps) =>
        button({ ...renderProps, variant: props.variant, className })
      )}
    />
  );
}
```

`composeRenderProps` merges React Aria's live state object (`isDisabled`, `isPressed`, `isFocusVisible`, etc.) directly into the `tv()` variant resolver. This means accessibility states and visual states are always in sync — no manual ARIA attribute wiring needed.

The helper `composeTailwindRenderProps` in `src/components/ui/utils.ts` is a thin wrapper around this pattern for cases that don't need `tv()`:

```ts
export function composeTailwindRenderProps<T>(
  className: string | ((v: T) => string) | undefined,
  tw: string,
): string | ((v: T) => string) {
  return composeRenderProps(className, (className) => cn(tw, className));
}
```

### Available UI Components

All components live in `src/components/ui/` and are individually exported via `index.ts` files:

| Component | React Aria Primitive |
|---|---|
| `Button` | `Button` |
| `TextField` | `TextField`, `Input` |
| `Select` | `Select`, `ListBox`, `Popover` |
| `ComboBox` | `ComboBox` |
| `MultiSelect` | `TagGroup`, `ListBox` |
| `Checkbox` | `Checkbox` |
| `Radio` / `RadioGroup` | `Radio`, `RadioGroup` |
| `DatePicker` | `DatePicker` |
| `DateField` | `DateField` |
| `DateRangePicker` | `DateRangePicker` |
| `Calendar` | `Calendar` |
| `Dialog` / `Modal` | `Dialog`, `Modal` |
| `Popover` | `Popover` |
| `Menu` | `Menu` |
| `ListBox` | `ListBox` |
| `Slider` | `Slider` |
| `Switch` | `Switch` |
| `ProgressBar` | `ProgressBar` |
| `Badge` | `<span>` (via `forwardRef`) |
| `Avatar` / `AvatarStack` | Custom |
| `Card` | Custom |
| `Spinner` | Custom |
| `Separator` | Custom |
| `Tag` | Custom |

Form field primitives (`Label`, `Description`, `FieldError`, `FieldGroup`, `Input`) are shared across all form components from `src/components/ui/Field/Field.tsx`, keeping validation display and focus styling consistent.

### Field Validation Styling

Invalid state is driven by CSS tokens rather than hardcoded colours:

```tsx
// border-input-border-invalid resolves to --input-border-invalid
// which is red-600 (light) or red-300 (dark)
isFocusWithin: { true: "border-input-border-focused" },
isInvalid:     { true: "border-input-border-invalid" },
isDisabled:    { true: "border-input-border-disabled" },
```

### Icon System

Icons are sourced from **Iconify** collections via `unplugin-icons`. The `@iconify-json/tabler` and `@iconify-json/ic` sets are installed. Icons are imported as virtual modules:

```tsx
import TablerChevronDown from "~icons/tabler/chevron-down";
import TablerSun from "~icons/tabler/sun";
```

The Vite config registers the `Icons` plugin with `compiler: "jsx"` so icons render as React SVG components.

---

## Storybook

Storybook is used to develop and document UI components in isolation.

**Run:**

```bash
pnpm storybook
```

Starts on **port 6006** (`http://localhost:6006`).

**Build static site:**

```bash
pnpm build-storybook
```

### Configuration

`.storybook/main.ts` uses the `@storybook/react-vite` framework, so it shares the same Vite config (including Tailwind and icon plugins) as the main app.

`.storybook/preview.ts` imports `src/index.css` so all theme tokens and Tailwind utilities are available in every story:

```ts
import "../src/index.css";
```

### Story Location

Stories are colocated with their components following the pattern:

```
src/components/ui/Button/Button.stories.tsx
src/components/ui/Badge/Badge.stories.tsx
src/components/AppShell/stories/NavLink.stories.tsx
```

The glob in `main.ts` picks up all `*.stories.@(js|jsx|mjs|ts|tsx)` files under `src/`. Stories use the `"UI/ComponentName"` title convention (e.g. `title: "UI/Button"`) and include `tags: ["autodocs"]` to generate an automatic props table.

Stories are also used as Storybook interaction test targets via `@storybook/experimental-addon-test`.

---

## Design Conventions

- **Color naming follows semantic role, not hue.** Prefer `bg-button-primary` over `bg-blue-600`. This allows dark mode to swap values without touching component code.
- **Status colors are defined centrally.** Application statuses (accepted, rejected, waitlisted, etc.) map to dedicated badge and event-button token families. Never use ad-hoc colour classes for status indicators.
- **`color-mix(in oklab, ...)` is used for translucent variants** in event-button and outline-button tokens — this keeps tints perceptually uniform across light and dark contexts.
- **`oklch` is the preferred colour space** for bespoke values (surface, dark-mode badge backgrounds) to ensure predictable perceptual lightness.
- **Dark mode is class-based**, not `prefers-color-scheme` media query. The `ThemeProvider` resolves system preference at runtime and applies `.dark` to `<html>`, allowing user override.
- **`tailwind-merge` is always used** when class strings are conditionally composed, preventing specificity conflicts from duplicate utilities.

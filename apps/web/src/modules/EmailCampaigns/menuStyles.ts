/**
 * Menu item styling shared by the campaign menus.
 *
 * The shared dropdownItemStyles highlights on *focus* (blue), and react-aria
 * menus keep the last item focused after the pointer leaves, so a focus-based
 * highlight looks stuck. Highlight on hover instead, keep a highlight for
 * keyboard focus (focused + focus-visible), and neutralise the focus blue.
 * It also sets cursor-default, so restore the pointer.
 */
export const campaignMenuItemClasses = [
  "cursor-pointer",
  // kill the focus-pinned blue
  "data-[focused]:bg-transparent data-[focused]:text-zinc-700 dark:data-[focused]:text-zinc-200",
  // mouse highlight
  "data-[hovered]:bg-zinc-100 data-[hovered]:text-zinc-900 dark:data-[hovered]:bg-zinc-700 dark:data-[hovered]:text-zinc-100",
  // keyboard highlight (stacked variant outranks the transparent override)
  "data-[focused]:data-[focus-visible]:bg-zinc-100 dark:data-[focused]:data-[focus-visible]:bg-zinc-700",
].join(" ");

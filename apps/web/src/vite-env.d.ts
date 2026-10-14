/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />
import "@tanstack/react-table";

type FilterOption = {
  value: string;
  label: string;
};

declare module "@tanstack/react-table" {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface ColumnMeta<TData extends RowData, TValue> {
    filterType?: "text" | "select";
    filterOptions?: FilterOption[];
    responsiveClass?: string;
  }
}

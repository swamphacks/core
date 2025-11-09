/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />
import "@tanstack/react-table";

type FilterOption = {
  value: string;
  label: string;
};

declare module "@tanstack/react-table" {
  interface ColumnMeta<TData extends RowData, TValue> {
    // eslint-disable-line @typescript-eslint/no-unused-vars
    filterType?: "text" | "select";
    filterOptions?: FilterOption[];
    responsiveClass?: string;
  }
}

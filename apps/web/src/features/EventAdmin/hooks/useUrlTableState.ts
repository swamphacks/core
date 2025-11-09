import { useState, useMemo, useEffect } from "react";
import debounce from "lodash.debounce";
import type { ColumnFiltersState, SortingState } from "@/components/ui/Table"; // Adjust this path as needed

interface TableState {
  filters: ColumnFiltersState;
  sorting: SortingState;
}

// Helper functions for handling URL filter encoding
function parseTableState(encodedString: string | undefined): TableState {
  const defaults: TableState = {
    filters: [],
    sorting: [],
  };
  if (encodedString) {
    try {
      const decoded = atob(encodedString);
      const parsed = JSON.parse(decoded);
      return { ...defaults, ...parsed };
    } catch (e) {
      console.error("Failed to decode filters from URL:", e);
    }
  }
  return defaults;
}

function encodeTableState(state: TableState): string {
  if (state.filters.length === 0 && state.sorting.length === 0) {
    return ""; // Return empty string to clear URL param
  }
  return btoa(JSON.stringify(state));
}

interface UseUrlStateProps<
  TSearch extends { tableState?: string | undefined },
> {
  search: TSearch;
  navigate: (options: {
    search: (prev: TSearch) => TSearch;
    replace: true;
  }) => void;
  debounceMs?: number;
}

export function useUrlTableState<
  TSearch extends { tableState?: string | undefined },
>({ search, navigate, debounceMs = 300 }: UseUrlStateProps<TSearch>) {
  const initialState = useMemo(
    () => parseTableState(search.tableState),
    [search.tableState],
  );

  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>(
    initialState.filters,
  );

  const [sorting, setSorting] = useState<SortingState>(initialState.sorting);
  const debouncedUrlUpdate = useMemo(
    () =>
      debounce(
        (
          filters: ColumnFiltersState,
          sort: SortingState,
          currentSearch: TSearch,
        ) => {
          const newState: TableState = { filters, sorting: sort };
          const newSearchState = encodeTableState(newState);

          const newSearchParam = newSearchState ? newSearchState : undefined;

          if (newSearchParam !== currentSearch.tableState) {
            navigate({
              search: (prev) => ({
                ...prev,
                tableState: newSearchParam,
              }),
              replace: true,
            });
          }
        },
        debounceMs,
      ),
    [navigate, debounceMs],
  );

  useEffect(() => {
    debouncedUrlUpdate(columnFilters, sorting, search);
  }, [columnFilters, sorting, search, debouncedUrlUpdate]);

  return {
    columnFilters,
    setColumnFilters,
    sorting,
    setSorting,
  };
}

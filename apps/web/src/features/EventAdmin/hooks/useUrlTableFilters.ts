import { useState, useMemo, useEffect } from "react";
import debounce from "lodash.debounce";
import type { ColumnFiltersState } from "@/components/ui/Table"; // Adjust this path as needed

// Helper functions for handling URL filter encoding
function parseEncodedFilters(
  encodedString: string | undefined,
): ColumnFiltersState {
  if (encodedString) {
    try {
      const decoded = atob(encodedString);
      return JSON.parse(decoded);
    } catch (e) {
      console.error("Failed to decode filters from URL:", e);
    }
  }
  return [];
}

function encodeFiltersForUrl(filters: ColumnFiltersState): string {
  return btoa(JSON.stringify(filters));
}

interface UseUrlFiltersProps<TSearch extends { filters?: string | undefined }> {
  search: TSearch;
  navigate: (options: {
    search: (prev: TSearch) => TSearch;
    replace: true;
  }) => void;
  /**
   * Debounce time in milliseconds.
   * @default 300
   */
  debounceMs?: number;
}

export function useUrlBoundTableFilters<
  TSearch extends { filters?: string | undefined },
>({ search, navigate, debounceMs = 300 }: UseUrlFiltersProps<TSearch>) {
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>(() => {
    return parseEncodedFilters(search.filters);
  });

  const debouncedUrlUpdate = useMemo(
    () =>
      debounce((filters: ColumnFiltersState, currentSearch: TSearch) => {
        const newSearchFilter =
          filters.length > 0 ? encodeFiltersForUrl(filters) : undefined;

        if (newSearchFilter !== currentSearch.filters) {
          navigate({
            search: (prev) => ({
              ...prev,
              filters: newSearchFilter,
            }),
            replace: true,
          });
        }
      }, debounceMs),
    [navigate, debounceMs],
  );

  useEffect(() => {
    debouncedUrlUpdate(columnFilters, search);
  }, [columnFilters, search, debouncedUrlUpdate]);

  return { columnFilters, setColumnFilters };
}

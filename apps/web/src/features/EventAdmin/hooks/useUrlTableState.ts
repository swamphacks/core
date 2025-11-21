import { useState, useMemo, useEffect } from "react";
import debounce from "lodash.debounce";
import type {
  ColumnFiltersState,
  SortingState,
  PaginationState,
} from "@tanstack/react-table";

const RESERVED_KEYS = ["page", "limit", "sort"];

interface UseUrlStateProps<TSearch> {
  search: TSearch;
  navigate: (options: {
    search: (prev: TSearch) => TSearch;
    replace: true;
  }) => void;
  debounceMs?: number;
}

export function useUrlTableState<TSearch extends Record<string, any>>({
  search,
  navigate,
  debounceMs = 300,
}: UseUrlStateProps<TSearch>) {
  const initialState = useMemo(() => {
    const pageIndex = Number(search.page) ? Number(search.page) - 1 : 0;
    const pageSize = Number(search.limit) ? Number(search.limit) : 10;

    let sorting: SortingState = [];
    if (search.sort && typeof search.sort === "string") {
      const [id, dir] = search.sort.split(".");
      sorting = [{ id, desc: dir === "desc" }];
    }

    // Map any non-reserved key to a column filter
    const filters: ColumnFiltersState = Object.keys(search)
      .filter(
        (key) => !RESERVED_KEYS.includes(key) && search[key] !== undefined,
      )
      .map((key) => ({
        id: key,
        value: search[key],
      }));

    return {
      pagination: { pageIndex, pageSize },
      sorting,
      filters,
    };
  }, [search]);

  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>(
    initialState.filters,
  );
  const [sorting, setSorting] = useState<SortingState>(initialState.sorting);
  const [pagination, setPagination] = useState<PaginationState>(
    initialState.pagination,
  );

  const updateUrl = useMemo(
    () =>
      debounce(
        (
          currentFilters: ColumnFiltersState,
          currentSorting: SortingState,
          currentPagination: PaginationState,
          currentSearch: TSearch,
        ) => {
          const newParams: Record<string, any> = {};

          if (currentPagination.pageIndex > 0) {
            newParams.page = currentPagination.pageIndex + 1;
          } else {
            newParams.page = undefined; // Remove from URL
          }

          if (currentPagination.pageSize !== 10) {
            newParams.limit = currentPagination.pageSize;
          } else {
            newParams.limit = undefined;
          }

          if (currentSorting.length > 0) {
            const { id, desc } = currentSorting[0];
            newParams.sort = `${id}.${desc ? "desc" : "asc"}`;
          } else {
            newParams.sort = undefined;
          }

          // Remove "stale" filters that aren't in url params
          Object.keys(currentSearch).forEach((key) => {
            if (!RESERVED_KEYS.includes(key)) {
              newParams[key] = undefined;
            }
          });

          currentFilters.forEach((filter) => {
            if (filter.value !== undefined && filter.value !== "") {
              newParams[filter.id] = filter.value;
            }
          });

          // Check if anything actually changed to avoid redundant navigations
          navigate({
            search: (prev) => ({
              ...prev,
              ...newParams,
            }),
            replace: true,
          });
        },
        debounceMs,
      ),
    [navigate, debounceMs],
  );

  useEffect(() => {
    updateUrl(columnFilters, sorting, pagination, search);
  }, [columnFilters, sorting, pagination, search, updateUrl]);

  return {
    columnFilters,
    setColumnFilters,
    sorting,
    setSorting,
    pagination,
    setPagination,
  };
}
